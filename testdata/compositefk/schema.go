package compositefk

import "github.com/mackee/go-genddl/index"

//genddl:table items
type Item struct {
	ShopID int64  `db:"shop_id"`
	ID     int64  `db:"id"`
	Name   string `db:"name"`
}

func (s Item) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.PrimaryKey(s.ID),
		methods.Unique(s.ShopID, s.ID),
	}
}

//genddl:table a_item_review
type ItemReview struct {
	ID     int64 `db:"id"`
	ShopID int64 `db:"shop_id"`
	ItemID int64 `db:"item_id"`
}

func (s ItemReview) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.PrimaryKey(s.ID),
		methods.ForeignKey(
			index.Columns(s.ShopID, s.ItemID),
			index.Columns(Item{}.ShopID, Item{}.ID),
			index.ForeignKeyDeleteNoAction,
		),
	}
}
