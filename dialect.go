package genddl

import "io"

type Dialect interface {
	ToSqlType(col *ColumnMap) string
	AutoIncrStr() string
	AutoIncrInsertSuffix(col *ColumnMap) string
	CreateTableSuffix() string
	QuoteField(field string) string
	RestartIdentityClause(table string) string
	DriverName() string

	WriteDDL(w io.Writer) error
}
