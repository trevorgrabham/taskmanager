package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// UpdateTaskIfOwned updates the task identified by task.ID if it is owned by task.UserID.
//
// If Repo is not initialized, returns an ErrNotConnected.
// If taskID is empty, returns an ErrTaskNotExist.
// If userID is empty, returns an ErrUserNotExist.
// If title is empty, returns an ErrInvalidTask.
// If the task is not owned by task.UserID, returns an ErrNotOwner.
// If the no task exists with id = task.ID, returns an ErrTaskNotExist.
// If an error occurs in the Repo, returns an ErrInternalRepo.
// If an error occurs commiting the Repo transaction, an ErrTransactionCommit is returned.
func (r *Repo) UpdateTaskIfOwned(t Task, userID int) (updatedTask Task, err error) {
	var (
		tx           *sql.Tx
		res          sql.Result
		rowsAffected int64
		row          *sql.Row
	)
	if !r.isConnected() {
		return Task{}, ErrNotConnected
	}
	if t.ID < 1 {
		return Task{}, ErrTaskNotExist
	}
	if t.UserID < 1 {
		return Task{}, ErrUserNotExist
	}
	if t.Title == "" {
		return Task{}, ErrInvalidTask
	}

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
      done = ?,
      completion_date = ?,
      recurring_id = ?
    WHERE id = ? AND user_id = ?`,
		t.Title, t.Category, t.Description, t.DueDate, t.Done, t.CompletionDate, t.RecurringID, t.ID, userID)
	if err != nil {
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
				return Task{}, fmt.Errorf("%w for id %d", ErrTaskNotExist, t.ID)
			}
			return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}
		if updatedTask.UserID != t.UserID {
			return Task{}, ErrNotOwner
		}

		return Task{}, ErrUnknown
	}

	if updatedTask, err = scanRow(row); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if err = tx.Commit(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrTransactionCommit, err)
	}

	return updatedTask, nil
}
