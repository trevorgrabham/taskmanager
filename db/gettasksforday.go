package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// GetTasksForDay returns the tasks for userID due on day.
//
// If Repo is not initialized, returns an ErrNotConnected.
// If userID is empty, an ErrUserNotExist is returned.
// If day is empty, an ErrEmptyDate is returned.
// If an error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) GetTasksForDay(userID int, day time.Time) (tasks []Task, err error) {
	var (
		rows       *sql.Rows
		start, end time.Time
	)
	if !r.isConnected() {
		return nil, ErrNotConnected
	}
	if userID < 1 {
		return nil, ErrUserNotExist
	}
	if day.IsZero() {
		return nil, ErrEmptyDate
	}

	start = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
	end = time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, time.Local)
	rows, err = r.db.Query(fmt.Sprintf(`
		%s
		WHERE done = 0 AND due_date >= ? AND due_date <= ? AND user_id = ?`, defaultTaskSelect),
		start.Unix(), end.Unix(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]Task, 0), nil
		}
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if tasks, err = scanRows(rows); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return tasks, nil
}
