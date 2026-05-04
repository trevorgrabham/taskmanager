package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TaskUpdateDueDateIfOwned updates the DueDate for the task identified by task.ID if it is owned by task.UserID. If DueDate was not NULL, then the date changes, but the time stays the same. Returns the updated Task.
//
// If taskID is empty or doesn't exist, no-op.
// If userID is empty, doesn't exist, or isn't the owner, returns ErrNotOwner.
// If an error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) TaskUpdateDueDateIfOwned(taskID, userID int, day time.Time) (task Task, err error) {
	var (
		tx             *sql.Tx
		res            sql.Result
		rowsAffected   int64
		row            *sql.Row
		unixOldDueDate sql.NullInt64
		oldDueDate     time.Time
	)
	if tx, err = r.db.Begin(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	defer func() { _ = tx.Rollback() }()

	if day.IsZero() {
		res, err = tx.Exec(`UPDATE task SET due_date = NULL WHERE id = ? AND user_id = ?`, taskID, userID)
	} else {
		// store old value to keep the same time on the new date
		if err = tx.QueryRow(`SELECT due_date FROM task WHERE id = ?`, taskID).Scan(&unixOldDueDate); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return Task{}, nil
			}
			return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}

		if unixOldDueDate.Valid {
			oldDueDate = time.Unix(unixOldDueDate.Int64, 0)
		} 
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
			return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}

		return task, nil
	}

	if task, err = scanRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, nil
		}
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if task.UserID != userID {
		return Task{}, ErrNotOwner
	}

	return Task{}, ErrInternalRepo
}
