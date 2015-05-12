package genddl

import (
	"go/ast"
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

	for _, field := range structType.Fields {
		tableMap.addColumn(field)
		tableMap.addIndex(field)
	}

	return tableMap
}

func (tm *tableMap) addColumnOrIndex(field *ast.Field) {
	if field.Tag == nil {
		return
	}
	tagMap := tm.parseTag(field.Tag.Value)

	tm.addColumn(columnMap, tagMap)
	tm.addIndex(columnMap, tagMap)
}

type ColumnMap struct {
	Name   string
	Type   *ast.Ident
	TagMap map[string]string
}

func (tm *tableMap) addColumn(field *ast.Field, tagMap map[string]string) {
	columnMap = new(ColumnMap)

	if name, ok := tagMap["db"]; ok {
		columnMap.Name = name
	} else {
		return
	}

	if t, ok = field.Type.(*ast.Ident); ok {
		columnMap.Type = t
	} else {
		return
	}

	columnMap.TagMap = tagMap

	tm.Columns = append(tm.Columns, columnMap)
}

type IndexMap struct {
	IndexName  string
	Columns    []string
	Unique     bool
	PrimaryKey bool
}

func (tm *tableMap) addIndex(field *ast.Field, tagMap map[string]string) {
	indexMap = new(indexMap)

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

	if t, ok = field.Type.(*ast.Ident); ok {
		indexMap.Type = t
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
