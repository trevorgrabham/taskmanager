package db

import (
	"database/sql"
	"fmt"
)

func PushTask(db *sql.DB, id, numDays int) error {
	if db == nil {
		return fmt.Errorf("pushing task: cannot push task for nil database")
	}
	if id < 0 {
		return fmt.Errorf("pushing task: cannot push a task with id %d", id)
	}

	_, err := db.Exec(fmt.Sprintf(`UPDATE %s SET due_date = due_date + ? WHERE id = ?;`, taskTableName), numDays*24*60*60, id)
	if err != nil {
		return fmt.Errorf("updating task db: %s", err)
	}

	return nil
}
