package parse

import (
	"fmt"
	"strings"
	"time"
)

func ParseDayOfWeek(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("parsing day of week: cannot parse an empty string")
	}

	s = strings.ToLower(s)
	now := time.Now()
	var dayDiff int
	switch s {
	case "sunday", "sun", "su":
		dayDiff = int(time.Sunday-now.Weekday()+7) % 7
	case "monday", "mon", "m":
		dayDiff = int(time.Monday-now.Weekday()+7) % 7
	case "tuesday", "tues", "tue", "tu":
		dayDiff = int(time.Tuesday-now.Weekday()+7) % 7
	case "wednesday", "wed", "w":
		dayDiff = int(time.Wednesday-now.Weekday()+7) % 7
	case "thursday", "thurs", "thur", "th":
		dayDiff = int(time.Thursday-now.Weekday()+7) % 7
	case "friday", "fri", "f":
		dayDiff = int(time.Friday-now.Weekday()+7) % 7
	case "saturday", "sat", "sa":
		dayDiff = int(time.Saturday-now.Weekday()+7) % 7
	default:
		return time.Time{}, fmt.Errorf("parsing day of week: %s is not a recognized day of the week", s)
	}
	if dayDiff == 0 {
		dayDiff = 7
	}

	return time.Date(now.Year(), now.Month(), now.Day()+dayDiff, 23, 59, 0, 0, now.Location()), nil
}
