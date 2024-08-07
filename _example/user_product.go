package example

import "database/sql"

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
  INNER JOIN user AS u ON p.user_id = u.id
  LEFT JOIN user AS ru ON p.received_user_id = ru.id
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
  INNER JOIN user AS u ON p.user_id = u.id
  LEFT JOIN user AS ru ON p.received_user_id = ru.id
	`
}
