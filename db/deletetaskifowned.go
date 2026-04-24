package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repo) DeleteTaskIfOwned(taskID, userID int) (err error) {
	var (
		caller       = "DeleteTaskIfOwned"
		res          sql.Result
		rowsAffected int64
		t            Task
	)
	if !r.isConnected() {
		return fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	res, err = r.db.Exec(`DELETE FROM task WHERE id = ? AND user_id = ?`, taskID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	// Happy path
	if rowsAffected > 0 {
		return nil
	}

	if t, err = r.GetTaskByID(taskID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", caller, NewErrTaskNotExist(taskID))
		}
		return fmt.Errorf("%s: %w", caller, err)
	}

	if t.UserID != userID {
		return fmt.Errorf("%s: %w", caller, NewErrNotOwner(userID, taskID))
	}

	return fmt.Errorf("%s: %w", caller, ErrUnknown)
}
