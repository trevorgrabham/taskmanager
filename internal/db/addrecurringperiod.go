package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
)

func AddRecurringPeriod(db *sql.DB, t task.Task) error {
	if db == nil {
		return fmt.Errorf("adding recurring period: cannot add into a nil database")
	}
	if t.IsZero() {
		return fmt.Errorf("adding recurring period: need a period and id to add a recurring period")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("adding recurring period: %s", err)
	}
	defer func() { _ = tx.Rollback() }()

	if t.RecurringPeriod != "" {
		var existed sql.NullBool
		err = tx.QueryRow(`UPDATE recurring SET period = ? WHERE task_id = ? RETURNING 1`, t.RecurringPeriod, t.ID).Scan(&existed)
		if !existed.Valid {
			_, err = tx.Exec(`INSERT INTO recurring(task_id, period) SELECT ?, ?`, t.ID, t.RecurringPeriod)
		}
	} else {
		_, err = tx.Exec(`DELETE FROM recurring WHERE task_id = ?`, t.ID)
	}
	if err != nil {
		return fmt.Errorf("adding recurring period: %s", err)
	}

	_ = tx.Commit()
	return nil
}
