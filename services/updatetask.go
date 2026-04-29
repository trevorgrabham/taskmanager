package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
)

// UpdateTask updates a task identified by task.ID.
//
// If task.ID is empty, an ErrInvalidTaskID is returned. Consider using AddTask(Task) instead.
// If task.UserID is empty, an ErrInvalidUserID is returned.
// If task.Title is empty, an ErrInvalidTitle is returned.
// If task.RecurringPeriod is not a valid format, an ErrInvalidRecurringPeriod is returned.
// If task.UserID is not the owner of the task, an ErrNotOwner is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) UpdateTask(t Task) (updatedTask Task, err error) {
	var (
		caller   = "UpdateTask"
		repoTask sqlite.Task
	)
	if t.ID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrInvalidTaskID)
	}
	if t.UserID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrInvalidUserID)
	}
	if t.Title == "" {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrInvalidTitle)
	}
	// Consider making Category required. If we do, check t.Category == "" here
	if err = validateRecurringPeriod(t.RecurringPeriod); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}

	repoTask = parseTaskToRepoTask(t)
	if repoTask, err = s.repo.UpdateTaskIfOwned(repoTask, t.UserID); err != nil {
		if errors.Is(err, sqlite.ErrNotOwner) { return Task{}, fmt.Errorf("%s: %w", ErrNotOwner) }
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}

	updatedTask = parseRepoTaskToTask(repoTask)
	return updatedTask, nil
}
