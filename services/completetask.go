package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
)

func (s Service) CompleteTask(taskID, userID int) error {
	var (
		caller = "CompleteTask"
		err    error
	)
	if taskID < 1 {
		return fmt.Errorf("%s: %w", caller, ErrTaskBadID)
	}
	if userID < 1 {
		return fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}

	if err = s.repo.CompleteTaskIfOwned(taskID, userID); err != nil {
		if targetErr := (sqlite.ErrNotOwner{}); errors.As(err, &targetErr) {
			return fmt.Errorf("%s: %w", caller, ErrUserWrongID)
		}
		return fmt.Errorf("%s: %w", caller, err)
	}

	return nil
}
