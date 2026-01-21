package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

var SaveFileName = "test.json"

type TaskType string

const (
	Due      TaskType = "Due"
	Upcoming TaskType = "Upcoming"
	Done     TaskType = "Done"
)

type Task struct {
	Title          string      `json:"title"`
	Category       string      `json:"category"`
	Description    string      `json:"description"`
	DueDate        TaskDueDate `json:"due-date"`
	Done           bool        `json:"done"`
	CompletionDate TaskDueDate `json:"completion-date,omitempty"`
	Type           TaskType    `json:"-"`
}

func (t Task) String() string {
	var status string
	if t.Done {
		status = "Done!"
	} else if time.Time(t.DueDate).After(time.Now()) {
		status = "Upcoming"
	} else {
		status = "Due"
	}
	var formattedString string
	if t.Category == "" {
		formattedString = fmt.Sprintf("%s - %s", t.Title, status)
	} else {
		formattedString = fmt.Sprintf("%s [%s] - %s", t.Title, t.Category, status)
	}
	if t.Description != "" {
		formattedString = fmt.Sprintf("%s\n%s", formattedString, t.Description)
	}
	if t.Done {
		formattedString = fmt.Sprintf("%s\nCompleted on %s", formattedString, t.CompletionDate.String())
	}
	return formattedString
}

func MarshalTasks(out *os.File, tasks ...*Task) error {
	if out == nil {
		return errors.New("cannot marshal to a nil out file")
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	_, err = out.Write(data)
	return err
}

func UnmarshalTasks() (tasks TaskList, err error) {
	if _, err = os.Stat(SaveFileName); err != nil {
		return TaskList{}, nil
	}
	var data []byte
	data, err = os.ReadFile(SaveFileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	for i := range tasks {
		if tasks[i].Done {
			tasks[i].Type = Done
			continue
		}
		if !time.Now().After(time.Time(tasks[i].DueDate)) {
			tasks[i].Type = Due
			continue
		}
		tasks[i].Type = Upcoming
	}
	return tasks, nil
}
