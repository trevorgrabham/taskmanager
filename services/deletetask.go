package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// DeleteTask deletes the task identified by taskID if it is owned by userID.
//
// If taskID is empty, an ErrInvalidTaskID is returned.
// If userID is empty, an ErrInvalidUserID is returned.
// If userID is not the owner of taskID, an ErrNotOwner is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) DeleteTask(taskID, userID int) error {
	var err error
	if taskID < 1 {
		return &ErrTaskValidation{
			Field:   TaskFieldTaskID,
			Message: MessageInvalidTaskID,
		}
	}
	if userID < 1 {
		return &ErrUserValidation{
			Field:   UserFieldID,
			Message: MessageInvalidUserID,
		}
	}

	if err = s.repo.DeleteTaskIfOwned(taskID, userID); err != nil {
		switch sqlite.MapErrorToTypeName(err) {
		case "ErrNotOwner":
			return &ErrTaskValidation{
				Field: TaskFieldTaskID,
				Message: MessageNotOwner,
			}
		default:
			return fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}
	}

	return nil
}
