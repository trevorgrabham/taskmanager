package db

import (
	"database/sql"
	"fmt"
)

func DeleteTask(db *sql.DB, id int) error {
	if db == nil {
		return fmt.Errorf("completing task: cannot complete task for nil database")
	}
	if id < 0 {
		return fmt.Errorf("completing task: cannot complete a task with id %d", id)
	}

	_, err := db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, taskTableName), id)
	if err != nil {
		return fmt.Errorf("deleting from task db: %s", err)
	}

	return nil
}
