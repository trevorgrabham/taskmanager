package db

import (
	"database/sql"
	"fmt"
)

func (r *Repo) UpdateTaskIfOwned(t Task, userID int) (updatedTask Task, err error) {
	var (
		caller = "UpdateTaskIfOwned"
		tx     *sql.Tx
		row    *sql.Row
	)
	if !r.isConnected() {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	if t.ID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, NewErrTaskNotExist(t.ID))
	}

	if tx, err = r.db.Begin(); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}
	defer func() { _ = tx.Rollback() }()

	row = tx.QueryRow(`
    UPDATE task 
    SET 
      title = ?,
      category = ?,
      description = ?,
      due_date = ?,
      done = ?,
      completion_date = ?,
      recurring_id = ?
    WHERE id = ? AND user_id = ?
    RETURNING id, title, category, description, due_date, completion_date, done, recurring_id, user_id, (
      SELECT period FROM recurring WHERE id = task.recurring_id
    )`,
		t.Title, t.Category, t.Description, t.DueDate, t.Done, t.CompletionDate, t.RecurringID, t.ID, userID)

	if updatedTask, err = scanRow(row); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	if err = tx.Commit(); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, NewErrTransactionCommit(err))
	}

	return updatedTask, nil
}
