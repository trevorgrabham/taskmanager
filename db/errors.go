package db

import (
	"errors"
	"fmt"
)

var (
	ErrNotConnected    = errors.New("repo not connected")
	ErrTaskAlreadyDone = errors.New("task already complete")
	ErrNilSQLRow       = errors.New("row is nil")
	ErrUnknown = errors.New("unknown error")
	ErrUsernameTaken = errors.New("username already taken")
	ErrTaskBadRecurringPeriod = errors.New("badly formed recurring period")
)

type ErrTransactionCommit struct {
	Err error
}

func (e *ErrTransactionCommit) Error() string { return fmt.Sprintf("commit error: %s", e.Err) }
func NewErrTransactionCommit(err error) error { return &ErrTransactionCommit{Err: err} }

type ErrRepo struct {
	Err error
}

func (e *ErrRepo) Error() string { return fmt.Sprintf("internal repo error: %s", e.Err) }
func NewErrRepo(err error) error {
	return &ErrRepo{Err: err}
}

type ErrTaskNotExist struct {
	TaskID int
}
func (e *ErrTaskNotExist) Error() string { return fmt.Sprintf("taskID %d does not exist", e.TaskID) }
func NewErrTaskNotExist(taskID int) error { return &ErrTaskNotExist{TaskID: taskID} }

type ErrNotOwner struct {
	UserID int
	TaskID int
}
func (e *ErrNotOwner) Error() string { return fmt.Sprintf("user #%d is not the owner of task #%d", e.UserID, e.TaskID) }
func NewErrNotOwner(userID, taskID int) error { return &ErrNotOwner{UserID: userID, TaskID: taskID} }
