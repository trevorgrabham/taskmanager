package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// GetOverdueTasks returns the list of tasks associated with userID that have DueDates before today at 00:00.
//
// If userID is empty, an ErrInvalidUserID is returned.
func (s Service) GetOverdueTasks(userID int) (tasks TaskList, err error) {
	var repoTasks []sqlite.Task
	if userID < 1 {
		return nil, &ErrUserValidation{
			Field:   UserFieldID,
			Message: MessageInvalidUserID,
		}
	}

	if repoTasks, err = s.repo.GetOverdueTasks(userID); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	tasks = parseRepoTaskListToTaskList(repoTasks)

	return tasks, nil
}
