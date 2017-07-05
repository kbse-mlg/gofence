package models

import (
	"fmt"

	"github.com/revel/revel"
)

type Area struct {
	AreaID  int    `json:"area_id"`
	Name    string `json:"name"`
	Geodata string `json:"geodata"`
	Type    int    `json:"type"`
	Group   string `json:"group"`
}

type Profile struct {
	Color  string
	Active bool
}

func (u *Area) String() string {
	return fmt.Sprintf("Object(%s)", u.Name)
}

func (area *Area) Validate(v *revel.Validation) {
	v.Check(area.Name,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
		revel.Match{nameRegex},
	)

	v.Check(area.Geodata,
		revel.Required{},
	)

	v.Check(area.Group,
		revel.Required{},
	)
}

type AreaCollection struct {
	Areas         []*Area
	CurrentSearch string
	Size          int64
	Page          int64
	NextPage      int64
}
