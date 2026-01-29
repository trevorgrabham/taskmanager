package db

import (
	"database/sql"
	"fmt"
	"time"
)

func CompleteTask(db *sql.DB, id int) error {
	if db == nil {
		return fmt.Errorf("completing task: cannot complete task for nil database")
	}
	if id < 0 {
		return fmt.Errorf("completing task: cannot complete a task with id %d", id)
	}

	_, err := db.Exec(`UPDATE tasks SET done = 1, completion_date = ? WHERE id = ?;`, time.Now().Unix(), id)
	if err != nil {
		return fmt.Errorf("updating task db: %s", err)
	}

	return nil
}
