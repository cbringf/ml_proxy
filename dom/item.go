package dom

import "time"

// Item represents relevant information of ML Item
type Item struct {
	ID         string       `json:"id" db:"id"`
	Title      string       `json:"title" db:"title"`
	CategoryID string       `json:"category_id" db:"category_id"`
	Price      float32      `json:"price" db:"price"`
	CurrencyID string       `json:"currency_id" db:"currency_id"`
	StartTime  time.Time    `json:"start_time" db:"start_time"`
	StopTime   time.Time    `json:"stop_time" db:"stop_time"`
	Children   []*ItemChild `json:"children"`
}

// ItemService represents a service to resolve an Item by Id
type ItemService interface {
	Item(id string) (*Item, *Error)
}

// CacheItemService represents a service to locally store relevant information
// of requested ML Item
type CacheItemService interface {
	Write(item *Item) *Error
}
