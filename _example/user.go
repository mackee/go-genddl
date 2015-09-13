package example

//go:generate go run ../cmd/genddl/main.go -outpath=./mysql.sql

//+table: user
type User struct {
	Id   uint32 `db:"user,primarykey"`
	Name string `db:"name,unique,size=255"`
}
