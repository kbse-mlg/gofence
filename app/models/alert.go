package models

import (
	"fmt"

	"github.com/kbse-mlg/gofence/utility/types"
)

const (
	Info    = 1
	Warning = 2
	Danger  = 3
	Success = 4
)

type Alert struct {
	ID      string         `db:"id" json:"id"`
	Level   int            `db:"level" json:"level"`
	Info    string         `db:"info" json:"info"`
	Created types.DateTime `db:"created" json:"created"`
}

func (u *Alert) String() string {
	return fmt.Sprintf("Alert(%s)", u.Level)
}

func NewAlert(level int, info string) Alert {
	return Alert{Level: level, Info: info}
}
