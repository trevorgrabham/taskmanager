package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// GetTaskByID retrieves that task identified by taskID if userID is the owner.
//
// If taskID is empty, an ErrInvalidTaskID is returned.
// If userID is empty, an ErrInvalidUserID is returned.
// If userID is not the owner of taskID, an ErrNotOwner is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) GetTaskByID(taskID, userID int) (t Task, err error) {
	var repoTask sqlite.Task
	if taskID < 1 {
		return Task{}, &ErrTaskValidation{
			Field:   TaskFieldTaskID,
			Message: MessageInvalidTaskID,
		}
	}
	if userID < 1 {
		return Task{}, &ErrUserValidation{
			Field:   UserFieldID,
			Message: MessageInvalidUserID,
		}
	}

	if repoTask, err = s.repo.GetTaskByID(taskID, userID); err != nil {
		switch sqlite.MapErrorToTypeName(err) {
		case "ErrNotOwner":
			return Task{}, &ErrTaskValidation{
				Field:   TaskFieldTaskID,
				Message: MessageNotOwner,
			}
		default:
			return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}
	}
	if repoTask == (sqlite.Task{}) {
		return Task{}, &ErrTaskValidation{
			Field:   TaskFieldTaskID,
			Message: MessageTaskNotFound,
		}
	}

	t = parseRepoTaskToTask(repoTask)

	if t.UserID != userID {
		return Task{}, &ErrTaskValidation{
			Field:   TaskFieldTaskID,
			Message: MessageNotOwner,
		}
	}

	return t, nil
}
