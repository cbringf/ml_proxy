package dom

import "time"

// ItemChild represents a child of ML Item
type ItemChild struct {
	ID       string    `json:"id" db:"id"`
	ItemID   string    `json:"item_id" db:"item_id"`
	StopTime time.Time `json:"stop_time" db:"stop_time"`
}
