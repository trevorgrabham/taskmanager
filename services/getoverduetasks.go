package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// GetOverdueTasks returns the list of tasks associated with userID that have DueDates before today at 00:00.
//
// If userID is empty, an ErrInvalidUserID is returned.
func (s Service) GetOverdueTasks(userID int) (tasks TaskList, err error) {
	var (
		caller    = "GetOverdueTasks"
		repoTasks []sqlite.Task
	)
	if userID < 1 { return nil, fmt.Errorf("%s: %w", caller, ErrInvalidUserID) }

	if repoTasks, err = s.repo.GetOverdueTasks(userID); err != nil {
		return nil, fmt.Errorf("%s: %w", caller, err)
	}
	tasks = parseRepoTaskListToTaskList(repoTasks)

	return tasks, nil
}
