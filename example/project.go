package example

//+table: project
type Project struct {
	Id   uint32 `db:"user", primaryKey`
	Name string `db:"name", unique, size:"255"`
	User uint32 `db:"user_id", default:"0"`
}
