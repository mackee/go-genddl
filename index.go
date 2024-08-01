package genddl

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
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
	indexNameMaxLength   = 32
)

type indexType int

const (
	indexUnique indexType = iota + 1
	indexPrimaryKey
	indexComplex
	indexForeign
	indexSpatial
	indexFulltext
)

type indexer interface {
	IsOuterOfCreateTable() bool
	IsPlaceOnEndOfDDLFile() bool
	Index(dialect Dialect, tables map[*ast.StructType]string) string
}

type indexIdent struct {
	Struct             *ast.StructType
	Type               indexType
	Column             []indexColumn
	References         []indexColumn
	ForeignKeyOptions  []index.ForeignKeyOption
	InnerComplexIndex  bool
	UniqueWithName     bool
	ForeignKeyWithName bool
	OuterForeignKey    bool
}

func (si indexIdent) IsOuterOfCreateTable() bool {
	if si.Type == indexComplex && !si.InnerComplexIndex {
		return true
	}
	if si.Type == indexForeign && si.OuterForeignKey {
		return true
	}
	return false
}

func (si indexIdent) IsPlaceOnEndOfDDLFile() bool {
	return si.Type == indexForeign && si.OuterForeignKey
}

func (si indexIdent) Index(dialect Dialect, tables map[*ast.StructType]string) string {
	bs := &bytes.Buffer{}
	switch si.Type {
	case indexUnique:
		if si.UniqueWithName {
			fmt.Fprintf(bs, "    UNIQUE %s (", joinAndStripName(si.Name()))
		} else {
			bs.WriteString("    UNIQUE (")
		}
	case indexPrimaryKey:
		bs.WriteString("    PRIMARY KEY (")
	case indexComplex:
		if si.InnerComplexIndex {
			fmt.Fprintf(bs, "    INDEX %s (", joinAndStripName(si.Name()))
		} else {
			tableName := tables[si.Struct]
			fmt.Fprintf(bs, "CREATE INDEX %s ON %s (", tableName+"_"+joinAndStripName(si.Name()), tableName)
		}
	case indexForeign:
		if si.OuterForeignKey {
			tableName := tables[si.Struct]
			fmt.Fprintf(bs, "ALTER TABLE %s ADD CONSTRAINT %s ", tableName, joinAndStripName(si.Name()))
		} else {
			bs.WriteString("    ")
		}
		if si.ForeignKeyWithName {
			fmt.Fprintf(bs, "FOREIGN KEY %s (", joinAndStripName(si.Name()))
		} else {
			bs.WriteString("FOREIGN KEY (")
		}
	case indexSpatial:
		fmt.Fprintf(bs, "    SPATIAL KEY %s (", joinAndStripName(si.Name()))
	case indexFulltext:
		fmt.Fprintf(bs, "    FULLTEXT KEY %s (", joinAndStripName(si.Name()))
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

func (si indexIdent) Name() string {
	var columnNames []string
	for _, column := range si.Column {
		columnNames = append(columnNames, column.ColumnName())
	}
	return strings.Join(columnNames, "_")
}

func joinAndStripName(s string) string {
	if len(s) <= indexNameMaxLength {
		return s
	}
	hs := sha1.Sum([]byte(s))
	he := hex.EncodeToString(hs[:])
	return s[:indexNameMaxLength-8] + he[:8]
}

type rawIndex string

func (rs rawIndex) IsOuterOfCreateTable() bool {
	return false
}

func (rs rawIndex) IsPlaceOnEndOfDDLFile() bool {
	return false
}

func (rs rawIndex) Index(dialect Dialect, tables map[*ast.StructType]string) string {
	return "    " + string(rs)
}

type indexColumn interface {
	Column(dialect Dialect, me *ast.StructType, tables map[*ast.StructType]string) (string, error)
	ColumnName() string
}

type unresolvedIndexColumn struct {
	StructName string
	Struct     *ast.StructType
	Field      *ast.Field
}

func (c unresolvedIndexColumn) ColumnName() string {
	field := c.Field

	tv, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		log.Fatalf("[ERROR] struct tag is not valid: %s", err)
	}
	tag := reflect.StructTag(tv)
	info := strings.SplitN(tag.Get("db"), ",", 2)
	columnName := info[0]
	return columnName
}

func (c unresolvedIndexColumn) bareColumn(dialect Dialect) string {
	return dialect.QuoteField(c.ColumnName())
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
