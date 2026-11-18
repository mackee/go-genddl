package selfref

import "github.com/mackee/go-genddl/index"

//genddl:table category
type Category struct {
	ID       int64 `db:"id,primarykey"`
	ParentID int64 `db:"parent_id"`
}

func (s Category) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.ForeignKey(s.ParentID, s.ID, index.ForeignKeyDeleteCascade),
	}
}
