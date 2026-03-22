package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
)

func DeleteTask(db *sql.DB, taskToDelete task.Task) error {
	if db == nil {
		return fmt.Errorf("deleting task: cannot complete task for nil database")
	}
	var defaultTask task.Task
	if taskToDelete == defaultTask {
		return fmt.Errorf("deleting task: cannot delete an emtpy task")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("deleting task: %s", err)
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.Exec(`DELETE FROM recurring WHERE task_id = ?`, taskToDelete.ID)
	if err != nil {
		return fmt.Errorf("deleting task: %s", err)
	}

	_, err = tx.Exec(`DELETE FROM task WHERE id = ?`, taskToDelete.ID)
	if err != nil {
		return fmt.Errorf("deleting task: %s", err)
	}

	_ = tx.Commit()

	return nil
}
