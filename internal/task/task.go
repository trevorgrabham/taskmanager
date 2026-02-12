package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var SaveFileName = "/home/trevorgrabham/.config/taskmanager/tasks.json"

type WhichTasks int

const (
	All WhichTasks = iota
	Inc
	Comp
)

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
func MarshalTasks(out *os.File, tasks ...*Task) error {
	if out == nil {
		return errors.New("cannot marshal to a nil out file")
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(out.Name())

	tmpFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}

	tmpFileName := tmpFile.Name()
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFileName)
	}()

	_, err = tmpFile.Write(data)
	if err != nil {
		return err
	}
	err = tmpFile.Sync()
	if err != nil {
		return err
	}

	stat, err := tmpFile.Stat()
	if err != nil {
		return err
	}
	if stat.Size() != int64(len(data)) {
		return fmt.Errorf("incomplete write: expected %d bytes, wrote %d", len(data), stat.Size())
	}

	err = tmpFile.Close()
	if err != nil {
		return err
	}

	err = os.Chmod(tmpFileName, 0644)
	if err != nil {
		return err
	}

	err = os.Rename(tmpFileName, out.Name())
	if err != nil {
		return err
	}

	dirFd, err := os.Open(dir)
	if err != nil {
		return err
	}

	defer dirFd.Close()

	return dirFd.Sync()
}

func UnmarshalTasks(saveFile string) (tasks TaskList, err error) {
	if _, err = os.Stat(saveFile); err != nil {
		return TaskList{}, nil
	}
	var data []byte
	data, err = os.ReadFile(saveFile)
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
		if time.Now().After(time.Time(tasks[i].DueDate)) {
			tasks[i].Type = Due
			continue
		}
		tasks[i].Type = Upcoming
	}
	return tasks, nil
}
