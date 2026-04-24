package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
)

func (s Service) TaskToggleComplete(taskID, userID int) (updatedTask Task, err error) {
	var (
		caller   = "TaskToggleComplete"
		repoTask sqlite.Task
	)
	if taskID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskBadID)
	}
	if userID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}

	if repoTask, err = s.repo.TaskToggleCompleteIfOwned(taskID, userID); err != nil {
		if targetErr := new(sqlite.ErrNotOwner); errors.As(err, &targetErr) {
			return Task{}, fmt.Errorf("%s: %w", caller, ErrUserWrongID)
		}
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}

	updatedTask = parseRepoTaskToTask(repoTask)
	return updatedTask, nil
}
