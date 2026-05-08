package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
	"time"
)

// GetWeekOfTasks returns the tasks for userID for the range [startDay, startDay+7].
//
// If userID is empty, an ErrInvalidUserID is returned.
// If startDay is empty, an ErrInvalidDate is returned.
func (s Service) GetWeekOfTasks(userID int, startDay time.Time) (tasks TaskList, err error) {
	var repoTasks []sqlite.Task
	if userID < 1 {
		return nil, &ErrUserValidation{
			Field:   UserFieldID,
			Message: MessageInvalidUserID,
		}
	}

	if startDay.IsZero() {
		return nil, ErrEmptyDay
	}

	if repoTasks, err = s.repo.GetWeekOfTasks(userID, startDay); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	tasks = parseRepoTaskListToTaskList(repoTasks)

	return tasks, nil
}
