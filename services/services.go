// Package services provides business logic functions and their supporting types.
package services

import (
	sqlite "local/taskmanager/db"
	"time"
)

// Service is a wrapper around a database layer interface. It exposes methods to implement buisness logic.
type Service struct {
	repo Repo
}

// NewService initializes a Service using r to perform our database methods.
func NewService(r Repo) Service { return Service{repo: r} }

// Repo defines the database layer functions needed by Service to respond to requests.
type Repo interface {
	AddTask(sqlite.Task) (sqlite.Task, error)
	CompleteTaskIfOwned(taskID, userID int) error
	DeleteTaskIfOwned(taskID, userID int) error
	DeleteSession(sessionID string) error
	GetOverdueTasks(userID int) ([]sqlite.Task, error)
	GetTaskByID(taskID, userID int) (sqlite.Task, error)
	GetTasksForDay(userID int, day time.Time) ([]sqlite.Task, error)
	GetUnscheduledTasks(userID int) ([]sqlite.Task, error)
	GetUserBySessionID(sessionID string) (sqlite.User, error)
	GetUserByUsername(username string) (sqlite.User, error)
	GetUserCategorySuggestions(userID int) ([]string, error)
	GetWeekOfTasks(userID int, startDay time.Time) ([]sqlite.Task, error)
	SignupUserIfNotTaken(username, password string) (sqlite.User, error) // need to check that the username is not already taken
	StartSession(sessionID string, userID int) error
	TaskToggleCompleteIfOwned(taskID, userID int) error                              // if we uncomplete a recurring delete all matching future tasks (id > taskID)
	TaskUpdateDueDateIfOwned(taskID, userID int, day time.Time) (sqlite.Task, error) // need to make sure we get the old value, so we can keep the time of day the same. Only update the day, not the time
	UpdateTaskIfOwned(t sqlite.Task) (sqlite.Task, error)
}
