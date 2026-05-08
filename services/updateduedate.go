package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
	"time"
)

// UpdateDueDate updates the DueDate for the task identified by taskID.
//
// If taskID is empty, an ErrInvalidTaskID is returned. Consider using AddTask(Task) instead.
// If userID is empty, an ErrInvalidUserID is returned.
// If userID is not the owner of the task, an ErrNotOwner is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) UpdateDueDate(taskID, userID int, dueDay time.Time) (updatedTask Task, err error) {
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

	if repoTask, err = s.repo.TaskUpdateDueDateIfOwned(taskID, userID, dueDay); err != nil {
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

	updatedTask = parseRepoTaskToTask(repoTask)
	return updatedTask, nil
}
