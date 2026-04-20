package db

import (
	"database/sql"
	"local/taskmanager/internal/task"
)

func DeleteTask(db *sql.DB, taskToDelete task.Task) error {
	var (
		err error
	)
	if err = checkDBConnection(db); err != nil {
		return err
	}
	if err = checkTaskNotEmpty(taskToDelete); err != nil {
		return err
	}

	return deleteTask(db, taskToDelete.ID)
}
