package types

import (
	"bytes"
	"io"
)

var (
	formatMappings map[string]string
)

// FromCommonDateFormat provides simple common date format recognizer
// and convert it into internal Go `time` standard package
// e.g. MM/DD/YYYY -> 01/02/2006
func FromCommonDateTimeFormat(f string) string {
	if len(f) == 0 {
		return f
	}
	result := bytes.Buffer{}
	tmp := bytes.Buffer{}
	input := bytes.Buffer{}
	input.WriteString(f)
scanning:
	for {
		c, err := input.ReadByte()
		if err != nil {
			break
		}

		// common js date time format characters
		// simple rules:
		// on valid character, do collecting until encounter different character
		switch string(c) {
		case "a", "A", "h", "H", "m", "s", "Y", "M", "D":
			// when match valid character, store it and start collecting
			tmp.WriteByte(c)
			for {
				cc, err := input.ReadByte()
				if err != nil && err == io.EOF {
					m, ok := formatMappings[tmp.String()]
					if ok {
						result.WriteString(m)
					}
					break
				}
				// found different character
				if cc != c {
					input.UnreadByte()
					m, ok := formatMappings[tmp.String()]
					if ok {
						result.WriteString(m)
					}
					tmp.Reset()
					continue scanning
				} else {
					tmp.WriteByte(cc)
				}
			}
		default:
			// fallback just store as-is
			result.WriteByte(c)
		}
	}
	return result.String()
}

func init() {
	formatMappings = map[string]string{
		"a":    "pm",
		"A":    "PM",
		"h":    "3",
		"hh":   "03",
		"HH":   "15",
		"m":    "4",
		"mm":   "04",
		"s":    "5",
		"ss":   "05",
		"YY":   "06",
		"YYYY": "2006",
		"M":    "1",
		"MM":   "01",
		"MMM":  "Jan",
		"MMMM": "January",
		"D":    "2",
		"DD":   "02",
		"DDD":  "Mon",
		"DDDD": "Monday",
	}
}
