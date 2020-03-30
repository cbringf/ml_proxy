package mysql

import (
	"database/sql"
	"fmt"

	"github.com/cbringf/proxy/dom"

	_ "github.com/go-sql-driver/mysql"
)

// ItemService is a MySql implementation of ItemService
type ItemService struct {
	DB *sql.DB
}

// Item resolves an Item by ID from local Db
func (s ItemService) Item(id string) (*dom.Item, *dom.Error) {
	var item dom.Item
	var children = make([]*dom.ItemChild, 0)
	var queryItem = `
		SELECT id, category_id, title, price, currency_id, start_time, stop_time FROM item WHERE id = ?
	`
	var queryItemChild = `
		SELECT id, item_id, stop_time FROM child WHERE item_id = ?
	`
	err := s.DB.QueryRow(queryItem, id).Scan(&item.ID, &item.CategoryID, &item.Title, &item.Price, &item.CurrencyID, &item.StartTime, &item.StopTime)

	if err == sql.ErrNoRows {
		return nil, &dom.Error{
			Status:  404,
			Error:   "not_found",
			Message: fmt.Sprintf("Item with ID %s not found", id),
		}
	}

	rows, err := s.DB.Query(queryItemChild, id)

	if err != nil {
		return nil, dom.UnknownError()
	}

	defer rows.Close()

	for rows.Next() {
		var aux dom.ItemChild

		rows.Scan(&aux.ID, &aux.ItemID, &aux.StopTime)

		children = append(children, &aux)
	}
	item.Children = children

	return &item, nil
}

// Write writes remote resolved Item to local Db
func (s ItemService) Write(item *dom.Item) *dom.Error {
	stmtIn, err := s.DB.Prepare("INSERT INTO item (id, category_id, title, price, currency_id, start_time, stop_time) VALUES (?,?,?,?,?,?,?)")

	if err != nil {
		return dom.UnknownError()
	}

	defer stmtIn.Close()

	_, err = stmtIn.Exec(item.ID, item.CategoryID, item.Title, item.Price, item.CurrencyID, item.StartTime, item.StopTime)

	if err != nil {
		return dom.UnknownError()
	}

	stmtInCl, err := s.DB.Prepare("INSERT INTO child (id, item_id, stop_time) VALUES (?,?,?)")

	if err != nil {
		// TODO Remove Item from DB
		return dom.UnknownError()
	}

	defer stmtInCl.Close()

	for _, c := range item.Children {
		child := dom.ItemChild{
			ID:       c.ID,
			ItemID:   item.ID,
			StopTime: c.StopTime,
		}
		_, err := stmtInCl.Exec(child.ID, child.ItemID, child.StopTime)

		if err != nil {
			// TODO Remove Item from DB and All inserted Children
			return dom.UnknownError()
		}
	}
	return nil
}
