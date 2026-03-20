package userinput

import (
	"bufio"
	"database/sql"
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"os"
	"strings"
)

func PushTask(db *sql.DB, category string) error {
	if db == nil {
		return fmt.Errorf("pushing task: cannot complete task for nil database")
	}

	queryParams := sqlite.TaskQueryParams{WhichTasks: sqlite.IncTasks, Category: category}
	tasks, err := sqlite.QueryTasks(db, queryParams)
	if err != nil {
		return fmt.Errorf("pushing task: %s", err)
	}

pushLoop:
	for {
		var taskToPush task.Task
		taskToPush, err = GetTaskSelection("Which task would you like to push?", tasks)
		if err != nil {
			return err
		}
		if taskToPush.IsZero() {
			return nil
		}

		var dueDate task.TaskDueDate
		dueDate, err = promptDueDate("New due date:\t")
		if err != nil {
			return fmt.Errorf("pushing task: %s", err)
		}

		taskToPush.DueDate = dueDate
		err = sqlite.PushTask(db, taskToPush)
		if err != nil {
			return fmt.Errorf("pushing task: %s", err)
		}

		fmt.Print("Anything else to push? (y/n)\t")
		var (
			more   string
			reader = bufio.NewReader(os.Stdin)
		)
		more, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("pushing task: %s", err)
		}

		more = strings.TrimSpace(strings.ToLower(more))
		switch more {
		case "n", "no":
			break pushLoop
		default:
			fmt.Println()
			fmt.Println()
		}
	}
	return nil
}
