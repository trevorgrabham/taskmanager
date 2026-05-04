package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetTaskByID returns the task identified by taskID.
//
// If taskID is empty, or no task matches, acts as a no-op.
// If userID is not the owner, return ErrNotOwner.
// If a transient error occurs in the Repo, returns ErrInternalRepo.
func (r *Repo) GetTaskByID(taskID, userID int) (t Task, err error) {
	row := r.db.QueryRow(fmt.Sprintf(`
		%s
		WHERE task.id = ? AND user_id = ?`, defaultTaskSelect),
		taskID, userID)
	if t, err = scanRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var taskExists int 
			if err = r.db.QueryRow(`SELECT 1 FROM task WHERE id = ?`, taskID).Scan(&taskExists); err != nil {
				if errors.Is(err, sql.ErrNoRows) { return Task{}, nil }			// TaskID doesn't exist
				return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err) 
			}

			return  Task{}, ErrNotOwner
		}
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return t, nil
}
