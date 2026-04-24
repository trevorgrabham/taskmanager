package db

import (
	"database/sql"
	"fmt"
)

func (r *Repo) AddTask(t Task) (addedTask Task, err error) {
	var (
		caller = "AddTask"
		tx *sql.Tx
		recurringID int
		row *sql.Row
	)
	if !r.isConnected() {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}
	if t.Done {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskAlreadyDone)
	}
	if t.CompletionDate.Valid {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskAlreadyDone)
	}

	if tx, err = r.db.Begin(); err != nil {
		return Task{}, NewErrRepo(err)
	}
	defer func() {_ = tx.Rollback() }()

	if t.RecurringPeriod.Valid {
		if recurringID, err = insertRecurringPeriodIfNotExists(tx, t.RecurringPeriod.String); err != nil {
			return Task{}, NewErrRepo(err)
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
		return Task{}, NewErrRepo(err)
	}

	if err = tx.Commit(); err != nil {
		return Task{}, NewErrTransactionCommit(err)
	}

	return addedTask, nil
}
