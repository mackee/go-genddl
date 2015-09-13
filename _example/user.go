package example

//go:generate go run ../cmd/genddl/main.go -outpath=./mysql.sql
//go:generate go run ../cmd/genddl/main.go -outpath=./sqlite3.sql -driver=sqlite3

//+table: user
type User struct {
	Id   uint32 `db:"id,primarykey,autoincrement"`
	Name string `db:"name,unique,size=255"`
}
