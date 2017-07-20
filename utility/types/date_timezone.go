package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

const defaultTimeFormat = "Mon, 01/_2/2006 15:04:05 PM"

type DateTimezone struct {
	Int64          int64
	Valid          bool
	TimezoneOffset int
	Relative       bool
	CommonFormat   string
}

func (n *DateTimezone) Scan(value interface{}) error {
	if value == nil {
		n.Int64, n.Valid = int64(0), false
		return nil
	}
	n.Valid = true
	n.Int64 = value.(int64)
	return nil
}

func (n DateTimezone) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Int64), nil
}

// TODO: format in relative time representation
// e.g. 2 days ago
func (n DateTimezone) formatRelative(t time.Time) string {
	nsec := int64(n.TimezoneOffset * 6E10)
	now := time.Unix(0, time.Now().UTC().UnixNano()-nsec).In(time.FixedZone("", -n.TimezoneOffset))
	d := now.Sub(t)
	sec := int64(d.Seconds())
	days := d.Hours() / 24.0
	months := days / 30.0
	switch {
	case sec < 60:
		if sec <= 1 {
			return "one second ago"
		}
		return fmt.Sprintf("%d seconds ago", sec)
	case sec < 120:
		return "a minute ago"
	case sec < 2700:
		return fmt.Sprintf("%d minutes ago", int64(d.Minutes()))
	case sec < 5400:
		return "an hour ago"
	case sec < 86400:
		return fmt.Sprintf("%d hours ago", int64(d.Hours()))
	case sec < 172800:
		return "yesterday"
	case sec < 2592000:
		return fmt.Sprintf("%d days ago", int64(days))
	case sec < 31104000:
		if months <= 1 {
			return "one month ago"
		}
		return fmt.Sprintf("%d months ago", int64(months))
	}
	years := days / 365.0
	if years <= 1 {
		return "one year ago"
	}
	return fmt.Sprintf("%d years ago", int64(years))
}

func (n DateTimezone) dateTimeFormat() string {
	if n.CommonFormat != "" {
		return FromCommonDateTimeFormat(n.CommonFormat)
	}
	return defaultTimeFormat
}

// MarshalJSON converts the int64 into string representation
// TODO: need many enhancements
func (n DateTimezone) MarshalJSON() ([]byte, error) {
	if n.Int64 == 0 {
		return []byte(strconv.Quote("")), nil
	}
	nsec := int64(n.TimezoneOffset * 6E10)
	t := time.Unix(0, n.Int64-nsec).In(time.FixedZone("", -n.TimezoneOffset))
	if n.Relative {
		x := n.formatRelative(t)
		return []byte(strconv.Quote(x)), nil
	}
	s := strconv.Quote(t.Format(n.dateTimeFormat()))
	return []byte(s), nil
}

// UnmarshalJSON parse the string representation into in64 value
func (n *DateTimezone) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(n.dateTimeFormat(), s)
	if err != nil {
		return err
	}
	n.Int64, n.Valid = t.UTC().UnixNano(), true
	return nil

}
