package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
)

func TaskByID(db *sql.DB, taskID int) (task.Task, error) {
	var (
		err error
		t   task.Task
	)
	if err = checkDBConnection(db); err != nil {
		return t, err
	}
	if taskID < 1 {
		return t, fmt.Errorf("TaskByID: no id")
	}

	return getTaskByID(db, taskID)
}
