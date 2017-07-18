package models

import (
	"fmt"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type Area struct {
	AreaID   int    `json:"id"`
	Name     string `json:"name"`
	Geodata  string `json:"geodata"`
	Type     int    `json:"type"`
	Group    string `json:"group"`
	Active   bool   `json:"active"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
}

type Profile struct {
	Color  string
	Active bool
}

func (u *Area) String() string {
	return fmt.Sprintf("Object(%s)", u.Name)
}

func (area *Area) PreInsert(_ gorp.SqlExecutor) error {
	area.Created = time.Now().UnixNano()
	area.Modified = area.Created
	return nil
}

func (area *Area) PreUpdate(_ gorp.SqlExecutor) error {
	area.Modified = time.Now().UnixNano()
	return nil
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

// AreaCollection data holder for multiple area collection in pagination view
type AreaCollection struct {
	Areas         []*Area
	CurrentSearch string
	Size          int64
	Page          int64
	NextPage      int64
}
