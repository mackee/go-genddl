package examplepg

import (
	"database/sql"
	"time"

	"github.com/mackee/go-genddl/index"
)

//go:generate go run ../../cmd/genddl/main.go -outpath=./postgresql.sql -outerforeignkey -withoutdroptable -driver=pg

type UserID int64

// User is a user of the service
//
//genddl:table user
type User struct {
	ID        UserID         `db:"id,primarykey,autoincrement"`
	Name      string         `db:"name,unique,size=255"`
	Age       sql.NullInt64  `db:"age"`
	Message   sql.NullString `db:"message"`
	IconImage []byte         `db:"icon_image"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

//genddl:view user_product
type UserProduct struct {
	ID               int64          `db:"u_id"`
	UserName         string         `db:"u_name"`
	ReceivedUserName sql.NullString `db:"ru_name"`
	ProductID        int64          `db:"p_id"`
	Type             int64          `db:"p_type"`
}

func (u UserProduct) _selectStatement() string {
	return `
SELECT u.id, u.name, ru.name, p.id, p.type FROM product AS p
  INNER JOIN "user" AS u ON p.user_id = u.id
  LEFT JOIN "user" AS ru ON p.received_user_id = ru.id
	`
}

//genddl:view user_product_structured
type UserProductStructured struct {
	Product      Product `db:"p_,nested"`
	User         User    `db:"u_,nested"`
	ReceivedUser User    `db:"ru_,nested"`
}

func (u UserProductStructured) _selectStatement() string {
	return `
SELECT p.*, u.*, ru.* FROM product AS p
  INNER JOIN "user" AS u ON p.user_id = u.id
  LEFT JOIN "user" AS ru ON p.received_user_id = ru.id
	`
}

// Location is location of place.
//
//genddl:table location
type Location struct {
	ID          int64  `db:"id,primarykey,autoincrement"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Place       []byte `db:"place"`
}

func (l Location) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{}
}

// Product is product of user
//
//genddl:table product
type Product struct {
	ID              uint32              `db:"id,primarykey,autoincrement"`
	Name            string              `db:"name"`
	Type            uint32              `db:"type"`
	UserID          uint32              `db:"user_id"`
	ReceivedUserID  sql.Null[UserID]    `db:"received_user_id"`
	Description     string              `db:"description,text"`
	FullDescription string              `db:"full_description,mediumtext"`
	Size            sql.NullInt16       `db:"size"`
	Status          uint8               `db:"status"`
	Category        int8                `db:"category"`
	CreatedAt       time.Time           `db:"created_at,precision=6"`
	UpdatedAt       sql.NullTime        `db:"updated_at,precision=6"`
	RemovedAt       sql.Null[time.Time] `db:"removed_at,precision=6"`

	Hidden           string `db:"-"`       // this field is ignored
	ExportedOtherTag int32  `json:"other"` // no tag field is ignored
	ExportedIgnore   string // no tag field is ignored
	unexported       bool   // unexported field is also ignored
}

func (s Product) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.Unique(s.UserID, s.Type),
		methods.Complex(s.UserID, s.CreatedAt),
		methods.ForeignKey(s.UserID, User{}.ID, index.ForeignKeyDeleteCascade, index.ForeignKeyUpdateCascade),
	}
}
