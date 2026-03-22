package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
)

func AddRecurringTask(db *sql.DB, recurringTask task.RecurringTask) error {
	if db == nil {
		return fmt.Errorf("adding recurring task: cannot add into nil database")
	}
	if recurringTask.IsZero() {
		return fmt.Errorf("adding recurring task: cannot add empty task")
	}

	_, err := db.Exec(`INSERT INTO recurring(task_id, period) VALUES (?, ?)`, recurringTask.TaskID, recurringTask.Period)
	if err != nil {
		return fmt.Errorf("adding recurring task: %s", err)
	}

	return nil
}
