package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
)

func UpdateTask(db *sql.DB, t task.Task) error {
	if db == nil {
		return fmt.Errorf("updating task: cannot update for a nil database")
	}
	if t.IsZero() {
		return fmt.Errorf("updating task: cannot update without a task")
	}

	_, err := db.Exec(`UPDATE task SET title = ?, category = ?, description = ?, due_date = ? WHERE id = ?`, t.Title, t.Category, t.Description, t.DueDate.Unix(), t.ID)
	if err != nil {
		return fmt.Errorf("updating task: %s", err)
	}

	err = AddRecurringPeriod(db, t)
	if err != nil {
		return fmt.Errorf("updating task: %s", err)
	}

	return nil
}
