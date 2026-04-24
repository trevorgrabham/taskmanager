package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (r *Repo) GetWeekOfTasks(userID int, startDay time.Time) (tasks []Task, err error) {
	var (
		caller     = "GetWeekOfTasks"
		rows       *sql.Rows
		start, end time.Time
	)
	if !r.isConnected() {
		return nil, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	start = time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, time.Local)
	end = time.Date(startDay.Year(), startDay.Month(), startDay.Day()+6, 23, 59, 59, 0, time.Local)
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
