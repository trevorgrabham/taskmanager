package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
)

// TaskToggleComplete toggles the completion status of the task identified by taskID if it is owned by the user identified by userID.
//
// If taskID is empty, an ErrInvalidTaskID is returned.
// If userID is empty, an ErrInvalidUserID is returned.
// If userID is not the owner of taskID, an ErrNotOwner is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) TaskToggleComplete(taskID, userID int) (updatedTask Task, err error) {
	var (
		caller   = "TaskToggleComplete"
		repoTask sqlite.Task
	)
	if taskID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrInvalidTaskID)
	}
	if userID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrInvalidUserID)
	}

	if repoTask, err = s.repo.TaskToggleCompleteIfOwned(taskID, userID); err != nil {
		if errors.Is(err, sqlite.ErrNotOwner) {
			return Task{}, fmt.Errorf("%s: %w", caller, ErrNotOwner)
		}
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}

	updatedTask = parseRepoTaskToTask(repoTask)

	return updatedTask, nil
}
