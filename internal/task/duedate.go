package task

import (
	"encoding/json"
	"time"
)

var DueDateFormatString = "Mon Jan 2 2006 3:04 PM"

type TaskDueDate time.Time

func (t TaskDueDate) String() string {
	if t == (TaskDueDate{}) {
		return ""
	}
	return time.Time(t).Format(DueDateFormatString)
}

func (t TaskDueDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TaskDueDate) UnmarshalJSON(data []byte) error {
	var dateString string
	err := json.Unmarshal(data, &dateString)
	if err != nil {
		return err
	}

	if dateString == "" {
		*t = TaskDueDate{}
		return nil
	}

	dueDate, err := time.Parse(DueDateFormatString, dateString)
	if err != nil {
		return err
	}

	*t = TaskDueDate(dueDate)
	return nil
}
