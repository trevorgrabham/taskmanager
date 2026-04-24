package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

func (s Service) AddTask(t Task, userID int) (Task, error) {
	var (
		caller = "AddTask"
		err error 
		addedTask sqlite.Task
		repoTask sqlite.Task
	)
	if t.ID > 0 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskAlreadyExists)
	}
	if userID < 0 {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrUserBadID)
	}
	if t.Title == "" {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskNoTitle)
	}
	// Consider error on t.Category == ""

	t.UserID = userID
	repoTask = parseTaskToRepoTask(t)
	if addedTask, err = s.repo.AddTask(repoTask); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}

	return parseRepoTaskToTask(addedTask), nil
}
