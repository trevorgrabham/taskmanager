package db

import (
	"database/sql"
	"fmt"
)

// DeleteTaskIfOwned deletes the task identified by taskID if it is owned by userID.
//
// If TaskID does not exist, treats it as a no-op.
// If a transient error occurs in the Repo, returns ErrInternalRepo.
// If UserID is not the owner, or doesn't exist, returns ErrNotOwner.
func (r *Repo) DeleteTaskIfOwned(taskID, userID int) (err error) {
	var (
		res          sql.Result
		rowsAffected int64
		t            Task
	)
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

	if t, err = r.GetTaskByID(taskID, userID); err != nil {
		return err
	}
	if t == (Task{}) {
		return nil
	}

	return ErrInternalRepo
}
