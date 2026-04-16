package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
)

func UpdateDueDate(db *sql.DB, taskToUpdate task.Task) error {
	if db == nil {
		return fmt.Errorf("updating due date: cannot update a nil database")
	}
	if taskToUpdate.IsZero() {
		return fmt.Errorf("updating due date: cannot update without a task")
	}

	var err error
	if taskToUpdate.DueDate.IsZero() {
		_, err = db.Exec(`UPDATE task SET due_date = NULL WHERE id = ?`, taskToUpdate.ID)
	} else {
		_, err = db.Exec(`UPDATE task SET due_date = ? WHERE id = ?`, taskToUpdate.DueDate.Unix(), taskToUpdate.ID)
	}
	if err != nil {
		return fmt.Errorf("updating due date: %s", err)
	}

	return nil
}
