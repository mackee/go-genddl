package example

import (
	"github.com/mackee/go-genddl/index"
)

// Location is location of place.
// +table: location
type Location struct {
	ID          uint64 `db:"id,primarykey,autoincrement"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Place       []byte `db:"place,type=GEOMETRY,srid=4326"`
}

func (l Location) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.Spatial(l.Place),
		methods.Fulltext(l.Description),
	}
}
