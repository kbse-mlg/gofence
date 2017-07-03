package models

import (
	"fmt"
	"regexp"

	"github.com/revel/revel"
)

type Object struct {
	ObjectID int     `json:"object_id"`
	Group    string  `json:"group"`
	Name     string  `json:"name"`
	Long     float64 `json:"long"`
	Lat      float64 `json:"lat"`
	Type     int     `json:"type"`
}

func (u *Object) String() string {
	return fmt.Sprintf("Object(%s)", u.Name)
}

var nameRegex = regexp.MustCompile("^\\w*$")

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
}
