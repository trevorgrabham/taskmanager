package db

import (
	"errors"
)

var (
	ErrInternalRepo           = errors.New("repo error")
	ErrNotOwner               = errors.New("user does not own task")
	ErrConstraintFailure      = errors.New("constraint failed")
	ErrTaskNotFound           = errors.New("task not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrUsernameTaken          = errors.New("username taken")
	ErrInvalidRecurringPeriod = errors.New("recurring period is not valid")
	ErrInvalidRecurringValue  = errors.New("recurring value is not valid")
	ErrInvalidRecurringUnit   = errors.New("recurring unit is not valid")
)

var (
	errInternalRepo           = "ErrInternalRepo"
	errNotOwner               = "ErrNotOwner"
	errTaskNotFound           = "ErrTaskNotFound"
	errUserNotFound           = "ErrUserNotFound"
	errConstraintFailure      = "ErrConstraintFailure"
	errUsernameTaken          = "ErrUsernameTaken"
	errInvalidRecurringPeriod = "ErrInvalidRecurringPeriod"
	errInvalidRecurringValue  = "ErrInvalidRecurringValue"
	errInvalidRecurringUnit   = "ErrInvalidRecurringUnit"
)

func MapErrorToTypeName(err error) string {
	if errors.Is(err, ErrInternalRepo) {
		return errInternalRepo
	} else if errors.Is(err, ErrNotOwner) {
		return errNotOwner
	} else if errors.Is(err, ErrConstraintFailure) {
		return errConstraintFailure
	} else if errors.Is(err, ErrTaskNotFound) {
		return errTaskNotFound
	} else if errors.Is(err, ErrUserNotFound) {
		return errUserNotFound
	} else if errors.Is(err, ErrInvalidRecurringPeriod) {
		return errInvalidRecurringPeriod
	} else if errors.Is(err, ErrInvalidRecurringValue) {
		return errInvalidRecurringValue
	} else if errors.Is(err, ErrInvalidRecurringUnit) {
		return errInvalidRecurringUnit
	} else if errors.Is(err, ErrUsernameTaken) {
		return errUsernameTaken
	}
	return ""
}
