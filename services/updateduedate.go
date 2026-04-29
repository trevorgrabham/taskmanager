package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
	"time"
)

// UpdateDueDate updates the DueDate for the task identified by taskID.
//
// If taskID is empty, an ErrInvalidTaskID is returned. Consider using AddTask(Task) instead.
// If userID is empty, an ErrInvalidUserID is returned.
// If userID is not the owner of the task, an ErrNotOwner is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) UpdateDueDate(taskID, userID int, dueDay time.Time) (updatedTask Task, err error) {
	var (
		caller = "UpdateDueDate"
		repoTask sqlite.Task
	)
	if taskID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrInvalidTaskID)
	}
	if userID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrInvalidUserID)
	}

	if repoTask, err = s.repo.TaskUpdateDueDateIfOwned(taskID, userID, dueDay); err != nil {
		if errors.Is(err, sqlite.ErrNotOwner) { return Task{}, fmt.Errorf("%s: %w", caller, ErrNotOwner) }
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}

	updatedTask = parseRepoTaskToTask(repoTask)
	return updatedTask, nil
}
