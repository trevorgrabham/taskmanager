package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"time"
)

func CompleteTask(db *sql.DB, taskToComplete task.Task) error {
	if db == nil {
		return fmt.Errorf("completing task: cannot complete task for nil database")
	}
	var defaultTask task.Task
	if taskToComplete == defaultTask {
		return fmt.Errorf("completing task: cannot complete an empty task")
	}

	_, err := db.Exec(`UPDATE task SET done = 1, completion_date = ? WHERE id = ?;`, time.Now().Unix(), taskToComplete.ID)
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	return nil
}
