package db

import (
	"database/sql"
	"fmt"
	"time"
)

// CompleteTaskIfOwned completes the task identified by taskID if it is owned by userID.
//
// If the task identified by TaskID is already completed, it's treated as a no-op.
//
// If taskID is empty or no task matches, returns ErrTaskNotFound.
// If userID is empty, not the owner, or does not exist, returns ErrNotOwner.
// If a transient error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) CompleteTaskIfOwned(taskID, userID int) (err error) {
	var (
		t            Task
		res          sql.Result
		rowsAffected int64
	)
	res, err = r.db.Exec(`
		UPDATE task 
		SET completion_date = ?, done = 1 
		WHERE done = 0 AND completion_date IS NULL AND id = ? AND user_id = ?`,
		time.Now().Unix(), taskID, userID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if rowsAffected, err = res.RowsAffected(); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	// Successful update
	if rowsAffected == 1 {
		return nil
	}

	if t, err = r.GetTaskByID(taskID, userID); err != nil {
		return err
	}
	if t == (Task{}) {
		return ErrTaskNotFound
	}

	// Task was already completed
	return nil
}
