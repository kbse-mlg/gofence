package models

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/kbse-mlg/gofence/utility/types"
	"github.com/revel/revel"
)

type Object struct {
	ObjectID int     `json:"id"`
	Group    string  `json:"group"`
	Name     string  `json:"name"`
	Long     float64 `json:"long"`
	Lat      float64 `json:"lat"`
	Type     int     `json:"type"`
	Created  int64   `json:"created"`
	Modified int64   `json:"modified"`
}

func (u *Object) String() string {
	return fmt.Sprintf("Object(%s)", u.Name)
}

var nameRegex = regexp.MustCompile("^\\w*$")

// SameLoc check same location
func (obj *Object) SameLoc(newObj *Object) bool {
	return obj.Long == newObj.Long && obj.Lat == newObj.Lat
}

func (obj *Object) PostUpdate(exe gorp.SqlExecutor) error {
	history := MoveHistory{
		ObjectID: obj.ObjectID,
		Long:     obj.Long,
		Lat:      obj.Lat,
		Created:  types.DateTime{Int64: time.Now().UnixNano(), Valid: true},
	}
	return exe.Insert(&history)
}

func (obj *Object) PreInsert(_ gorp.SqlExecutor) error {
	obj.Created = time.Now().UnixNano()
	obj.Modified = obj.Created
	return nil
}

func (obj *Object) PreUpdate(_ gorp.SqlExecutor) error {
	obj.Modified = time.Now().UnixNano()
	return nil
}

func (obj *Object) Validate(v *revel.Validation) {
	v.Check(obj.Name,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
		revel.Match{nameRegex},
	)
	v.Check(obj.Group,
		revel.Required{},
		revel.MaxSize{100},
		revel.MinSize{4},
	)
}

type ObjectCollection struct {
	Objects       []*Object
	CurrentSearch string
	Size          int64
	Page          int64
	NextPage      int64
}
