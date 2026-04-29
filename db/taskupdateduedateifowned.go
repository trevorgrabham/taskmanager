package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TaskUpdateDueDateIfOwned updates the DueDate for the task identified by task.ID if it is owned by task.UserID. If DueDate was not NULL, then the date changes, but the time stays the same. Returns the updated Task.
//
// If Repo is not initialized, returns an ErrNotConnected.
// If taskID is empty, returns an ErrTaskNotExist.
// If userID is empty, returns an ErrUserNotExist.
// If the task is not owned by task.UserID, returns an ErrNotOwner.
// If the no task exists with id = task.ID, returns an ErrTaskNotExist.
// If an error occurs in the Repo, returns an ErrInternalRepo.
// If an error occurs commiting the Repo transaction, an ErrTransactionCommit is returned.
func (r *Repo) TaskUpdateDueDateIfOwned(taskID, userID int, day time.Time) (task Task, err error) {
	var (
		tx             *sql.Tx
		res            sql.Result
		rowsAffected   int64
		row            *sql.Row
		unixOldDueDate int64
		oldDueDate     time.Time
	)
	if !r.isConnected() {
		return Task{}, ErrNotConnected
	}
	if taskID < 1 {
		return Task{}, ErrTaskNotExist
	}
	if userID < 1 {
		return Task{}, ErrUserNotExist
	}

	if tx, err = r.db.Begin(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err = tx.QueryRow(`SELECT due_date FROM task WHERE id = ?`, taskID).Scan(&unixOldDueDate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, fmt.Errorf("%w for %d", ErrTaskNotExist, taskID)
		}
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if day.IsZero() {
		res, err = tx.Exec(`UPDATE task SET due_date = NULL WHERE id = ? AND user_id = ?`, taskID, userID)
	} else {
		oldDueDate = time.Unix(unixOldDueDate, 0)
		day = time.Date(day.Year(), day.Month(), day.Day(), oldDueDate.Hour(), oldDueDate.Minute(), 0, 0, time.Local)
		res, err = tx.Exec(`UPDATE task SET due_date = ? WHERE id = ? AND user_id = ?`, day.Unix(), taskID, userID)
	}

	if rowsAffected, err = res.RowsAffected(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	row = tx.QueryRow(fmt.Sprintf(`
		%s 
		WHERE task.id = ?`, defaultTaskSelect),
		taskID)

	// Happy path
	if rowsAffected == 1 {
		if task, err = scanRow(row); err != nil {
			return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}
		if err = tx.Commit(); err != nil {
			return Task{}, fmt.Errorf("%w: %s", ErrTransactionCommit, err)
		}

		return task, nil
	}

	if task, err = scanRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, fmt.Errorf("%w for %d", ErrTaskNotExist, taskID)
		}
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if task.UserID != userID {
		return Task{}, fmt.Errorf("%w user: %d, task: %d", ErrNotOwner, userID, taskID)
	}

	return Task{}, ErrUnknown
}
