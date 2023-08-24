package example

import "database/sql"

//genddl:view user_product
type UserProduct struct {
	ID               int64          `db:"p.id"`
	UserName         string         `db:"u.name"`
	ReceivedUserName sql.NullString `db:"ru.name"`
	ProductID        int64          `db:"p.id"`
	Type             int64          `db:"p.type"`
}

func (u UserProduct) _selectStatement() string {
	return `
SELECT p.id, u.name, ru.name, p.id, p.type FROM product AS p
  INNER JOIN user AS u ON p.user_id = u.id
  LEFT JOIN user AS ru ON p.received_user_id = ru.id
	`
}
