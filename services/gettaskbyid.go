package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// GetTaskByID retrieves that task identified by taskID if userID is the owner.
//
// If taskID is empty, an ErrInvalidTaskID is returned.
// If userID is empty, an ErrInvalidUserID is returned.
// If userID is not the owner of taskID, an ErrNotOwner is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) GetTaskByID(taskID, userID int) (t Task, err error) {
	var (
		caller   = "GetTaskByID"
		repoTask sqlite.Task
	)
	if taskID < 1 {
		return Task{}, ErrInvalidTaskID
	}
	if userID < 1 {
		return Task{}, ErrInvalidUserID
	}

	if repoTask, err = s.repo.GetTaskByID(taskID); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}
	t = parseRepoTaskToTask(repoTask)

	if t.UserID != userID {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrNotOwner)
	}

	return t, nil
}
