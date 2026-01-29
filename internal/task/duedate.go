package task

import (
	"encoding/json"
	"fmt"
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
		return fmt.Errorf("unmarshaling date: %s", err)
	}

	if dateString == "" {
		*t = TaskDueDate{}
		return nil
	}

	timeZone, err := time.LoadLocation("Local")
	if err != nil {
		return fmt.Errorf("loading timezone: %s", err)
	}

	dueDate, err := time.ParseInLocation(DueDateFormatString, dateString, timeZone)
	if err != nil {
		return fmt.Errorf("parsing DueDate: %s", err)
	}

	*t = TaskDueDate(dueDate)
	return nil
}
