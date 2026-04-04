package tui

import (
	"bufio"
	"database/sql"
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/parse"
	"local/taskmanager/internal/task"
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

		var id int
		id, err = sqlite.AddTask(db, task.Task{Title: title, Category: category, Description: description, DueDate: dueDate})
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		var recurring, period string
		fmt.Print("Recurring?\t")
		recurring, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		recurring = strings.TrimSpace(strings.ToLower(recurring))
		switch recurring {
		case "y", "yes":
			fmt.Print("How often?\t")
			recurring, err = reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}

			period, err = parse.ParsePeriod(recurring)
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}

			err = sqlite.AddRecurringTask(db, task.RecurringTask{TaskID: id, Period: period})
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}
		}

		fmt.Print("Anything else to add? (y/n)\t")
		var more string
		more, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		more = strings.TrimSpace(strings.ToLower(more))
		switch more {
		case "y", "yes":
			fmt.Println()
			fmt.Println()
		default:
			break addLoop
		}
	}
	return nil
}
