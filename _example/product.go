package example

import (
	"time"

	"github.com/mackee/go-genddl/index"
)

//+table: product
type Product struct {
	ID        uint32    `db:"id,autoincrement"`
	Name      string    `db:"name"`
	Type      uint32    `db:"type"`
	UserID    uint32    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (s Product) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.PrimaryKey(s.ID, s.CreatedAt),
		methods.Unique(s.UserID, s.Type),
		methods.Complex(s.UserID, s.CreatedAt),
		methods.ForeignKey(s.UserID, User{}.ID, index.ForeignKeyDeleteCascade, index.ForeignKeyUpdateCascade),
	}
}
