package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// GetUnscheduledTasks returns the unscheduled tasks associated with the userID.
//
// If userID is empty, an ErrInvalidUserID is returned.
func (s Service) GetUnscheduledTasks(userID int) (tasks TaskList, err error) {
	var (
		caller    = "GetUnscheduledTasks"
		repoTasks []sqlite.Task
	)
	if userID < 1 {
		return nil, fmt.Errorf("%s: %w", caller, ErrInvalidUserID)
	}

	if repoTasks, err = s.repo.GetUnscheduledTasks(userID); err != nil {
		return nil, fmt.Errorf("%s: %w", caller, err)
	}

	tasks = parseRepoTaskListToTaskList(repoTasks)

	return tasks, nil
}
