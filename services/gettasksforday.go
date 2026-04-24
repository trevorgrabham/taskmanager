package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
	"time"
)

func (s Service) GetTasksForDay(userID int, day time.Time) (tasks TaskList, err error) {
	var (
		caller    = "GetTasksForDay"
		repoTasks []sqlite.Task
	)
	if userID < 1 {
		return nil, fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}
	if day.IsZero() {
		return nil, fmt.Errorf("%s: %w", caller, ErrNoDate)
	}

	if repoTasks, err = s.repo.GetTasksForDay(userID, day); err != nil {
		return nil, fmt.Errorf("%s: %w", caller, err)
	}
	tasks = parseRepoTaskListToTaskList(repoTasks)

	return tasks, nil
}
