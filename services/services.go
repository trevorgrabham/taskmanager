package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
	"time"
)

type Repo interface {
	AddTask(sqlite.Task) (sqlite.Task, error)
	CompleteTaskIfOwned(taskID, userID int) error
	DeleteTaskIfOwned(taskID, userID int) error
	GetOverdueTasks(userID int) ([]sqlite.Task, error)
	GetTaskByID(taskID int) (sqlite.Task, error)
	GetTasksForDay(userID int, day time.Time) ([]sqlite.Task, error)
	GetUnscheduledTasks(userID int) ([]sqlite.Task, error)
	GetUserBySessionID(sessionID string) (sqlite.User, error)
	GetUserByUsername(username string) (sqlite.User, error)
	GetUserCategorySuggestions(userID int) ([]string, error)
	GetWeekOfTasks(userID int, startDay time.Time) ([]sqlite.Task, error)
	SignupUserIfNotTaken(username, password string) (sqlite.User, error) // need to check that the username is not already taken
	StartSession(userID int) (sessionID string, err error)
	TaskToggleCompleteIfOwned(taskID, userID int) (sqlite.Task, error)               // if we uncomplete a recurring delete all matching future tasks (id > taskID)
	TaskUpdateDueDateIfOwned(taskID, userID int, day time.Time) (sqlite.Task, error) // need to make sure we get the old value, so we can keep the time of day the same. Only update the day, not the time
	UpdateTaskIfOwned(t sqlite.Task, userID int) (sqlite.Task, error)
}

type Service struct {
	repo Repo
}

func NewService(r Repo) Service { return Service{repo: r} }

var (
	ErrNoDate              = errors.New("no date")
	ErrSessionNoID         = errors.New("no session ID")
	ErrTaskAlreadyExists   = errors.New("task already exists")
	ErrTaskAlreadyComplete = errors.New("task already done")
	ErrTaskBadID           = errors.New("bad task id")
	ErrTaskNotExist        = errors.New("task does not exist")
	ErrTaskNoTitle         = errors.New("no title for task")
	ErrUserBadID           = errors.New("bad user id")
	ErrUserNoPassword      = errors.New("no password")
	ErrUserNoUsername      = errors.New("no username")
	ErrUserWrongPassword   = errors.New("incorrect password")
	ErrUserWrongID         = errors.New("task doesn't belong to user")
)

type ErrBcrypt struct {
	Err error
}

func (e *ErrBcrypt) Error() string { return fmt.Sprintf("bcrypt: %s", e.Err) }
func NewErrBcrypt(err error) error {
	return &ErrBcrypt{Err: err}
}
