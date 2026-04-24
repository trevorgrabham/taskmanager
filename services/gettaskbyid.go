package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

func (s Service) GetTaskByID(taskID, userID int) (t Task, err error) {
	var (
		caller   = "GetTaskByID"
		repoTask sqlite.Task
	)
	if taskID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskBadID)
	}
	if userID < 1 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}

	if repoTask, err = s.repo.GetTaskByID(taskID); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}
	t = parseRepoTaskToTask(repoTask)

	if t.IsZero() {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskNotExist)
	}
	if t.UserID != userID {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrUserWrongID)
	}

	return t, nil
}
