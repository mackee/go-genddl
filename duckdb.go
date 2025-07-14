package genddl

import (
	"fmt"

	"github.com/mackee/go-genddl/index"
)

type DuckDBDialect struct{}

func (m DuckDBDialect) DriverName() string { return "duckdb" }

func (m DuckDBDialect) ToSqlType(col *ColumnMap) (string, error) {
	column := ""

	switch col.TypeName {
	case "bool":
		column = "BOOLEAN"
	case "int8":
		column = "TINYINT"
	case "uint8":
		column = "UTINYINT"
	case "int16":
		column = "SMALLINT"
	case "uint16":
		column = "USMALLINT"
	case "int", "int32":
		column = "INTEGER"
	case "uint32":
		column = "UINTEGER"
	case "int64":
		column = "BIGINT"
	case "uint64":
		column = "UBIGINT"
	case "float32":
		column = "FLOAT"
	case "float64":
		column = "DOUBLE"
	case "string":
		column = "VARCHAR"
	case "time.Time":
		column = "DATETIME"
	case "[]byte":
		if v, ok := col.TagMap["type"]; ok {
			column = v
		} else {
			column = "BLOB"
		}
	default:
		if v, ok := col.TagMap["type"]; ok {
			column = v
		} else {
			return "", fmt.Errorf("unsupported types: %s, column=%s", col.TypeName, col.Name)
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
		seqName := "seq_" + col.TableMap.Name + "_" + col.Name
		column += " DEFAULT nextval('" + seqName + "')"
	}

	return column, nil
}

func (m DuckDBDialect) CreateTableSuffix() string {
	return ""
}

func (m DuckDBDialect) QuoteField(field string) string {
	return `"` + field + `"`
}

func (m DuckDBDialect) ForeignKey(option index.ForeignKeyOption) string {
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

func (m DuckDBDialect) BeforeDefinitionStatement(tm *TableMap, cm *ColumnMap) string {
	if _, ok := cm.TagMap["autoincrement"]; !ok {
		return ""
	}
	stmt := fmt.Sprintf("CREATE SEQUENCE IF NOT EXISTS %s START WITH %d INCREMENT BY %d;",
		"seq_"+tm.Name+"_"+cm.Name, 1, 1)
	return stmt
}
