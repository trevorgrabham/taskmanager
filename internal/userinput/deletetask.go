package userinput

import (
	"fmt"
	"local/taskmanager2.0/internal/task"
	sqlite "local/taskmanager2.0/internal/db"
	"log"
	"database/sql"
	"slices"
)

func DeleteTask(db *sql.DB) error {
	if db == nil { return fmt.Errorf("deleting task: cannot delete from a nil database") }

	tasks, err := sqlite.QueryIncompleteTasks(db)
	if err != nil { return fmt.Errorf("deleting task: %s", err) }

removeLoop:
	for {
		var id int
		id, err = GetTaskSelection("Which task would you like to remove?", tasks)
		if err != nil { return fmt.Errorf("deleting task: %s" , err) }

		if id == -1 { break removeLoop }

		err = sqlite.DeleteTask(db, id)
		if err != nil { return fmt.Errorf("deleting task: %s" , err) }

		tasks = slices.DeleteFunc(tasks, func(t *task.Task) bool { return t.ID == id })

		fmt.Print("Anything else to remove? (y/n)\t")
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
			break removeLoop
		}
	}
	return nil
}
