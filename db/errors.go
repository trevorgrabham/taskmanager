package db

import (
	"errors"
)

var (
	ErrNoLiveSession          = errors.New("no live session")
	ErrTaskNotExist           = errors.New("task does not exits")
	ErrUserNotExist           = errors.New("user does not exits")
	ErrUnknown                = errors.New("unknown error")
	ErrEmptyDate              = errors.New("empty date")
	ErrTransactionCommit      = errors.New("error commiting transaction")
	ErrInvalidTask            = errors.New("invalid task")
	ErrTaskCompleted          = errors.New("task already completed")
	ErrUsernameExists         = errors.New("username exists already")
	ErrInvalidUsername        = errors.New("username is not valid")
	ErrInvalidPassword        = errors.New("password is not valid")
	ErrInvalidRecurringPeriod = errors.New("recurring period is not valid")
	ErrInvalidRecurringValue  = errors.New("recurring value is not valid")
	ErrInvalidRecurringUnit   = errors.New("recurring unit is not valid")
)

var (
	ErrNotConnected         = errors.New("repo not connected")
	ErrInvalidTaskID        = errors.New("invalid taskID")
	ErrInvalidUserID        = errors.New("invalid userID")
	ErrEmptyTitle           = errors.New("empty title")
	ErrInternalRepo         = errors.New("repo error")
	ErrNotOwner             = errors.New("user does not own task")
	ErrTaskNotFound         = errors.New("task not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrAlreadyCompleted     = errors.New("already completed")
	ErrEmptySessionID       = errors.New("empty sessionID")
	ErrSessionIDNotFound    = errors.New("sessionID not found")
	ErrSessionExpired       = errors.New("session expired")
	ErrSessionAlreadyExists = errors.New("session already exists")
	ErrInvalidDay           = errors.New("invalid day")
	ErrEmptyUsername        = errors.New("empty username")
	ErrUsernameNotFound     = errors.New("username not found")
	ErrUsernameTaken        = errors.New("username taken")
	ErrEmptyPassword        = errors.New("empty password")
)

var (
	errNotConnected         = "ErrNotConnected"
	errInvalidTaskID        = "ErrInvalidTaskID"
	errInvalidUserID        = "ErrInvalidUserID"
	errEmptyTitle           = "ErrEmptyTitle"
	errInternalRepo         = "ErrInternalRepo"
	errNotOwner             = "ErrNotOwner"
	errTaskNotFound         = "ErrTaskNotFound"
	errUserNotFound         = "ErrUserNotFound"
	errAlreadyCompleted     = "ErrAlreadyCompleted"
	errEmptySessionID       = "ErrEmptySessionID"
	errSessionIDNotFound    = "ErrSessionIDNotFound"
	errSessionExpired       = "ErrSessionExpired"
	errSessionAlreadyExists = "ErrSessionAlreadyExists"
	errInvalidDay           = "ErrInvalidDay"
	errEmptyUsername        = "ErrEmptyUsername"
	errUsernameNotFound     = "ErrUsernameNotFound"
	errUsernameTaken        = "ErrUsernameTaken"
	errEmptyPassword        = "ErrEmptyPassword"
)

func MapErrorToTypeName(err error) string {
	if errors.Is(err, ErrNotConnected) {
		return errNotConnected
	} else if errors.Is(err, ErrInvalidTaskID) {
		return errInvalidTaskID
	} else if errors.Is(err, ErrInvalidUserID) {
		return errInvalidUserID
	} else if errors.Is(err, ErrEmptyTitle) {
		return errEmptyTitle
	} else if errors.Is(err, ErrInternalRepo) {
		return errInternalRepo
	} else if errors.Is(err, ErrNotOwner) {
		return errNotOwner
	} else if errors.Is(err, ErrTaskNotFound) {
		return errTaskNotFound
	} else if errors.Is(err, ErrUserNotFound) {
		return errUserNotFound
	} else if errors.Is(err, ErrAlreadyCompleted) {
		return errAlreadyCompleted
	} else if errors.Is(err, ErrEmptySessionID) {
		return errEmptySessionID
	} else if errors.Is(err, ErrSessionIDNotFound) {
		return errSessionIDNotFound
	} else if errors.Is(err, ErrSessionExpired) {
		return errSessionExpired
	} else if errors.Is(err, ErrSessionAlreadyExists) {
		return errSessionAlreadyExists
	} else if errors.Is(err, ErrInvalidDay) {
		return errInvalidDay
	} else if errors.Is(err, ErrEmptyUsername) {
		return errEmptyUsername
	} else if errors.Is(err, ErrUsernameNotFound) {
		return errUsernameNotFound
	} else if errors.Is(err, ErrUsernameTaken) {
		return errUsernameTaken
	} else if errors.Is(err, ErrEmptyPassword) {
		return errEmptyPassword
	}
	return ""
}
