package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
)

// CompleteTask completes the task identified by taskID if it is owned by the user identified by userID.
//
// If taskID is empty, an ErrInvalidTaskID is returned.
// If userID is empty, an ErrInvalidUserID is returned.
// If userID is not the owner of taskID, an ErrNotOwner is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) CompleteTask(taskID, userID int) error {
	var (
		caller = "CompleteTask"
		err    error
	)
	if taskID < 1 {
		return fmt.Errorf("%s: %w", caller, ErrInvalidTaskID)
	}
	if userID < 1 {
		return fmt.Errorf("%s: %w", caller, ErrInvalidUserID)
	}

	if err = s.repo.CompleteTaskIfOwned(taskID, userID); err != nil {
		if errors.Is(err, sqlite.ErrNotOwner) { return fmt.Errorf("%s: %w", caller, ErrNotOwner) }
		return fmt.Errorf("%s: %w", caller, err)
	}

	return nil
}
