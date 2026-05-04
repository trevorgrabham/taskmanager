package db

import (
	"database/sql"
	"fmt"
	"time"
)

// GetWeekOfTasks returns the tasks for userID for the range [startDay, startDay+7].
//
// If userID is empty or doesn't exist, treats it as a no-op.
// If startDay is empty, treats it as a no-op.
// If an error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) GetWeekOfTasks(userID int, startDay time.Time) (tasks []Task, err error) {
	var (
		rows       *sql.Rows
		start, end time.Time
	)
	start = time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, time.Local)
	end = time.Date(startDay.Year(), startDay.Month(), startDay.Day()+6, 23, 59, 59, 0, time.Local)
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
