package genddl

import (
	"bytes"
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/mackee/go-genddl/index"
)

const (
	indexFuncName        = "_schemaIndex"
	indexReturnSliceType = "Definition"
)

type indexType int

const (
	indexUnique indexType = iota + 1
	indexPrimaryKey
	indexComplex
	indexForeign
)

type indexer interface {
	Index(dialect Dialect, tables map[*ast.StructType]string) string
}

type indexIdent struct {
	Struct            *ast.StructType
	Type              indexType
	Column            []indexColumn
	References        []indexColumn
	ForeignKeyOptions []index.ForeignKeyOption
}

func (si indexIdent) Index(dialect Dialect, tables map[*ast.StructType]string) string {
	bs := &bytes.Buffer{}
	bs.WriteString("    ")
	switch si.Type {
	case indexUnique:
		bs.WriteString("UNIQUE (")
	case indexPrimaryKey:
		bs.WriteString("PRIMARY KEY (")
	case indexComplex:
		bs.WriteString("INDEX (")
	case indexForeign:
		bs.WriteString("FOREIGN (")
	}
	columns := []string{}
	for _, column := range si.Column {
		columnName, err := column.Column(dialect, si.Struct, tables)
		if err != nil {
			log.Fatalf("[ERROR] cannot resolve column error: %s", err)
		}
		columns = append(columns, columnName)
	}
	bs.WriteString(strings.Join(columns, ", "))
	bs.WriteString(")")

	if si.Type == indexForeign {
		bs.WriteString(" REFERENCES ")
		if len(si.References) == 0 {
			log.Fatalf("[ERROR] specified references column is invalid")
		}
		references, err := si.References[0].Column(dialect, si.Struct, tables)
		if err != nil {
			log.Fatalf("[ERROR] cannot resolve foreign references column error: %s", err)
		}
		bs.WriteString(references)
		bs.WriteString(" ")

		var options []string
		for _, option := range si.ForeignKeyOptions {
			o := dialect.ForeignKey(option)
			options = append(options, o)
		}
		s := strings.Join(options, " ")
		bs.WriteString(s)
	}

	return bs.String()
}

type rawIndex string

func (rs rawIndex) Index(dialect Dialect, tables map[*ast.StructType]string) string {
	return "    " + string(rs)
}

type indexColumn interface {
	Column(dialect Dialect, me *ast.StructType, tables map[*ast.StructType]string) (string, error)
}

type unresolvedIndexColumn struct {
	StructName string
	Struct     *ast.StructType
	Field      *ast.Field
}

func (c unresolvedIndexColumn) bareColumn(dialect Dialect) string {
	field := c.Field

	tv, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		log.Fatalf("[ERROR] struct tag is not valid: %s", err)
	}
	tag := reflect.StructTag(tv)
	info := strings.SplitN(tag.Get("db"), ",", 2)
	columnName := info[0]
	return dialect.QuoteField(columnName)
}

func (c unresolvedIndexColumn) Column(dialect Dialect, me *ast.StructType, tables map[*ast.StructType]string) (string, error) {
	bareColumn := c.bareColumn(dialect)
	if me == c.Struct {
		return bareColumn, nil
	}
	if tableName, ok := tables[c.Struct]; ok {
		return tableName + "(" + bareColumn + ")", nil
	}

	return "", fmt.Errorf("specified column is not define table struct: %s.%s", c.StructName, c.Field.Names[0].Name)
}
