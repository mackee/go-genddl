package invalidfk

import "github.com/mackee/go-genddl/index"

//genddl:table workspace
type Workspace struct {
	ShopID int64 `db:"shop_id"`
	ID     int64 `db:"id"`
}

//genddl:table member
type Member struct {
	ID int64 `db:"id"`
}

//genddl:table count_mismatch
type CountMismatch struct {
	ShopID      int64 `db:"shop_id"`
	WorkspaceID int64 `db:"workspace_id"`
}

func (s CountMismatch) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.ForeignKey(index.Columns(s.ShopID, s.WorkspaceID), index.Columns(Workspace{}.ID)),
	}
}

//genddl:table no_columns
type NoColumns struct {
	ID int64 `db:"id"`
}

func (s NoColumns) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.ForeignKey(index.Columns(), index.Columns()),
	}
}

//genddl:table unresolved
type Unresolved struct {
	ShopID      int64 `db:"shop_id"`
	WorkspaceID int64 `db:"workspace_id"`
}

func (s Unresolved) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.ForeignKey(index.Columns(s.ShopID, "workspace_id"), index.Columns(Workspace{}.ShopID, Workspace{}.ID)),
	}
}

//genddl:table span_structs
type SpanStructs struct {
	WorkspaceID int64 `db:"workspace_id"`
	MemberID    int64 `db:"member_id"`
}

func (s SpanStructs) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.ForeignKey(index.Columns(s.WorkspaceID, s.MemberID), index.Columns(Workspace{}.ID, Member{}.ID)),
	}
}
