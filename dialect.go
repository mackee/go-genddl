package genddl

type Dialect interface {
	ToSqlType(col *ColumnMap) string
	CreateTableSuffix() string
	QuoteField(field string) string
	DriverName() string
}
