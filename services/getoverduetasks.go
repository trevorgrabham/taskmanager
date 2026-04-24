package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

func (s Service) GetOverdueTasks(userID int) (tasks TaskList, err error) {
	var (
		caller    = "GetOverdueTasks"
		repoTasks []sqlite.Task
	)
	if userID < 1 {
		return nil, fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}

	if repoTasks, err = s.repo.GetOverdueTasks(userID); err != nil {
		return nil, fmt.Errorf("%s: %w", caller, err)
	}
	tasks = parseRepoTaskListToTaskList(repoTasks)

	return tasks, nil
}
