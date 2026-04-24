package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
	"time"
)

func (s Service) UpdateDueDate(taskID, userID int, dueDay time.Time) (updatedTask Task, err error) {
	var (
		caller   = "UpdateDueDate"
		repoTask sqlite.Task
	)
	if taskID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskBadID)
	}
	if userID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}
	if dueDay.IsZero() {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrNoDate)
	}

	if repoTask, err = s.repo.TaskUpdateDueDateIfOwned(taskID, userID, dueDay); err != nil {
		if targetErr := (sqlite.ErrNotOwner{}); errors.Is(err, &targetErr) {
			return Task{}, fmt.Errorf("%s: %w", caller, err)
		}
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}

	updatedTask = parseRepoTaskToTask(repoTask)
	return updatedTask, nil
}
