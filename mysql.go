package genddl

import (
	"fmt"
	"strings"

	"github.com/mackee/go-genddl/index"
)

const (
	MYSQL_DEFAULT_VARCHAR_SIZE = "191"
)

type MysqlDialect struct {
	Collate string
}

func (m MysqlDialect) DriverName() string { return "mysql" }

func (m MysqlDialect) ToSqlType(col *ColumnMap) string {
	column := ""

	switch col.TypeName {
	case "bool", "sql.NullBool":
		column = "BOOLEAN"
	case "int8":
		column = "TINYINT"
	case "uint8":
		column = "TINYINT unsigned"
	case "int", "int16", "sql.NullInt16", "int32", "sql.NullInt32", "sql.NullByte":
		column = "INTEGER"
	case "uint16", "uint32":
		column = "INTEGER unsigned"
	case "int64", "sql.NullInt64":
		column = "BIGINT"
	case "uint64":
		column = "BIGINT unsigned"
	case "float64", "sql.NullFloat64":
		column = "DOUBLE"
	case "float32":
		column = "FLOAT"
	case "string", "sql.NullString":
		if _, ok := col.TagMap["text"]; ok {
			column = "TEXT"
		} else if _, ok := col.TagMap["mediumtext"]; ok {
			column = "MEDIUMTEXT"
		} else {
			size := MYSQL_DEFAULT_VARCHAR_SIZE
			if v, ok := col.TagMap["size"]; ok {
				size = v
			}
			column = "VARCHAR(" + size + ")"
		}
		if m.Collate != "" {
			column += " COLLATE " + m.Collate
		}
	case "time.Time", "sql.NullTime", "mysql.NullTime":
		column = "DATETIME"
		if v, ok := col.TagMap["precision"]; ok {
			column += "(" + v + ")"
		}
	case "[]byte":
		if v, ok := col.TagMap["type"]; ok {
			column = v
		} else {
			column = "BLOB"
		}
	}

	if _, ok := col.TagMap["null"]; ok || strings.HasPrefix(col.TypeName, "sql.Null") || col.TypeName == "mysql.NullTime" {
		column += " NULL"
	} else {
		column += " NOT NULL"
	}

	if v, ok := col.TagMap["srid"]; ok {
		column += " SRID " + v
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
		column += " AUTO_INCREMENT"
	}

	return column
}

func (m MysqlDialect) CreateTableSuffix() string {
	if m.Collate != "" {
		return fmt.Sprintf("ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=%s", m.Collate)
	}
	return "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
}

func (m MysqlDialect) QuoteField(field string) string {
	return "`" + field + "`"
}

func (m MysqlDialect) ForeignKey(option index.ForeignKeyOption) string {
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
