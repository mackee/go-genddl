package genddl

import (
	"go/ast"
	"io"
	"strings"
)

type TableMap struct {
	Name    string
	Columns []*ColumnMap
	Indexes []*IndexMap
}

func NewTableMap(name string, structType *ast.StructType) *TableMap {
	tableMap := new(TableMap)
	tableMap.Name = name

	for _, field := range structType.Fields.List {
		tableMap.addColumnOrIndex(field)
	}

	return tableMap
}

func (tm *TableMap) WriteDDL(w io.Writer, dialect Dialect) error {
	tableName := dialect.QuoteField(strings.TrimSpace(tm.Name))

	_, err := io.WriteString(w, "DROP TABLE IF EXISTS "+tableName+";\n\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "CREATE TABLE "+tableName+" (\n")
	if err != nil {
		return err
	}

	for i, cm := range tm.Columns {
		columnName := dialect.QuoteField(cm.Name)
		columnType := dialect.ToSqlType(cm)

		comma := ",\n"
		if i == len(tm.Columns)-1 {
			comma = "\n"
		}

		_, err = io.WriteString(w, "    "+columnName+" "+columnType+comma)
		if err != nil {
			return err
		}
	}

	suffix := dialect.CreateTableSuffix()
	_, err = io.WriteString(w, ")"+suffix+";\n")
	if err != nil {
		return err
	}
	return nil
}

func (tm *TableMap) addColumnOrIndex(field *ast.Field) {
	if field.Tag == nil {
		return
	}
	tagMap := tm.parseTag(field.Tag.Value)

	tm.addColumn(field, tagMap)
	tm.addIndex(field, tagMap)
}

type ColumnMap struct {
	Name   string
	Type   *ast.Ident
	TagMap map[string]string
}

func (tm *TableMap) addColumn(field *ast.Field, tagMap map[string]string) {
	columnMap := new(ColumnMap)

	if name, ok := tagMap["db"]; ok {
		columnMap.Name = name
	} else {
		return
	}

	if t, ok := field.Type.(*ast.Ident); ok {
		columnMap.Type = t
	} else {
		return
	}

	columnMap.TagMap = tagMap

	tm.Columns = append(tm.Columns, columnMap)
}

type IndexMap struct {
	Name       string
	Columns    []string
	Unique     bool
	PrimaryKey bool
	TagMap     map[string]string
}

func (tm *TableMap) addIndex(field *ast.Field, tagMap map[string]string) {
	indexMap := new(IndexMap)

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

	tm.Indexes = append(tm.Indexes, indexMap)
}

func (tm *TableMap) parseTag(v string) map[string]string {
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
