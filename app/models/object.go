package models

import (
	"fmt"
	"regexp"

	"github.com/revel/revel"
)

type Object struct {
	ObjectID int
	Name     string
	Long     float64
	Lat      float64
	Type     int
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
}
