package models

import (
	"fmt"

	"github.com/revel/revel"
)

type Area struct {
	AreaID  int
	Name    string
	Geodata string
	Type    int
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
}
