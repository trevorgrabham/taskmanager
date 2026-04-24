package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repo) GetUnscheduledTasks(userID int) (tasks []Task, err error) {
	var (
		caller = "GetUnscheduledTasks"
		rows   *sql.Rows
	)
	if !r.isConnected() {
		return nil, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	rows, err = r.db.Query(fmt.Sprintf(`
		%s 
		WHERE done = 0 AND due_date IS NULL AND user_id = ?`, defaultTaskSelect),
		userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]Task, 0), nil
		}
		return nil, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	if tasks, err = scanRows(rows); err != nil {
		return nil, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	return tasks, nil
}
