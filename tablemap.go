package genddl

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/mackee/go-genddl/index"
)

type TableMap struct {
	Name                string
	Columns             []*ColumnMap
	ColumnIndexes       []*IndexMap
	Indexes             []indexer
	Tables              map[*ast.StructType]string
	EndOfDDLFileIndexes []indexer
}

type tableMapOption struct {
	innerIndexDef      bool
	uniqueWithName     bool
	foreignKeyWithName bool
	outerForeignKey    bool
	outerUniqueKey     bool
	withoutDropTable   bool
}

type tableMapOptionFunc func(*tableMapOption)

type tableMapOptionFuncs []tableMapOptionFunc

func (fs tableMapOptionFuncs) apply(o *tableMapOption) {
	for _, f := range fs {
		f(o)
	}
}

func withInnerIndexDef() tableMapOptionFunc {
	return func(o *tableMapOption) {
		o.innerIndexDef = true
	}
}

func withUniqueWithName() tableMapOptionFunc {
	return func(o *tableMapOption) {
		o.uniqueWithName = true
	}
}

func withForeignKeyWithName() tableMapOptionFunc {
	return func(o *tableMapOption) {
		o.foreignKeyWithName = true
	}
}

func withOuterForeignKey() tableMapOptionFunc {
	return func(o *tableMapOption) {
		o.outerForeignKey = true
	}
}

func withOuterUniqueKey() tableMapOptionFunc {
	return func(o *tableMapOption) {
		o.outerUniqueKey = true
	}
}

func NewTableMap(name string, structType *ast.StructType, funcs []*ast.FuncDecl, tables map[*ast.StructType]string, ti *types.Info, opts *tableMapOption) (*TableMap, error) {
	tableMap := new(TableMap)
	tableMap.Name = name

	tableMap.Indexes = retrieveIndexesByFuncs(funcs, structType, opts)
	tableMap.Tables = tables

	for _, field := range structType.Fields.List {
		if err := tableMap.addColumnOrIndex(field, ti); err != nil {
			return nil, fmt.Errorf("failed to add column or index: %w", err)
		}
	}

	return tableMap, nil
}

func retrieveIndexesByFuncs(funcs []*ast.FuncDecl, me *ast.StructType, opts *tableMapOption) []indexer {
	var f *ast.FuncDecl
	for _, funcDecl := range funcs {
		if funcDecl.Name.String() != indexFuncName {
			continue
		}
		res := funcDecl.Type.Results.List[0]
		ra, ok := res.Type.(*ast.ArrayType)
		if !ok {
			continue
		}
		if ra.Len != nil {
			continue
		}
		se, ok := ra.Elt.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if se.Sel.String() != indexReturnSliceType {
			continue
		}
		f = funcDecl
	}
	if f == nil {
		return nil
	}
	body := f.Body
	var rt *ast.ReturnStmt
	ast.Inspect(body, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.ReturnStmt:
			rt = t
			return false
		default:
			return true
		}
	})

	sis := make([]indexer, 0)
	for _, lit := range rt.Results {
		cl, ok := lit.(*ast.CompositeLit)
		if !ok {
			continue
		}
	OUTER:
		for _, elt := range cl.Elts {
			var si indexIdent
			switch n := elt.(type) {
			case *ast.CallExpr:
				se, ok := n.Fun.(*ast.SelectorExpr)
				if !ok {
					break OUTER
				}
				methodName := se.Sel.Name
				switch methodName {
				case "PrimaryKey":
					si = indexIdent{
						Struct: me,
						Type:   indexPrimaryKey,
						Column: retrieveIndexColumnByExpr(n.Args),
					}
				case "Unique":
					si = indexIdent{
						Struct:         me,
						Type:           indexUnique,
						Column:         retrieveIndexColumnByExpr(n.Args),
						UniqueWithName: opts.uniqueWithName,
						OuterUniqueKey: opts.outerUniqueKey,
					}
				case "Complex":
					si = indexIdent{
						Struct:            me,
						Type:              indexComplex,
						Column:            retrieveIndexColumnByExpr(n.Args),
						InnerComplexIndex: opts.innerIndexDef,
					}
				case "ForeignKey":
					options := make([]index.ForeignKeyOption, 0)
					if len(n.Args) >= 2 {
						options = retrieveIndexForeignKeyOptionByExpr(n.Args[2:])
					}
					si = indexIdent{
						Struct:             me,
						Type:               indexForeign,
						Column:             retrieveIndexColumnByExpr([]ast.Expr{n.Args[0]}),
						References:         retrieveIndexColumnByExpr([]ast.Expr{n.Args[1]}),
						ForeignKeyOptions:  options,
						ForeignKeyWithName: opts.foreignKeyWithName,
						OuterForeignKey:    opts.outerForeignKey,
					}
				case "Spatial":
					si = indexIdent{
						Struct:            me,
						Type:              indexSpatial,
						Column:            retrieveIndexColumnByExpr(n.Args),
						InnerComplexIndex: opts.innerIndexDef,
					}
				case "Fulltext":
					si = indexIdent{
						Struct:            me,
						Type:              indexFulltext,
						Column:            retrieveIndexColumnByExpr(n.Args),
						InnerComplexIndex: opts.innerIndexDef,
					}
				default:
					break OUTER
				}
				sis = append(sis, si)
			case *ast.BasicLit:
				if n.Kind != token.STRING {
					break OUTER
				}
				v := n.Value
				v, _ = strconv.Unquote(v)
				sis = append(sis, rawIndex(v))
			}
		}
	}
	return sis
}

func retrieveIndexColumnByExpr(exprs []ast.Expr) []indexColumn {
	sc := make([]indexColumn, 0)
OUTER:
	for _, expr := range exprs {
		switch e := expr.(type) {
		case *ast.SelectorExpr:
			var st *ast.StructType
			var structName string
			if ident, ok := e.X.(*ast.Ident); ok {
				if sf, ok := ident.Obj.Decl.(*ast.Field); ok {
					sident, ok := sf.Type.(*ast.Ident)
					if !ok {
						continue OUTER
					}
					structName = sident.Name
					st, ok = typeNameStructMap[structName]
					if !ok {
						continue OUTER
					}
				} else if as, ok := ident.Obj.Decl.(*ast.AssignStmt); ok {
					cl, ok := as.Rhs[0].(*ast.CompositeLit)
					if !ok {
						continue OUTER
					}
					cident, ok := cl.Type.(*ast.Ident)
					if !ok {
						continue OUTER
					}
					structName = cident.Name
					st, ok = typeNameStructMap[structName]
					if !ok {
						continue OUTER
					}
				} else {
					continue OUTER
				}
			} else if cl, ok := e.X.(*ast.CompositeLit); ok {
				cident, ok := cl.Type.(*ast.Ident)
				if !ok {
					continue OUTER
				}
				structName = cident.Name
				st, ok = typeNameStructMap[structName]
				if !ok {
					continue OUTER
				}
			} else {
				continue OUTER
			}

			fieldName := e.Sel.Name
			var selField *ast.Field
			for _, field := range st.Fields.List {
				if field.Names[0].Name == fieldName {
					selField = field
					break
				}
			}
			column := unresolvedIndexColumn{
				StructName: structName,
				Struct:     st,
				Field:      selField,
			}
			sc = append(sc, column)
		}
	}
	return sc
}

func retrieveIndexForeignKeyOptionByExpr(exprs []ast.Expr) []index.ForeignKeyOption {
	options := make([]index.ForeignKeyOption, 0, len(exprs))
	for _, expr := range exprs {
		se, ok := expr.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		xident, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}
		if xident.Name != "index" {
			continue
		}
		for i := 0; i <= int(index.ForeignKeyOptionConstMax); i++ {
			if se.Sel.Name == index.ForeignKeyOption(i).String() {
				options = append(options, index.ForeignKeyOption(i))
				break
			}
		}
	}
	return options
}

func (tm *TableMap) WriteDDL(w io.Writer, dialect Dialect, opts *tableMapOption) error {
	tableName := dialect.QuoteField(strings.TrimSpace(tm.Name))

	if !opts.withoutDropTable {
		if _, err := io.WriteString(w, "DROP TABLE IF EXISTS "+tableName+";\n\n"); err != nil {
			return fmt.Errorf("failed to write drop table string: %w", err)
		}
	}
	if _, err := io.WriteString(w, "CREATE TABLE "+tableName+" (\n"); err != nil {
		return fmt.Errorf("failed to write create table string: %w", err)
	}

	innerIndexes := make([]indexer, 0, len(tm.Indexes))
	outerIndexes := make([]indexer, 0, len(tm.Indexes)/2)
	for _, index := range tm.Indexes {
		if index.IsOuterOfCreateTable() {
			outerIndexes = append(outerIndexes, index)
		} else {
			innerIndexes = append(innerIndexes, index)
		}
	}
	comma := ",\n"
	remainLines := len(tm.Columns) + len(innerIndexes)
	for _, cm := range tm.Columns {
		remainLines--
		if remainLines == 0 {
			comma = "\n"
		}

		columnName := dialect.QuoteField(cm.Name)
		columnType, err := dialect.ToSqlType(cm)
		if err != nil {
			return fmt.Errorf("failed to convert to sql type: %w", err)
		}

		_, err = io.WriteString(w, "    "+columnName+" "+columnType+comma)
		if err != nil {
			return fmt.Errorf("failed to write column string: %w", err)
		}
	}

	for _, sf := range innerIndexes {
		remainLines--
		if remainLines == 0 {
			comma = "\n"
		}
		str := sf.Index(dialect, tm.Tables)
		_, err := io.WriteString(w, str+comma)
		if err != nil {
			return fmt.Errorf("failed to write index string: %w", err)
		}
	}

	suffix := dialect.CreateTableSuffix()
	if _, err := io.WriteString(w, ") "+suffix+";\n"); err != nil {
		return err
	}

	tm.EndOfDDLFileIndexes = make([]indexer, 0, len(outerIndexes))
	for _, sf := range outerIndexes {
		if sf.IsPlaceOnEndOfDDLFile() {
			tm.EndOfDDLFileIndexes = append(tm.EndOfDDLFileIndexes, sf)
			continue
		}
		str := sf.Index(dialect, tm.Tables)
		_, err := io.WriteString(w, str+";\n")
		if err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "\n"); err != nil {
		return err
	}

	return nil
}

func (tm *TableMap) addColumnOrIndex(field *ast.Field, ti *types.Info) error {
	if field.Tag == nil {
		return nil
	}
	tv := field.Tag.Value

	if err := tm.addColumn(field, tv, ti); err != nil {
		return fmt.Errorf("failed to add column: %w", err)
	}
	return nil
}

type ColumnMap struct {
	Name       string
	TypeName   string
	IsNullable bool
	TagMap     map[string]string
}

func (tm *TableMap) addColumn(field *ast.Field, tag string, ti *types.Info) error {
	t := ti.TypeOf(field.Type)
	cs, err := columnsByFields(t, ti, tag, "")
	if err != nil {
		return fmt.Errorf("failed to get columns by fields: %w", err)
	}
	for _, c := range cs {
		tm.Columns = append(tm.Columns, &ColumnMap{
			Name:       c.name,
			TypeName:   c.typeName,
			IsNullable: c.isNullable,
			TagMap:     c.tagMap,
		})
		tm.addIndex(c.tagMap)
	}
	return nil
}

type IndexMap struct {
	Name       string
	Columns    []string
	Unique     bool
	PrimaryKey bool
	TagMap     map[string]string
}

func (tm *TableMap) addIndex(tagMap map[string]string) {
	indexMap := &IndexMap{}

	if name, ok := tagMap["index"]; ok {
		indexMap.Name = name
	} else if name, ok := tagMap["unique"]; ok {
		indexMap.Name = name
		indexMap.Unique = true
	} else if name, ok := tagMap["primarykey"]; ok {
		indexMap.Name = name
		indexMap.PrimaryKey = true
	} else {
		return
	}

	indexMap.TagMap = tagMap

	tm.ColumnIndexes = append(tm.ColumnIndexes, indexMap)
}

func parseTag(v string) map[string]string {
	st := reflect.StructTag(strings.Replace(v, "`", "", 2))
	dbTag := st.Get("db")
	if dbTag == "" {
		return nil
	}
	tags := strings.Split(dbTag, ",")
	tagMap := map[string]string{}

	tagMap["db"] = tags[0]
	for _, tag := range tags[1:] {
		kv := strings.Split(tag, "=")
		key := strings.TrimSpace(kv[0])
		if len(kv) == 1 {
			tagMap[key] = ""
			continue
		}
		tagMap[key] = strings.TrimSpace(kv[1])
	}

	return tagMap
}
