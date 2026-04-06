package tui

import (
	"bufio"
	"database/sql"
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"os"
	"slices"
	"strings"
)

func CompleteTask(db *sql.DB, category string) error {
	if db == nil {
		return fmt.Errorf("completing task: cannot complete task for nil database")
	}
	queryParams := sqlite.TaskQueryParams{WhichTasks: sqlite.IncTasks, Category: category}
	tasks, err := sqlite.QueryTasks(db, queryParams)
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	for {
		var taskToComplete task.Task
		taskToComplete, err = GetTaskSelection("Which task would you like to complete?", tasks)
		if err != nil {
			return fmt.Errorf("completing task: %s", err)
		}

		if taskToComplete.IsZero() {
			return nil
		}

		err = sqlite.CompleteTask(db, taskToComplete)
		if err != nil {
			return fmt.Errorf("completing task: %s", err)
		}

		if len(tasks) < 2 {
			return nil
		}

		tasks = slices.DeleteFunc(tasks, func(t task.Task) bool { return t.ID == taskToComplete.ID })

		fmt.Print("Anything else to complete? (y/n)\t")
		var (
			more   string
			reader = bufio.NewReader(os.Stdin)
		)
		more, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("completing task: %s", err)
		}

		more = strings.TrimSpace(strings.ToLower(more))
		switch more {
		case "y", "yes":
			fmt.Println()
			fmt.Println()
		default:
			return nil
		}
	}
}
