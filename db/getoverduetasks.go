package db

import (
	"database/sql"
	"fmt"
	"time"
)

// GetOverdueTasks returns the tasks for userID that were due before today at 00:00.
//
// If userID is empty, or doesn't exist, treated as a no-op.
// If a transient error occurs in the Repo, returns ErrInternalRepo.
func (r *Repo) GetOverdueTasks(userID int) (tasks []Task, err error) {
	var (
		rows                     *sql.Rows
		now, yesterdayAtMidnight time.Time
	)
	now = time.Now()
	yesterdayAtMidnight = time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, time.Local)

	rows, err = r.db.Query(fmt.Sprintf(`
		%s
		WHERE done = 0 AND due_date <= ? AND user_id = ?`, defaultTaskSelect),
		yesterdayAtMidnight.Unix(), userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if tasks, err = scanRows(rows); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return tasks, nil
}
