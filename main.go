package genddl

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Field struct {
	ColumnDef string
}

type Table struct {
	Name       string
	Fields     []Field
	PrimaryKey string
}

type Driver interface {
	Template() string

	TypeText() string
	TypeString(map[string]string) string
	TypeUInt32() string
	TypeUInt64() string
	TypeInt32() string
	TypeInt64() string
	TypeDateTime() string

	NotNull() string
	Null() string
	AutoIncrement() string
	Unique() string
	DefaultValue(map[string]string) string
}

func Run() {
	var schemadir, outpath, driverName string
	flag.StringVar(&schemadir, "schemadir", "", "schema declaretion directory")
	flag.StringVar(&outpath, "outpath", "", "schema target path")
	flag.StringVar(&driverName, "driver", "mysql", "target driver")

	flag.Parse()

	var driver Driver
	switch driverName {
	case "mysql":
		driver = Mysql{}
	default:
		log.Fatalf("undefined driver name: %s", driver)
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

	tableFields := map[string][]*ast.Field{}
	for tableName, specs := range tables {
		for _, spec := range specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					fields := structType.Fields.List
					tableFields[tableName] = fields
				}
			}
		}
	}

	tmpl, err := template.New("ddl").Parse(driver.Template())
	if err != nil {
		log.Println("template parse error:", err)
	}

	absOutpath, err := filepath.Abs(outpath)
	if err != nil {
		log.Println("target out path error:", err)
	}

	file, err := os.Create(absOutpath)
	if err != nil {
		log.Println("target out path error:", err)
	}

	for tableName, fields := range tableFields {
		parsedTable := Table{Name: strings.TrimSpace(tableName)}
		var parsedFields []Field
		for _, field := range fields {
			if field.Tag == nil {
				continue
			}
			tagMap := parseTag(field.Tag.Value)
			parsedField := Field{}

			var schemaDef []string
			if t, ok := field.Type.(*ast.Ident); ok {
				schemaDef, err = buildColumnDef(driver, t.Name, tagMap)
				if err != nil {
					continue
				}
				if _, ok := tagMap["primaryKey"]; ok {
					parsedTable.PrimaryKey = fmt.Sprintf("`%s`", tagMap["db"])
				}
			} else {
				log.Fatalf("unsupported type: %s", field.Type)
			}

			parsedField.ColumnDef = strings.Join(schemaDef, " ")

			parsedFields = append(parsedFields, parsedField)
		}
		parsedTable.Fields = parsedFields
		tmpl.Execute(file, parsedTable)
	}
}

func buildColumnDef(driver Driver, typeName string, tagMap map[string]string) ([]string, error) {
	var schemaDef []string

	if columnName, ok := tagMap["db"]; ok {
		schemaDef = append(schemaDef, fmt.Sprintf("`%s`", columnName))
	} else {
		return []string{}, errors.New("build column error: not defined column name")
	}

	switch typeName {
	case "string":
		if _, ok := tagMap["text"]; ok {
			schemaDef = append(schemaDef, driver.TypeText())
		} else {
			schemaDef = append(schemaDef, driver.TypeString(tagMap))
		}
	case "int32":
		schemaDef = append(schemaDef, driver.TypeInt32())
	case "int64":
		schemaDef = append(schemaDef, driver.TypeInt64())
	case "uint32":
		schemaDef = append(schemaDef, driver.TypeUInt32())
	case "uint64":
		if _, ok := tagMap["datetime"]; ok {
			schemaDef = append(schemaDef, driver.TypeDateTime())
		} else {
			schemaDef = append(schemaDef, driver.TypeUInt64())
		}
	default:
		log.Fatalf("unsupported type: %s", typeName)
	}

	if _, ok := tagMap["null"]; ok {
		schemaDef = append(schemaDef, driver.Null())
	} else {
		schemaDef = append(schemaDef, driver.NotNull())
	}

	if _, ok := tagMap["primarykey"]; ok {
		schemaDef = append(schemaDef, driver.AutoIncrement())
	} else if _, ok := tagMap["unique"]; ok {
		schemaDef = append(schemaDef, driver.Unique())
	}

	if _, ok := tagMap["default"]; ok {
		schemaDef = append(schemaDef, driver.DefaultValue(tagMap))
	}

	return schemaDef, nil
}

func parseTag(v string) map[string]string {
	tags := strings.Split(v, ",")
	tagMap := map[string]string{}
	for _, tag := range tags {
		kv := strings.Split(tag, ":")
		key := strings.Trim(kv[0], "` ")
		if len(kv) == 1 {
			tagMap[key] = ""
			continue
		}
		value := strings.Trim(kv[1], "`")
		tagMap[key] = strings.Replace(value, `"`, "", -1)
	}

	return tagMap
}
