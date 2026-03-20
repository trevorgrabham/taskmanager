package userinput

import (
	"bufio"
	"database/sql"
	"fmt"
	sqlite "local/taskmanager2.0/internal/db"
	"local/taskmanager2.0/internal/task"
	"os"
	"strings"
)

func AddTask(db *sql.DB, category string) error {
	if db == nil {
		return fmt.Errorf("adding task: cannot add a task into a nil database")
	}

	reader := bufio.NewReader(os.Stdin)
addLoop:
	for {
		fmt.Printf("title:\t")
		title, err := reader.ReadString('\n')
		title = strings.TrimSpace(title)
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}
		if title == "" {
			return fmt.Errorf("adding task: task needs a title")
		}

		if category == "" {
			fmt.Printf("category:\t")
			category, err = reader.ReadString('\n')
			category = strings.TrimSpace(category)
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}
		}

		fmt.Printf("description:\t")
		description, err := reader.ReadString('\n')
		description = strings.TrimSpace(description)
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		var dueDate task.TaskDueDate
		dueDate, err = promptDueDate("due date (day of the week (Mon/Monday), or form of DD/MM [XX:XX XM]):\t")
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}
		err = sqlite.AddTask(db, task.Task{Title: title, Category: category, Description: description, DueDate: dueDate})
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		fmt.Print("Anything else to add? (y/n)\t")
		var more string
		more, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		more = strings.TrimSpace(strings.ToLower(more))
		switch more {
		case "n", "no":
			break addLoop
		default:
			fmt.Println()
			fmt.Println()
		}
	}
	return nil
}
