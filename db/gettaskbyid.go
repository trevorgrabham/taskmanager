package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetTaskByID returns the task identified by taskID.
//
// If Repo is not initialized, returns an ErrNotConnected.
// If taskID is empty or no task matches, returns an ErrTaskNotExist.
// If an error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) GetTaskByID(taskID int) (t Task, err error) {
	var row *sql.Row
	if !r.isConnected() {
		return Task{}, ErrNotConnected
	}
	if taskID < 1 {
		return Task{}, fmt.Errorf("%w for id %d", ErrTaskNotExist, taskID)
	}

	row = r.db.QueryRow(fmt.Sprintf(`
		%s
		WHERE task.id = ?`, defaultTaskSelect),
		taskID)
	if t, err = scanRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, fmt.Errorf("%w for id %d", ErrTaskNotExist, taskID)
		}
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return t, nil
}
