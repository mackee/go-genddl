package genddl

import (
	"fmt"

	"github.com/mackee/go-genddl/index"
)

type PostgresqlDialect struct {
	Collate string
}

func (m PostgresqlDialect) DriverName() string { return "pg" }

func (m PostgresqlDialect) ToSqlType(col *ColumnMap) (string, error) {
	column := ""
	if v, ok := col.TagMap["def"]; ok {
		return v, nil
	}

	switch col.TypeName {
	case "bool":
		column = "BOOLEAN"
	case "int8", "uint8":
		column = "SMALLINT"
	case "int", "int16", "int32", "uint", "uint16":
		column = "INTEGER"
	case "uint32", "int64":
		column = "BIGINT"
	case "uint64":
		column = "NUMERIC(20)"
	case "float64":
		column = "DOUBLE precision"
	case "float32":
		column = "REAL"
	case "string":
		if _, ok := col.TagMap["text"]; ok {
			column = "TEXT"
		} else {
			tname := "VARCHAR"
			if v, ok := col.TagMap["type"]; ok {
				tname = v
			}
			if size, ok := col.TagMap["size"]; ok {
				column = tname + "(" + size + ")"
			} else {
				column = tname
			}
		}
		if m.Collate != "" {
			column += " COLLATE " + m.Collate
		}
	case "time.Time":
		column = "TIMESTAMP"
		if v, ok := col.TagMap["precision"]; ok {
			column += "(" + v + ")"
		}
		if _, ok := col.TagMap["withouttimezone"]; ok {
			column += " WITHOUT TIME ZONE"
		} else {
			column += " WITH TIME ZONE"
		}
	case "[]byte":
		if v, ok := col.TagMap["type"]; ok {
			column = v
		} else {
			column = "BYTEA"
		}
	default:
		if v, ok := col.TagMap["type"]; ok {
			column = v
		} else {
			return "", fmt.Errorf("unsupported type: %s, column=%s", col.TypeName, col.Name)
		}
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
		column += " GENERATED BY DEFAULT AS IDENTITY"
	}

	return column, nil
}

func (m PostgresqlDialect) CreateTableSuffix() string {
	return ""
}

func (m PostgresqlDialect) QuoteField(field string) string {
	return `"` + field + `"`
}

func (m PostgresqlDialect) ForeignKey(option index.ForeignKeyOption) string {
	switch option {
	case index.ForeignKeyDeleteRestrict:
		return "ON DELETE RISTRICT"
	case index.ForeignKeyDeleteCascade:
		return "ON DELETE CASCADE"
	case index.ForeignKeyDeleteSetNull:
		return "ON DELETE SET NULL"
	case index.ForeignKeyDeleteSetDefault:
		return "ON DELETE SET DEFAULT"
	case index.ForeignKeyDeleteNoAction:
		return "ON DELETE NO ACTION"
	case index.ForeignKeyUpdateRestrict:
		return "ON UPDATE RISTRICT"
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
