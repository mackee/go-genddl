package genddl

import "github.com/mackee/go-genddl/index"

type Dialect interface {
	ToSqlType(col *ColumnMap) string
	CreateTableSuffix(collate string) string
	QuoteField(field string) string
	DriverName() string
	ForeignKey(index.ForeignKeyOption) string
}
