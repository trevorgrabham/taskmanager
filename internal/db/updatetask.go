package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
)

func UpdateTask(db *sql.DB, t task.Task) error {
	if db == nil {
		return fmt.Errorf("updating task: cannot update for a nil database")
	}
	if t.IsZero() {
		return fmt.Errorf("updating task: cannot update without a task")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("updating task: %s", err)
	}
	defer func() { _ = tx.Rollback() }()

	if t.RecurringPeriod != "" {
		var recurringID int
		err = tx.QueryRow(`INSERT INTO recurring(period) VALUES (?) ON CONFLICT(period) DO UPDATE SET period = excluded.period RETURNING id`, t.RecurringPeriod).Scan(&recurringID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("updating task: %s", err)
		}

		if t.DueDate.IsZero() {
			_, err = tx.Exec(`UPDATE task SET title = ?, category = ?, description = ?, due_date = NULL, recurring_id = ? WHERE id = ?`, t.Title, t.Category, t.Description, recurringID, t.ID)
		} else {
			_, err = tx.Exec(`UPDATE task SET title = ?, category = ?, description = ?, due_date = ?, recurring_id = ? WHERE id = ?`, t.Title, t.Category, t.Description, t.DueDate.Unix(), recurringID, t.ID)
		}
		if err != nil {
			return fmt.Errorf("updating task: %s", err)
		}
	} else {
		if t.DueDate.IsZero() {
			_, err = tx.Exec(`UPDATE task SET title = ?, category = ?, description = ?, due_date = NULL, recurring_id = NULL WHERE id = ?`, t.Title, t.Category, t.Description, t.ID)
		} else {
			_, err = tx.Exec(`UPDATE task SET title = ?, category = ?, description = ?, due_date = ?, recurring_id = NULL WHERE id = ?`, t.Title, t.Category, t.Description, t.DueDate.Unix(), t.ID)
		}
		if err != nil {
			return fmt.Errorf("updating task: %s", err)
		}
	}

	_ = tx.Commit()

	return nil
}
