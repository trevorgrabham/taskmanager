package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ParsePeriod(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("cannot parse an empty period")
	}

	s = strings.TrimSpace(strings.ToLower(s))

	re := regexp.MustCompile(`^(\d+)\s*([A-Za-z]+)$`)
	matches := re.FindStringSubmatch(s)

	if matches == nil {
		return "", fmt.Errorf("bad period format")
	}

	unit := matches[2]
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}

	var period string
	switch unit {
	case "d", "day", "days":
		period = fmt.Sprintf("%d days", value)
	case "w", "week", "weeks":
		period = fmt.Sprintf("%d weeks", value)
	case "m", "month", "months":
		period = fmt.Sprintf("%d months", value)
	default:
		return "", fmt.Errorf("unrecognized format")
	}

	return period, nil
}
