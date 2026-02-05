package userinput

import (
	"database/sql"
	"fmt"
	sqlite "local/taskmanager2.0/internal/db"
	"local/taskmanager2.0/internal/task"
	"slices"
)

func CompleteTask(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("completing task: cannot complete task for nil database")
	}
	tasks, err := sqlite.QueryIncompleteTasks(db)
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

completeLoop:
	for {
		var id int
		id, err = GetTaskSelection("Which task would you like to complete?", tasks)
		if err != nil {
			return fmt.Errorf("completing task: %s", err)
		}

		if id == -1 {
			break completeLoop
		}

		err = sqlite.CompleteTask(db, id)
		if err != nil {
			return fmt.Errorf("completing task: %s", err)
		}

		tasks = slices.DeleteFunc(tasks, func(t *task.Task) bool { return t.ID == id })

		fmt.Print("Anything else to complete? (y/n)\t")
		var more string
		_, err = fmt.Scan(&more)
		if err != nil {
			return fmt.Errorf("completing task: %s", err)
		}

		switch more {
		case "y", "Y", "YES", "Yes", "yes":
			fmt.Println()
			fmt.Println()
		default:
			break completeLoop
		}
	}
	return nil
}
