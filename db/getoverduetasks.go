package db

import (
	"database/sql"
	"fmt"
	"time"
)

// GetOverdueTasks returns the tasks for userID that were due before today at 00:00.
//
// If Repo is not initialized, returns an ErrNotConnected.
// If userID is empty, an ErrUserNotExist is returned.
// If an error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) GetOverdueTasks(userID int) (tasks []Task, err error) {
	var (
		rows                     *sql.Rows
		now, yesterdayAtMidnight time.Time
	)
	if !r.isConnected() {
		return nil, ErrNotConnected
	}
	if userID < 1 {
		return nil, fmt.Errorf("%w for id %d", ErrUserNotExist, userID)
	}

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
