package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// UpdateTask updates a task identified by task.ID.
//
// If task.ID is empty, an ErrInvalidTaskID is returned. Consider using AddTask(Task) instead.
// If task.UserID is empty, an ErrInvalidUserID is returned.
// If task.Title is empty, an ErrInvalidTitle is returned.
// If task.RecurringPeriod is not a valid format, an ErrInvalidRecurringPeriod is returned.
// If task.UserID is not the owner of the task, an ErrNotOwner is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) UpdateTask(t Task) (updatedTask Task, err error) {
	var repoTask sqlite.Task
	if t.ID < 1 {
		return Task{}, &ErrTaskValidation{
			Field:   TaskFieldTaskID,
			Message: MessageInvalidTaskID,
		}
	}
	if t.UserID < 1 {
		return Task{}, &ErrUserValidation{
			Field:   UserFieldID,
			Message: MessageInvalidUserID,
		}
	}
	if t.Title == "" {
		return Task{}, &ErrTaskValidation{
			Field:   TaskFieldTitle,
			Message: MessageEmptyTitle,
		}
	}
	if !t.CompletionDate.IsZero() {
		return Task{}, &ErrTaskValidation{
			Field: TaskFieldCompletionDate,
			Message: MessageTaskAlreadyCompleted,
		}
	}
	if t.Done {
		return Task{}, &ErrTaskValidation{
			Field: TaskFieldDone,
			Message: MessageTaskAlreadyCompleted,
		}
	}
	// Consider making Category required. If we do, check t.Category == "" here

	repoTask = parseTaskToRepoTask(t)
	if repoTask, err = s.repo.UpdateTaskIfOwned(repoTask); err != nil {
		switch sqlite.MapErrorToTypeName(err) {
		case "ErrNotOwner":
			return Task{}, &ErrTaskValidation{
				Field: TaskFieldTaskID,
				Message: MessageNotOwner,
			}
		default: 
			return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err) 
		}
	}
	if repoTask == (sqlite.Task{}) { 
		return Task{}, &ErrTaskValidation{
			Field: TaskFieldTaskID, 
			Message: MessageTaskNotFound,
		} 
	}

	updatedTask = parseRepoTaskToTask(repoTask)
	return updatedTask, nil
}
