package task

import (
	"fmt"
	"strings"
	"time"
)

var SaveFileName = "/home/trevorgrabham/.config/taskmanager/tasks.json"

type Task struct {
	ID             int         `json:"id"`
	Title          string      `json:"title"`
	Category       string      `json:"category"`
	Description    string      `json:"description"`
	DueDate        TaskDueDate `json:"due-date"`
	Done           bool        `json:"done"`
	CompletionDate TaskDueDate `json:"completion-date,omitempty"`
	Type           TaskStatus  `json:"-"`
}

func (t Task) String() string {
	var formattedString string
	if t.Category == "" {
		formattedString = fmt.Sprintf("\033[30;46m%s\033[0m - \033[%s%s\033[0m", t.Title, t.Type.ANSICode(), t.DueDate.String())
	} else {
		formattedString = fmt.Sprintf("\033[30;46m%s\033[0m [\033[95m%s\033[0m] - \033[%s%s\033[0m", t.Title, t.Category, t.Type.ANSICode(), t.DueDate.String())
	}
	if t.Description != "" {
		formattedString = fmt.Sprintf("%s\n%s", formattedString, t.Description)
	}
	if t.Done {
		formattedString = fmt.Sprintf("%s\nCompleted on %s", formattedString, t.CompletionDate.String())
	}

	return formattedString
}

// ================================================== WhichTasks ==================================================

type WhichTasks int

const (
	All WhichTasks = iota
	Inc
	Comp
)

// ================================================== TaskStatus ==================================================

type TaskStatus int

const (
	Due TaskStatus = iota
	Upcoming
	Done
)

func (t TaskStatus) ANSICode() string {
	switch t {
	case Due:
		return "31m"
	case Upcoming:
		return "33m"
	case Done:
		return "32m"
	default:
		return ""
	}
}

func (t TaskStatus) String() string {
	switch t {
	case Due:
		return "Due"
	case Upcoming:
		return "Upcoming"
	case Done:
		return "Done"
	default:
		return ""
	}
}

// ================================================== TaskList ==================================================

type TaskList []*Task

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

// configHome := os.Getenv("XDG_CONFIG_HOME")
// if configHome == "" {
// home, err := os.UserHomeDir()
// if err != nil {
// fmt.Println("Timer: unable to locate users home directory")
// os.Exit(1)
// }
// configHome = filepath.Join(home, ".config")
// }
// dir := filepath.Join(configHome, "timer")
// os.MkdirAll(dir, 0700)
// jobFile := filepath.Join(dir, "job-num")
