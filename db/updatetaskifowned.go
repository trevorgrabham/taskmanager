package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

// UpdateTaskIfOwned updates the task identified by task.ID if it is owned by task.UserID.
//
// If taskID is empty or doesn't exist, no-op.
// If userID is empty or not the owner, returns ErrNotOwner.
// If title is empty, returns ErrInternalRepo.
// If an error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) UpdateTaskIfOwned(t Task) (updatedTask Task, err error) {
	var (
		tx           *sql.Tx
		res          sql.Result
		rowsAffected int64
		row          *sql.Row
		targetErr    sqlite3.Error
	)
	if tx, err = r.db.Begin(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err = tx.Exec(`
    UPDATE task 
    SET 
      title = ?,
      category = ?,
      description = ?,
      due_date = ?,
      recurring_period = ?
    WHERE id = ? AND user_id = ?`,
		t.Title, t.Category, t.Description, t.DueDate, t.RecurringPeriod, t.ID, t.UserID)
	if err != nil {
		if errors.As(err, &targetErr) && targetErr.ExtendedCode == sqlite3.ErrConstraintCheck {
			return Task{}, fmt.Errorf("%w: %s", ErrConstraintFailure, err)
		}

		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	row = tx.QueryRow(fmt.Sprintf(`
		%s
		WHERE task.id = ?`, defaultTaskSelect),
		t.ID)

	if rowsAffected != 1 {
		if updatedTask, err = scanRow(row); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return Task{}, nil
			}
			return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}
		if updatedTask.UserID != t.UserID {
			return Task{}, ErrNotOwner
		}

		return Task{}, ErrInternalRepo
	}

	if updatedTask, err = scanRow(row); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if err = tx.Commit(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return updatedTask, nil
}
