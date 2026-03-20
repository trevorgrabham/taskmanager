package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager2.0/internal/task"
)

func DeleteTask(db *sql.DB, taskToDelete task.Task) error {
	if db == nil {
		return fmt.Errorf("deleting task: cannot complete task for nil database")
	}
	var defaultTask task.Task
	if taskToDelete == defaultTask {
		return fmt.Errorf("deleting task: cannot delete an emtpy task")
	}

	_, err := db.Exec(`DELETE FROM task WHERE id = ?`, taskToDelete.ID)
	if err != nil {
		return fmt.Errorf("deleting task: %s", err)
	}

	return nil
}
