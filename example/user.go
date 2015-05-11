package example

//go:generate go run ../cmd/main.go -schemadir=./ -outpath=./mysql.sql

//+table: user
type User struct {
	Id   uint32 `db:"user", primaryKey`
	Name string `db:"name", unique, size:"255"`
}
