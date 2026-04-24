package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
)

func (s Service) DeleteTask(taskID, userID int) error {
	var (
		caller = "DeleteTask"
		err    error
	)
	if taskID < 1 {
		return fmt.Errorf("%s: %w", caller, ErrTaskBadID)
	}
	if userID < 1 {
		return fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}

	if err = s.repo.DeleteTaskIfOwned(taskID, userID); err != nil {
		if targetErr := (sqlite.ErrNotOwner{}); errors.Is(err, &targetErr) {
			return fmt.Errorf("%s: %w", caller, ErrUserWrongID)
		}
		return fmt.Errorf("%s: %w", caller, err)
	}

	return nil
}
