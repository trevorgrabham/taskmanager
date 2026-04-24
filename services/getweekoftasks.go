package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
	"time"
)

func (s Service) GetWeekOfTasks(userID int, startDay time.Time) (tasks TaskList, err error) {
	var (
		caller    = "GetWeekOfTasks"
		repoTasks []sqlite.Task
	)
	if userID < 1 {
		return nil, fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}

	if startDay.IsZero() {
		return nil, fmt.Errorf("%s: %w", caller, ErrNoDate)
	}

	if repoTasks, err = s.repo.GetWeekOfTasks(userID, startDay); err != nil {
		return nil, fmt.Errorf("%s: %w", caller, err)
	}
	tasks = parseRepoTaskListToTaskList(repoTasks)

	return tasks, nil
}
