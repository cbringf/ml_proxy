package mysql

import (
	"database/sql"

	"github.com/cbringf/ml_proxy/main"

	_ "github.com/go-sql-driver/mysql"
)

type ItemService struct {
	DB *sql.DB
}

func (s *main.ItemService) Item(id string) (*main.Item, error) {
	var item main.Item
	var children = make(main.Child)
	var query = `
		SELECT i.*, c.id as cid, c.item_id, c.stop_time as cstop_time
		FROM item as i
		LEFT JOIN child as c ON i.id = c.item_id
		WHERE i.id = $1
	`

	rows, err := s.DB.Query(query, id)
	if err != nil {
		return nil, err
	} else {
		for rows.Next() {
			var child main.Child
			rows.Scan(_, _, _, _, _, _, _, &child.id, _, &child.stopTime)
			children = append(children, child)
		}
		item.children = children
	}
	return nil, item
}
