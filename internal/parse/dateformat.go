package parse

import (
	"fmt"
	"time"
)

func ParseDate(s string) (time.Time, error) {
	if s == "" { return time.Time{}, fmt.Errorf("parsing date: cannot parse an empty string") }


	date, err := time.Parse("02/01", s)
	if err != nil {
		date, err = time.Parse("02/01/06", s)
		if err != nil {
			return time.Time{}, err
		}
	}

	var year int
	if date.Year() == 0 {
		year = time.Now().Year()
	} else {
		year = date.Year()
	}
	return time.Date(year, date.Month(), date.Day(), 0, 0, 0, 0, date.Location()), nil
}
