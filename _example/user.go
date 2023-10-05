package example

import (
	"database/sql"
	"time"
)

//go:generate go run ../cmd/genddl/main.go -outpath=./mysql.sql -innerindex -uniquename -foreignkeyname -tablecollate=utf8mb4_general_ci
//go:generate go run ../cmd/genddl/main.go -outpath=./sqlite3.sql -driver=sqlite3

// User is a user of the service
// +table: user
type User struct {
	ID        uint32         `db:"id,primarykey,autoincrement"`
	Name      string         `db:"name,unique,size=255"`
	Age       sql.NullInt64  `db:"age"`
	Message   sql.NullString `db:"message"`
	IconImage []byte         `db:"icon_image"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}
