package db

import (
	"database/sql"
	"fmt"
)

// AddTask adds the task to the Repo.
//
// 'ID', 'Done', and 'CompletionDate' fields are ignored if set.
//
// If UserID does not exist, returns ErrInternalRepo (Foreign key constraint)
// If Title is empty, returns ErrInternalRepo (Check constraint)
// If a transient error occurs in the Repo, returns ErrInternalRepo.
func (r *Repo) AddTask(t Task) (addedTask Task, err error) {
	var (
		tx          *sql.Tx
		recurringID int
		row         *sql.Row
	)
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
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return addedTask, nil
}
