package genddl

const (
	MYSQL_DEFAULT_VARCHAR_SIZE = "191"
)

type Dialect interface {
	ToSqlType(col *ColumnMap) string
	CreateTableSuffix() string
	QuoteField(field string) string
	DriverName() string
}

type MysqlDialect struct {
}

func (m MysqlDialect) DriverName() string { return "mysql" }

func (m MysqlDialect) ToSqlType(col *ColumnMap) string {
	column := ""

	switch col.Type.Name {
	case "bool":
		column = "BOOLEAN"
	case "int", "int16", "int32":
		column = "INTEGER"
	case "uint16", "uint32":
		column = "INTEGER unsigned"
	case "int64":
		column = "BIGINT"
	case "uint64":
		column = "BIGINT unsigned"
	case "float64":
		column = "DOUBLE"
	case "float32":
		column = "FLOAT"
	case "string":
		size := MYSQL_DEFAULT_VARCHAR_SIZE
		if v, ok := col.TagMap["size"]; ok {
			size = v
		}
		column = "VARCHAR(" + size + ")"
	}

	if _, ok := col.TagMap["null"]; ok {
		column += " NULL"
	} else {
		column += " NOT NULL"
	}

	if v, ok := col.TagMap["default"]; ok {
		column += " DEFAULT" + v
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
	return "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
}

func (m MysqlDialect) QuoteField(field string) string {
	return "`" + field + "`"
}
