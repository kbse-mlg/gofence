package models

// MoveHistory Move history
type MoveHistory struct {
	HistoryID int     `json:"id"`
	ObjectID  int     `json:"object_id"`
	Long      float64 `json:"long"`
	Lat       float64 `json:"lat"`
	Created   int64   `json:"created"`

	// Transient
	Object *Object
}
