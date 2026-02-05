package userinput

import (
	"fmt"
	"local/taskmanager2.0/internal/task"
	sqlite "local/taskmanager2.0/internal/db"
	"database/sql"
	"log"
)

func PushTask(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("pushing task: cannot complete task for nil database")
	}
	tasks, err := sqlite.QueryIncompleteTasks(db)
	if err != nil {
		return fmt.Errorf("pushing task: %s", err)
	}

pushLoop:
	for {
		var id, numDays int
		id, err = GetTaskSelection("Which task would you like to push?", tasks)
		if err != nil { return err }
	
		fmt.Print("How many days should it be pushed back?\t")
		_, err = fmt.Scan(&numDays)
		if err != nil {
			return fmt.Errorf("pushing task: %s", err) 
		}
	
		if numDays < 1 {
			return fmt.Errorf("pushing task: cannot push back by a negative number of days")
		}

		err = sqlite.PushTask(db, id, numDays)
		if err != nil { return fmt.Errorf("pushing task: %s", err) }

		var updatedTask task.Task
		updatedTask, err = sqlite.QueryTaskByID(db, id)
		if err != nil { return fmt.Errorf("pushing task: %s", err) }

		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].DueDate = updatedTask.DueDate
				tasks[i].Type = updatedTask.Type
				break
			}
		}
	
		fmt.Print("Anything else to push? (y/n)\t")
		var more string
		_, err = fmt.Scan(&more)
		if err != nil {
			log.Fatal(err)
		}
		
		switch more {
		case "y", "Y", "YES", "Yes", "yes":
			fmt.Println()
			fmt.Println()
		default:
			break pushLoop
		}
	}
	return nil
}
