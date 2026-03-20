package userinput

import (
	"bufio"
	"fmt"
	"local/taskmanager2.0/internal/task"
	"os"
	"strings"
	"time"
)

func parseDayOfWeek(day string) (time.Time, error) {
	if day == "" {
		return time.Time{}, fmt.Errorf("parsing day of week: cannot parse an empty string")
	}

	day = strings.ToLower(day)
	now := time.Now()
	var dayDiff int
	switch day {
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
		return time.Time{}, fmt.Errorf("parsing day of week: %s is not a recognized day of the week", day)
	}
	if dayDiff == 0 {
		dayDiff = 7
	}

	return time.Date(now.Year(), now.Month(), now.Day()+dayDiff, 0, 0, 0, 0, now.Location()), nil
}

func promptDueDate(prompt string) (task.TaskDueDate, error) {
	if prompt != "" {
		fmt.Print(prompt)
	}
	reader := bufio.NewReader(os.Stdin)
	dueDateString, err := reader.ReadString('\n')
	if err != nil {
		return task.TaskDueDate{}, err
	}

	dueDateString = strings.TrimSpace(strings.ToLower(dueDateString))
	if dueDateString == "" || dueDateString == "n" || dueDateString == "no" || dueDateString == "none" {
		return task.TaskDueDate{}, nil
	}

	dueDate, err := time.Parse("02/01", dueDateString)
	if err != nil {
		dueDate, err = time.Parse("02/01 3:04 PM", dueDateString)
		if err != nil {
			dueDate, err = parseDayOfWeek(dueDateString)
			if err != nil {
				return task.TaskDueDate{}, err
			}
		}
	}

	now := time.Now()
	dueDate = time.Date(now.Year(), dueDate.Month(), dueDate.Day(), dueDate.Hour(), dueDate.Minute(), 0, 0, now.Location())
	if dueDate.Before(now) {
		dueDate = dueDate.AddDate(1, 0, 0)
	}
	return task.TaskDueDate(dueDate), nil
}
