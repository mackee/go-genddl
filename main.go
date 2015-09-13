package genddl

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Field struct {
	ColumnDef string
}

type Table struct {
	Name       string
	Fields     []Field
	PrimaryKey string
}

func Run(from string) {
	fromdir := filepath.Dir(from)

	var schemadir, outpath, driverName string
	flag.StringVar(&schemadir, "schemadir", fromdir, "schema declaretion directory")
	flag.StringVar(&outpath, "outpath", "", "schema target path")
	flag.StringVar(&driverName, "driver", "mysql", "target driver")

	flag.Parse()

	var dialect Dialect
	switch driverName {
	case "mysql":
		dialect = MysqlDialect{}
	case "sqlite3":
		dialect = Sqlite3Dialect{}
	default:
		log.Fatalf("undefined driver name: %s", driverName)
	}

	path, err := filepath.Abs(schemadir)
	if err != nil {
		log.Println("filepath error:", err)
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(
		fset,
		path,
		func(finfo os.FileInfo) bool { return true },
		parser.ParseComments,
	)

	if err != nil {
		log.Println("schema parse error:", err)
	}

	var decls []ast.Decl
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			decls = append(decls, file.Decls...)
		}
	}

	tables := map[string][]ast.Spec{}
	for _, decl := range decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Doc == nil {
				continue
			}

			comment := genDecl.Doc.List[0]
			if strings.HasPrefix(comment.Text, "//+table:") {
				tableName := strings.TrimPrefix(comment.Text, "//+table:")
				tables[tableName] = genDecl.Specs
			}
		}
	}

	file, err := os.Create(outpath)
	if err != nil {
		log.Fatal("invalid outpath error:", err)
	}

	for tableName, specs := range tables {
		for _, spec := range specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					tableMap := NewTableMap(tableName, structType)
					if tableMap != nil {
						tableMap.WriteDDL(file, dialect)
					}
				}
			}
		}
	}

}
