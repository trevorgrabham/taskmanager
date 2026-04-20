package db

import (
	"database/sql"
	"local/taskmanager/internal/task"
)

func CompleteTask(db *sql.DB, taskToComplete task.Task) error {
	var (
		err                            error
		tx                             *sql.Tx
		t task.Task
		columns, placeholders []string 
		args []any
	)
	if err = checkDBConnection(db); err != nil {
		return err
	}

	if err = checkTaskNotEmpty(taskToComplete); err != nil {
		return err
	}

	if tx, err = startDBSession(db); err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if t, err = completeTask(tx, taskToComplete.ID); err != nil {
		return err
	}

	if t.RecurringPeriod == "" {
		_ = tx.Commit()
		return nil
	}

	if t.DueDate, err = computeNextDueDate(t.RecurringPeriod); err != nil { return err }

	if columns, placeholders, args, err = setupQueryParams(db, t); err != nil { return err }

	if _, err = addNewTask(db, columns, placeholders, args); err != nil { return err }

	_ = tx.Commit()

	return nil
}
