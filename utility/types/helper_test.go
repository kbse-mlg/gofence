package types

import (
	"testing"
)

var testDates = []struct {
	js     string
	golang string
}{
	{
		js:     "YYYY/MM/DD",
		golang: "2006/01/02",
	},
	{
		js:     "HH:mm:ss",
		golang: "15:04:05",
	},
	{
		js:     "hh:mm:ss",
		golang: "03:04:05",
	},
	{
		js:     "DDDD, MMMM/DD/YYYY",
		golang: "Monday, January/02/2006",
	},
}

func TestCommonDateFormatConversion(t *testing.T) {
	for _, d := range testDates {
		g := FromCommonDateTimeFormat(d.js)
		if g != d.golang {
			t.Errorf("expect: %s got: %s\n", d.golang, g)
		}
	}
}
