package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"time"
)

func ToggleTaskComplete(db *sql.DB, taskToToggle task.Task) (t task.Task, err error) {
	if db == nil {
		return task.Task{}, fmt.Errorf("toggle task: cannot toggle task completion for a nil database")
	}
	if taskToToggle.IsZero() {
		return task.Task{}, fmt.Errorf("toggle task: cannot toggle task for an empty task")
	}

	row := db.QueryRow(`
		UPDATE task 
		SET done = (done + 1) % 2, completion_date = 
			CASE 
				WHEN completion_date IS NULL THEN ? 
				ELSE NULL
			END
		WHERE id = ?
		RETURNING id, title, category, description, due_date, completion_date, done, NULL
	`, time.Now().Unix(), taskToToggle.ID)
	t, err = scanTaskRow(row)
	if err != nil {
		return task.Task{}, fmt.Errorf("toggle task: %s", err)
	}

	return t, nil
}
