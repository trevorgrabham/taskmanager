package db

import (
	"database/sql"
	"fmt"
)

// GetUnscheduledTasks returns the unscheduled tasks for userID.
//
// If userID is empty, treated as a no-op.
// If an error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) GetUnscheduledTasks(userID int) (tasks []Task, err error) {
	var rows *sql.Rows
	rows, err = r.db.Query(fmt.Sprintf(`
		%s 
		WHERE done = 0 AND due_date IS NULL AND user_id = ?`, defaultTaskSelect),
		userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if tasks, err = scanRows(rows); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return tasks, nil
}
