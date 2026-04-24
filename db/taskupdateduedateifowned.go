package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (r *Repo) TaskUpdateDueDateIfOwned(taskID, userID int, day time.Time) (task Task, err error) {
	var (
		caller                 = "TaskUpdateDueDateIfOwned"
		tx                     *sql.Tx
		row                    *sql.Row
		oldDueDate, newDueDate time.Time
		actualUserID           int
	)
	if !r.isConnected() {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	if tx, err = r.db.Begin(); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}
	defer func() { _ = tx.Rollback() }()

	row = tx.QueryRow(fmt.Sprint(`
    %s 
    WHERE id = ? AND user_id = ?`, defaultTaskSelect),
		taskID, userID)

	if task, err = scanRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = tx.QueryRow(`SELECT user_id FROM task WHERE id = ?`, taskID).Scan(&actualUserID)
			if err != nil && errors.Is(err, sql.ErrNoRows) {
				return Task{}, fmt.Errorf("%s: %w", caller, NewErrTaskNotExist(taskID))
			}

			if actualUserID != userID {
				return Task{}, fmt.Errorf("%s: %w", caller, NewErrNotOwner(userID, taskID))
			}

			return Task{}, fmt.Errorf("%s: %w", caller, ErrUnknown)
		}

		return Task{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	if !task.DueDate.Valid {
		newDueDate = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
	} else {
		oldDueDate = time.Unix(task.DueDate.Int64, 0)
		newDueDate = time.Date(day.Year(), day.Month(), day.Day(), oldDueDate.Hour(), oldDueDate.Minute(), 0, 0, time.Local)
	}

	row = tx.QueryRow(`UPDATE task SET due_date = ? WHERE id = ? AND user_id = ?`, newDueDate.Unix(), taskID, userID)

	if task, err = scanRow(row); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	if err = tx.Commit(); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, NewErrTransactionCommit(err))
	}

	return task, nil
}
