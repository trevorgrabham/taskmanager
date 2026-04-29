package db

import (
	"database/sql"
	"fmt"
)

// AddTask adds the task to the Repo.
//
// If Repo is not initialized, returns an ErrNotConnected.
// If taskID is not empty, returns an ErrInvalidTask.
// If userID is empty, returns an ErrUserNotExist.
// If the task is completed (done is true, or completiondate is not zero), returns an ErrTaskCompleted.
// If title is empty, returns an ErrInvalidTask.
// If an error occurs in the Repo, returns an ErrInternalRepo.
// If an error occurs commiting the Repo transaction, an ErrTransactionCommit is returned.
func (r *Repo) AddTask(t Task) (addedTask Task, err error) {
	var (
		tx          *sql.Tx
		recurringID int
		row         *sql.Row
	)
	if !r.isConnected() { return Task{}, ErrNotConnected }
	if t.ID != 0 { return Task{}, ErrInvalidTask }
	if t.UserID < 1 { return Task{}, ErrUserNotExist }
	if t.Done || t.CompletionDate.Valid {
		return Task{}, ErrTaskCompleted
	}
	if t.Title == "" { return Task{}, ErrInvalidTask }

	if tx, err = r.db.Begin(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	defer func() { _ = tx.Rollback() }()

	if t.RecurringPeriod.Valid {
		if recurringID, err = insertRecurringPeriodIfNotExists(tx, t.RecurringPeriod.String); err != nil {
			return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}

		t.RecurringID.Valid = true
		t.RecurringID.Int64 = int64(recurringID)
	}

	row = tx.QueryRow(`
		INSERT INTO task(title, category, description, due_date, recurring_id, user_id) 
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING task.id, title, category, description, due_date, completion_date, done, recurring_id, user_id, (
			SELECT period FROM recurring WHERE recurring.id = task.recurring_id)`,
		t.Title, t.Category, t.Description, t.DueDate, t.RecurringID, t.UserID)

	if addedTask, err = scanRow(row); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if err = tx.Commit(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrTransactionCommit, err)
	}

	return addedTask, nil
}
