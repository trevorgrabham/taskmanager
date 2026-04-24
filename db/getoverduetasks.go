package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (r *Repo) GetOverdueTasks(userID int) (tasks []Task, err error) {
	var (
		caller                   = "GetOverdueTasks"
		rows                     *sql.Rows
		now, yesterdayAtMidnight time.Time
	)
	if !r.isConnected() {
		return nil, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	now = time.Now()
	yesterdayAtMidnight = time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, time.Local)

	rows, err = r.db.Query(fmt.Sprintf(`
		%s
		WHERE done = 0 AND due_date <= ? AND user_id = ?`, defaultTaskSelect), 
		yesterdayAtMidnight.Unix(), userID)
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
