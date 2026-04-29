package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// AddTask adds a new task to the Repo.
// 
// If task.ID is not empty, an ErrInvalidTaskID is returned. Consider using UpdateTask(Task) instead.
// If task.UserID is empty, an ErrInvalidUserID is returned.
// If task.Title is empty, an ErrInvalidTitle is returned.
// If task.RecurringPeriod is not a valid format, an ErrInvalidRecurringPeriod is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) AddTask(t Task) (Task, error) {
	var (
		caller    = "AddTask"
		err       error
		addedTask sqlite.Task
		repoTask  sqlite.Task
	)
	if t.ID > 0 { return Task{}, fmt.Errorf("%s: %w", caller, ErrInvalidTaskID) }
	if t.UserID < 1 { return Task{}, fmt.Errorf("%s: %w", caller, ErrInvalidUserID) }
	if t.Title == "" { return Task{}, fmt.Errorf("%s: %w", caller, ErrInvalidTitle) }
	// Consider error on t.Category == ""
	if err = validateRecurringPeriod(t.RecurringPeriod); err != nil { return Task{}, fmt.Errorf("%s: %w", caller, err) }

	repoTask = parseTaskToRepoTask(t)
	if addedTask, err = s.repo.AddTask(repoTask); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, err)
	}

	return parseRepoTaskToTask(addedTask), nil
}
