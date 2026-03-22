package parse

import (
	"fmt"
	"local/taskmanager/internal/task"
	"time"
)

func ParseDateAndTime(s string) (task.TaskDueDate, error) {
	if s == "" {
		return task.TaskDueDate{}, fmt.Errorf("cannot parse an empty string")
	}

	date, err := time.Parse("02/01", s)
	if err != nil {
		date, err = time.Parse("02/01 3:04 pm", s)
		if err != nil {
			date, err = ParseDayOfWeek(s)
			if err != nil {
				return task.TaskDueDate{}, err
			}
		}
	}
	now := time.Now()
	if date.Hour() == 0 && date.Minute() == 0 {
		date = time.Date(now.Year(), date.Month(), date.Day(), 23, 59, 0, 0, now.Location())
	} else {
		date = time.Date(now.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), 0, 0, now.Location())
	}
	if date.Before(now) {
		date = date.AddDate(1, 0, 0)
	}
	return task.TaskDueDate(date), nil
}
