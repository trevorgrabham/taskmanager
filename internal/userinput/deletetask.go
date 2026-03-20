package userinput

import (
	"bufio"
	"database/sql"
	"fmt"
	sqlite "local/taskmanager2.0/internal/db"
	"local/taskmanager2.0/internal/task"
	"os"
	"slices"
	"strings"
)

func DeleteTask(db *sql.DB, category string) error {
	if db == nil {
		return fmt.Errorf("deleting task: cannot delete from a nil database")
	}

	queryParams := sqlite.TaskQueryParams{WhichTasks: sqlite.IncTasks, Category: category}
	tasks, err := sqlite.QueryTasks(db, queryParams)
	if err != nil {
		return fmt.Errorf("deleting task: %s", err)
	}

removeLoop:
	for {
		var taskToRemove task.Task
		taskToRemove, err = GetTaskSelection("Which task would you like to remove?", tasks)
		if err != nil {
			return fmt.Errorf("deleting task: %s", err)
		}

		if taskToRemove.IsZero() {
			break removeLoop
		}

		err = sqlite.DeleteTask(db, taskToRemove)
		if err != nil {
			return fmt.Errorf("deleting task: %s", err)
		}

		tasks = slices.DeleteFunc(tasks, func(t task.Task) bool { return t.ID == taskToRemove.ID })
		if len(tasks) < 1 {
			break removeLoop
		}

		fmt.Print("Anything else to remove? (y/n)\t")
		var (
			more   string
			reader = bufio.NewReader(os.Stdin)
		)
		more, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("deleting task: %s", err)
		}

		more = strings.TrimSpace(strings.ToLower(more))
		switch more {
		case "n", "no":
			break removeLoop
		default:
			fmt.Println()
			fmt.Println()
		}
	}
	return nil
}
