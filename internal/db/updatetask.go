package db

import (
	"database/sql"
	"local/taskmanager/internal/task"
)

func UpdateTask(db *sql.DB, t task.Task) error {
	var (
		err error
		tx  *sql.Tx
	)
	if err = checkDBConnection(db); err != nil {
		return err
	}
	if err = checkTaskNotEmpty(t); err != nil {
		return err
	}
	if tx, err = startDBSession(db); err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err = updateTask(db, t); err != nil {
		return err
	}

	_ = tx.Commit()

	return nil
}
