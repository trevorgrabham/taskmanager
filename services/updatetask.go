package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
)

func (s Service) UpdateTask(t Task, userID int) (updatedTask Task, err error) {
	var (
		caller   = "UpdateTask"
		repoTask sqlite.Task
	)
	if t.ID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskBadID)
	}
	if userID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}
	if t.Title == "" {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskNoTitle)
	}
	// Consider making Category required. If we do, check t.Category == "" here
	repoTask = parseTaskToRepoTask(t)
	if repoTask, err = s.repo.UpdateTaskIfOwned(repoTask, userID); err != nil {
		if targetErr := (sqlite.ErrNotOwner{}); errors.Is(err, &targetErr) {
			return Task{}, fmt.Errorf("%s: %w", caller, err)
		}
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}

	updatedTask = parseRepoTaskToTask(repoTask)
	return updatedTask, nil
}
