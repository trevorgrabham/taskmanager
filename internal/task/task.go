package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID              int
	Title           string
	Category        string
	Description     string
	DueDate         TaskDueDate
	CompletionDate  TaskDueDate
	Done            bool
	RecurringPeriod string // n days || n weeks || n months
	UserID          int
}

func (t Task) ANSICode() string {
	null := time.Time{}
	// complete
	if !time.Time(t.CompletionDate).Equal(null) {
		return "32m"
	}
	// no due date
	if time.Time(t.DueDate).Equal(null) {
		return "0m"
	}
	// upcoming
	if time.Time(t.DueDate).After(time.Now()) {
		return "33m"
	}
	// due
	return "31m"
}

func (t Task) String() string {
	formattedString := fmt.Sprintf("\033[30;46m%s\033[0m", t.Title)
	if t.Category != "" {
		formattedString = fmt.Sprintf("%s [\033[95m%s\033[0m]", formattedString, t.Category)
	}
	if !t.DueDate.IsZero() {
		formattedString = fmt.Sprintf("%s - \033[%s%s\033[0m", formattedString, t.ANSICode(), t.DueDate.String())
	}
	if t.Description != "" {
		formattedString = fmt.Sprintf("%s\n%s", formattedString, t.Description)
	}
	if t.Done {
		formattedString = fmt.Sprintf("%s\n\033[32mCompleted on %s\033[0m", formattedString, t.CompletionDate.String())
	}

	return formattedString
}

func (t Task) IsZero() bool {
	return t.ID == 0 && t.Title == "" && t.Category == "" && t.Description == "" && t.DueDate.IsZero() && t.CompletionDate.IsZero() && !t.Done && t.RecurringPeriod == "" && t.UserID == 0
}

func (t Task) PeriodValue() (int, error) {
	valueString := strings.Split(t.RecurringPeriod, " ")[0]
	value, err := strconv.Atoi(valueString)
	if err != nil {
		return -1, fmt.Errorf("task.PeriodValue(): %s", err)
	}

	return value, nil
}

func (t Task) PeriodUnit() string {
	return strings.Split(t.RecurringPeriod, " ")[1]
}

// ================================================== TaskList ==================================================

type TaskList []Task

func (tl TaskList) String() string {
	taskStrings := make([]string, len(tl))
	for i := range tl {
		taskStrings[i] = tl[i].String()
	}
	return strings.Join(taskStrings, "\n\n")
}

// ================================================== TaskDueDate ==================================================

var DueDateFormatString = "Mon Jan 2 2006 3:04 PM"

type TaskDueDate time.Time

func (t TaskDueDate) String() string {
	if t == (TaskDueDate{}) {
		return ""
	}
	return time.Time(t).Format(DueDateFormatString)
}

func (t TaskDueDate) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t TaskDueDate) Unix() int64 {
	return time.Time(t).Unix()
}

func (t TaskDueDate) IsDue() bool {
	return !time.Now().After(time.Time(t))
}
