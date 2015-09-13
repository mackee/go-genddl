package genddl

type Sqlite3Dialect struct {
}

func (m Sqlite3Dialect) DriverName() string { return "sqlite3" }

func (m Sqlite3Dialect) ToSqlType(col *ColumnMap) string {
	column := ""

	switch col.Type.Name {
	case "bool", "int", "int16", "int32", "int64", "uint16", "uint32", "uint64":
		column = "INTEGER"
	case "float32", "float64":
		column = "REAL"
	case "string":
		column = "TEXT"
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
