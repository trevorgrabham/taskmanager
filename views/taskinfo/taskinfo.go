package taskinfo

import "time"

type TaskInfoViewData struct {
	ID              int
	Title           string
	Category        string
	Description     string
	DueDate         time.Time
	RecurringPeriod string
}
