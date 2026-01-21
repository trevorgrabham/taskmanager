package userinput

import (
	"bufio"
	"fmt"
	"local/taskmanager2.0/internal/task"
	"os"
	"strings"
	"time"
)

func AddTaskMenu() (task.Task, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("title:\t")
	title, err := reader.ReadString('\n')
	title = strings.TrimSpace(title)
	if err != nil {
		return task.Task{}, err
	}

	fmt.Printf("category:\t")
	category, err := reader.ReadString('\n')
	category = strings.TrimSpace(category)
	if err != nil {
		return task.Task{}, err
	}

	fmt.Printf("description:\t")
	description, err := reader.ReadString('\n')
	description = strings.TrimSpace(description)
	if err != nil {
		return task.Task{}, err
	}

	fmt.Printf("due date (\"now\" or form of DD/MM [XX:XX XM]):\t")
	dueDateString, err := reader.ReadString('\n')
	dueDateString = strings.TrimSpace(dueDateString)
	if dueDateString == "" || dueDateString == "now" || dueDateString == "NOW" || dueDateString == "Now" {
		return task.Task{Title: title, Category: category, Description: description, DueDate: task.TaskDueDate(time.Now())}, nil
	}
	if err != nil {
		return task.Task{}, err
	}

	dueDate, err := time.Parse("02/01", dueDateString)
	if err != nil {
		dueDate, err = time.Parse("02/01 3:04 PM", dueDateString)
		if err != nil {
			return task.Task{}, err
		}
	}

	now := time.Now()
	year := now.Year()
	if dueDate.Month() < now.Month() || (dueDate.Month() == now.Month() && dueDate.Day() < now.Day()) {
		year++
	}
	return task.Task{Title: title, Category: category, Description: description, DueDate: task.TaskDueDate(time.Date(year, dueDate.Month(), dueDate.Day(), dueDate.Hour(), dueDate.Minute(), 0, 0, now.Location()))}, nil
}
