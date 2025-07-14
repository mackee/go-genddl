package genddl

import "github.com/mackee/go-genddl/index"

type Dialect interface {
	ToSqlType(col *ColumnMap) (string, error)
	CreateTableSuffix() string
	QuoteField(field string) string
	DriverName() string
	ForeignKey(index.ForeignKeyOption) string
}

type DialectAddBeforeDefinitionStatement interface {
	Dialect
	BeforeDefinitionStatement(tm *TableMap, cm *ColumnMap) string
}
