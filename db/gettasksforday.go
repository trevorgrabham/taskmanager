package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (r *Repo) GetTasksForDay(userID int, day time.Time) (tasks []Task, err error) {
	var (
		caller          = "GetTasksForDay"
		rows            *sql.Rows
		now, start, end time.Time
	)
	if !r.isConnected() {
		return nil, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	now = time.Now()
	start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	rows, err = r.db.Query(fmt.Sprintf(`
		%s
		WHERE done = 0 AND due_date >= ? AND due_date <= ? AND user_id = ?`, defaultTaskSelect), 
		start.Unix(), end.Unix(), userID)
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
