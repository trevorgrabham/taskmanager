package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
	"time"
)

// GetTasksForDay returns the tasks associated with the userID with a DueDate on day.
//
// If userID is empty, an ErrInvalidUserID is returned.
// If day is empty, an ErrInvalidDate is returned.
func (s Service) GetTasksForDay(userID int, day time.Time) (tasks TaskList, err error) {
	var repoTasks []sqlite.Task
	if userID < 1 {
		return nil, &ErrUserValidation{
			Field:   UserFieldID,
			Message: MessageInvalidUserID,
		}
	}
	if day.IsZero() {
		return nil, ErrEmptyDay
	}

	if repoTasks, err = s.repo.GetTasksForDay(userID, day); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	tasks = parseRepoTaskListToTaskList(repoTasks)

	return tasks, nil
}
