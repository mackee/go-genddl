package genddl

import (
	"log"
	"strings"

	"github.com/mackee/go-genddl/index"
)

type Sqlite3Dialect struct {
}

func (m Sqlite3Dialect) DriverName() string { return "sqlite3" }

func (m Sqlite3Dialect) ToSqlType(col *ColumnMap) string {
	column := ""

	switch col.TypeName {
	case "bool", "int", "int16", "int32", "int64", "uint16", "uint32", "uint64", "sql.NullBool", "sql.NullInt64", "sql.NullInt32", "sql.NullInt16", "sql.NullByte":
		column = "INTEGER"
	case "float32", "float64", "sql.NullFloat64":
		column = "REAL"
	case "string", "sql.NullString":
		column = "TEXT"
	case "time.Time", "mysql.NullTime", "sql.NullTime":
		column = "DATETIME"
	case "[]byte":
		column = "BLOB"
	case "driver.Valuer":
		if v, ok := col.TagMap["type"]; ok {
			column = v
		} else {
			log.Printf("[ERROR] must be defined type value: %s", col.Name)
		}
	default:
		log.Printf("[ERROR] undefined types: %s", col.TypeName)
	}

	if _, ok := col.TagMap["null"]; ok || strings.HasPrefix(col.TypeName, "sql.Null") || col.TypeName == "mysql.NullTime" {
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

	return column
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
		return "ON DELETE RISTRICT"
	case index.ForeignKeyDeleteCascade:
		return "ON DELETE CASCADE"
	case index.ForeignKeyDeleteSetNull:
		return "ON DELETE SET NULL"
	case index.ForeignKeyDeleteSetDefault:
		return "ON DELETE DEFAULT"
	case index.ForeignKeyDeleteNoAction:
		return "ON DELETE NO ACTION"
	case index.ForeignKeyUpdateRestrict:
		return "ON UPDATE RISTRICT"
	case index.ForeignKeyUpdateCascade:
		return "ON UPDATE CASCADE"
	case index.ForeignKeyUpdateSetNull:
		return "ON UPDATE SET NULL"
	case index.ForeignKeyUpdateSetDefault:
		return "ON UPDATE DEFAULT"
	case index.ForeignKeyUpdateNoAction:
		return "ON UPDATE NO ACTION"
	}
	return ""
}
