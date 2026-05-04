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
		row         *sql.Row
	)
	if tx, err = r.db.Begin(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	defer func() { _ = tx.Rollback() }()

	row = tx.QueryRow(`
		INSERT INTO task(title, category, description, due_date, recurring_period, user_id) 
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING task.id, title, category, description, due_date, completion_date, done, recurring_period, user_id`,
		t.Title, t.Category, t.Description, t.DueDate, t.RecurringPeriod, t.UserID)

	if addedTask, err = scanRow(row); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if err = tx.Commit(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return addedTask, nil
}
