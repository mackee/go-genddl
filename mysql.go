package genddl

var (
	MysqlVarcharDefaultSize = "191"
)

type Mysql struct {
}

func (m Mysql) Template() string {
	return `DROP TABLE IF EXISTS {{ .Name }};

CREATE TABLE {{ .Name }} (
    {{ range .Fields }}{{ .ColumnDef }},
    {{ end }}PRIMARY KEY ({{ .PrimaryKey }})
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`

}

func (m Mysql) TypeText() string {
	return "TEXT"
}
func (m Mysql) TypeString(opts map[string]string) string {
	size := MysqlVarcharDefaultSize
	if v, ok := opts["size"]; ok {
		size = v
	}

	return "VARCHAR(" + size + ")"
}
func (m Mysql) TypeUInt32() string {
	return "INTEGER unsigned"
}
func (m Mysql) TypeInt32() string {
	return "INTEGER"
}
func (m Mysql) TypeUInt64() string {
	return "BIGINT unsigned"
}
func (m Mysql) TypeInt64() string {
	return "BIGINT"
}

func (m Mysql) TypeDateTime() string {
	return "DATETIME"
}

func (m Mysql) NotNull() string {
	return "NOT NULL"
}
func (m Mysql) Null() string {
	return "NULL"
}

func (m Mysql) AutoIncrement() string {
	return "AutoIncrement"
}
func (m Mysql) Unique() string {
	return "UNIQUE"
}
