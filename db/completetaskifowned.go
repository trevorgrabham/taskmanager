package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (r *Repo) CompleteTaskIfOwned(taskID, userID int) (err error) {
	var (
		caller       = "CompleteTaskIfOwned"
		t            Task
		res          sql.Result
		rowsAffected int64
	)
	if !r.isConnected() {
		return fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	res, err = r.db.Exec(`
		UPDATE task 
		SET completion_date = ?, done = 1 
		WHERE id = ? AND user_id = ?`,
		time.Now().Unix(), taskID, userID)
	if err != nil {
		return NewErrRepo(err)
	}

	if rowsAffected, err = res.RowsAffected(); err != nil {
		return NewErrRepo(err)
	}

	// Successful update
	if rowsAffected == 1 {
		return nil
	}

	if t, err = r.GetTaskByID(taskID); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return fmt.Errorf("%s: %w", caller, NewErrTaskNotExist(taskID)) }
		return fmt.Errorf("%s: %w", caller, err)
	}

	if t.UserID != userID {
		return fmt.Errorf("%s: %w", caller, NewErrNotOwner(userID, taskID))
	}

	return fmt.Errorf("%s: %w", caller, NewErrRepo(ErrUnknown))
}
