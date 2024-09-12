package genddl

import (
	"fmt"

	"github.com/mackee/go-genddl/index"
)

type Sqlite3Dialect struct{}

func (m Sqlite3Dialect) DriverName() string { return "sqlite3" }

func (m Sqlite3Dialect) ToSqlType(col *ColumnMap) (string, error) {
	column := ""

	switch col.TypeName {
	case "bool", "int", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
		column = "INTEGER"
	case "float32", "float64":
		column = "REAL"
	case "string":
		column = "TEXT"
	case "time.Time":
		column = "DATETIME"
	case "[]byte":
		if v, ok := col.TagMap["type"]; ok {
			column = v
		} else {
			column = "BLOB"
		}
	default:
		return "", fmt.Errorf("unsupported types: %s, column=%s", col.TypeName, col.Name)
	}

	if _, ok := col.TagMap["null"]; ok || col.IsNullable {
		column += " NULL"
	} else {
		column += " NOT NULL"
	}

	if v, ok := col.TagMap["default"]; ok {
		column += " DEFAULT " + v
	}
	if _, ok := col.TagMap["unique"]; ok {
		column += " UNIQUE"
	}
	if _, ok := col.TagMap["primarykey"]; ok {
		column += " PRIMARY KEY"
	}
	if _, ok := col.TagMap["autoincrement"]; ok {
		column += " AUTOINCREMENT"
	}

	return column, nil
}

func (m Sqlite3Dialect) CreateTableSuffix() string {
	return ""
}

func (m Sqlite3Dialect) QuoteField(field string) string {
	return `"` + field + `"`
}

func (m Sqlite3Dialect) ForeignKey(option index.ForeignKeyOption) string {
	switch option {
	case index.ForeignKeyDeleteRestrict:
		return "ON DELETE RESTRICT"
	case index.ForeignKeyDeleteCascade:
		return "ON DELETE CASCADE"
	case index.ForeignKeyDeleteSetNull:
		return "ON DELETE SET NULL"
	case index.ForeignKeyDeleteSetDefault:
		return "ON DELETE SET DEFAULT"
	case index.ForeignKeyDeleteNoAction:
		return "ON DELETE NO ACTION"
	case index.ForeignKeyUpdateRestrict:
		return "ON UPDATE RESTRICT"
	case index.ForeignKeyUpdateCascade:
		return "ON UPDATE CASCADE"
	case index.ForeignKeyUpdateSetNull:
		return "ON UPDATE SET NULL"
	case index.ForeignKeyUpdateSetDefault:
		return "ON UPDATE SET DEFAULT"
	case index.ForeignKeyUpdateNoAction:
		return "ON UPDATE NO ACTION"
	}
	return ""
}
