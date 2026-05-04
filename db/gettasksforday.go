package db

import (
	"database/sql"
	"fmt"
	"time"
)

// GetTasksForDay returns the tasks for userID due on day.
//
// If userID or day is empty, treated as a no-op.
// If an error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) GetTasksForDay(userID int, day time.Time) (tasks []Task, err error) {
	var (
		rows       *sql.Rows
		start, end time.Time
	)
	start = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
	end = time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, time.Local)
	rows, err = r.db.Query(fmt.Sprintf(`
		%s
		WHERE done = 0 AND due_date >= ? AND due_date <= ? AND user_id = ?`, defaultTaskSelect),
		start.Unix(), end.Unix(), userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if tasks, err = scanRows(rows); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return tasks, nil
}
