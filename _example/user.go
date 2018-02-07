package example

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

//go:generate go run ../cmd/genddl/main.go -outpath=./mysql.sql
//go:generate go run ../cmd/genddl/main.go -outpath=./sqlite3.sql -driver=sqlite3

//+table: user
type User struct {
	ID        uint32         `db:"id,primarykey,autoincrement"`
	Name      string         `db:"name,unique,size=255"`
	Age       sql.NullInt64  `db:"age"`
	Message   sql.NullString `db:"message"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"updated_at"`
}
