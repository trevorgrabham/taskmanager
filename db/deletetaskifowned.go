package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// DeleteTaskIfOwned deletes the task identified by taskID if it is owned by userID.
//
// If Repo is not initialized, returns an ErrNotConnected.
// If taskID is empty or no task matches, returns an ErrTaskNotExist.
// If userID is empty, returns an ErrUserNotExist.
// If an error occurs in the Repo, returns an ErrInternalRepo.
// If the task is not owned by userID, returns an ErrNotOwner.
// If no task was completed, but the task exists and is owned by userID, returns an ErrorUnknown. This should never happen.
func (r *Repo) DeleteTaskIfOwned(taskID, userID int) (err error) {
	var (
		res          sql.Result
		rowsAffected int64
		t            Task
	)
	if !r.isConnected() {
		return ErrNotConnected
	}
	if taskID < 1 {
		return fmt.Errorf("%w for id %d", ErrTaskNotExist, taskID)
	}
	if userID < 1 {
		return fmt.Errorf("%w for id %d", ErrUserNotExist, userID)
	}

	res, err = r.db.Exec(`DELETE FROM task WHERE id = ? AND user_id = ?`, taskID, userID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	// Happy path
	if rowsAffected > 0 {
		return nil
	}

	if t, err = r.GetTaskByID(taskID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w for id %d", ErrTaskNotExist, taskID)
		}
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if t.UserID != userID {
		return ErrNotOwner
	}

	return ErrUnknown
}
