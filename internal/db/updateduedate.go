package db

import (
	"database/sql"
	"local/taskmanager/internal/task"
)

func UpdateDueDate(db *sql.DB, taskToUpdate task.Task) error {
	var err error
	if err = checkDBConnection(db); err != nil {
		return err
	}
	if err = checkTaskNotEmpty(taskToUpdate); err != nil {
		return err
	}

	return updateDueDate(db, taskToUpdate.ID, taskToUpdate.DueDate.Unix())
}
