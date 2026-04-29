// Package services provides business logic functions and their supporting types.
package services

import (
	"errors"
	"fmt"
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
	GetTaskByID(taskID int) (sqlite.Task, error)
	GetTasksForDay(userID int, day time.Time) ([]sqlite.Task, error)
	GetUnscheduledTasks(userID int) ([]sqlite.Task, error)
	GetUserBySessionID(sessionID string) (sqlite.User, error)
	GetUserByUsername(username string) (sqlite.User, error)
	GetUserCategorySuggestions(userID int) ([]string, error)
	GetWeekOfTasks(userID int, startDay time.Time) ([]sqlite.Task, error)
	SignupUserIfNotTaken(username, password string) (sqlite.User, error) // need to check that the username is not already taken
	StartSession(sessionID string, userID int) error
	TaskToggleCompleteIfOwned(taskID, userID int) (sqlite.Task, error)               // if we uncomplete a recurring delete all matching future tasks (id > taskID)
	TaskUpdateDueDateIfOwned(taskID, userID int, day time.Time) (sqlite.Task, error) // need to make sure we get the old value, so we can keep the time of day the same. Only update the day, not the time
	UpdateTaskIfOwned(t sqlite.Task, userID int) (sqlite.Task, error)
}

var (
	ErrInvalidSessionID       = errors.New("no active sessions")
	ErrInternalRepo           = errors.New("repo error")
	ErrInvalidUserID          = errors.New("invalid userID")
	ErrInvalidTaskID          = errors.New("invalid taskID")
	ErrNotOwner               = errors.New("user does not own task")
	ErrInvalidUsername        = errors.New("invalid username")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrInvalidTitle           = errors.New("invalid title")
	ErrInvalidRecurringPeriod = errors.New("invalid recurring period")
	ErrInvalidRecurringValue  = errors.New("invalid recurring value")
	ErrInvalidRecurringUnit   = errors.New("invalid recurring unit")
	ErrInvalidDate            = errors.New("invalid date")
	ErrEncrypt                = errors.New("error during encryption")
	ErrUsernameExists         = errors.New("username already exists")
	// ErrNoDate              = errors.New("no date")
	// ErrSessionNoID         = errors.New("no session ID")
	// ErrTaskAlreadyExists   = errors.New("task already exists")
	// ErrTaskAlreadyComplete = errors.New("task already done")
	// ErrTaskBadID           = errors.New("bad task id")
	// ErrTaskNotExist        = errors.New("task does not exist")
	// ErrTaskNoTitle         = errors.New("no title for task")
	// ErrUserBadID           = errors.New("bad user id")
	// ErrUserNoPassword      = errors.New("no password")
	// ErrUserNoUsername      = errors.New("no username")
	// ErrUserWrongPassword   = errors.New("incorrect password")
	// ErrUserWrongID         = errors.New("task doesn't belong to user")
)

type ErrBcrypt struct {
	Err error
}

func (e *ErrBcrypt) Error() string { return fmt.Sprintf("bcrypt: %s", e.Err) }
func NewErrBcrypt(err error) error {
	return &ErrBcrypt{Err: err}
}
