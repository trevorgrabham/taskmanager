package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// AddTask adds a new task to the Repo.
//
// If task.ID is not empty, an ErrInvalidTaskID is returned. Consider using UpdateTask(Task) instead.
// If task.UserID is empty, an ErrInvalidUserID is returned.
// If task.Title is empty, an ErrInvalidTitle is returned.
// If task.RecurringPeriod is not a valid format, an ErrInvalidRecurringPeriod is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) AddTask(t Task) (Task, error) {
	var (
		err       error
		addedTask sqlite.Task
		repoTask  sqlite.Task
	)
	if t.ID != 0 {
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
	// Consider error on t.Category == ""
	if t.Done {
		return Task{}, &ErrTaskValidation{
			Field:   TaskFieldDone,
			Message: MessageTaskAlreadyCompleted,
		}
	}
	if !t.CompletionDate.IsZero() {
		return Task{}, &ErrTaskValidation{
			Field:   TaskFieldCompletionDate,
			Message: MessageTaskAlreadyCompleted,
		}
	}

	repoTask = parseTaskToRepoTask(t)
	if addedTask, err = s.repo.AddTask(repoTask); err != nil {
		switch sqlite.MapErrorToTypeName(err) {
		case "ErrConstraintFailure":
			return Task{}, &ErrUserValidation{Field: UserFieldID, Message: MessageInvalidUserID}
		default:
			return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}
	}

	return parseRepoTaskToTask(addedTask), nil
}
