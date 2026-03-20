package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager2.0/internal/task"
)

func PushTask(db *sql.DB, pushedTask task.Task) error {
	if db == nil {
		return fmt.Errorf("pushing task: cannot push task for nil database")
	}

	if pushedTask.IsZero() {
		return fmt.Errorf("pushing task: cannot push an empty task")
	}

	_, err := db.Exec(`UPDATE task SET due_date = ? WHERE id = ?;`, pushedTask.DueDate.Unix(), pushedTask.ID)
	if err != nil {
		return fmt.Errorf("pushing task: %s", err)
	}

	return nil
}
