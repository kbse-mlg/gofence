package models

import (
	"github.com/kbse-mlg/gofence/utility/types"
)

// MoveHistory Move history
type MoveHistory struct {
	HistoryID     int            `json:"id"`
	ObjectID      int            `json:"object_id"`
	Long          float64        `json:"long"`
	Lat           float64        `json:"lat"`
	Created       types.DateTime `json:"created"`
	StringCreated string         `json:"-" db:"-"`

	// Transient
	Object *Object
}

// AreaCollection data holder for multiple area collection in pagination view
type MoveHistoriesCollection struct {
	MoveHistories []*MoveHistory
	CurrentSearch string
	Size          int64
	Page          int64
	NextPage      int64
}
