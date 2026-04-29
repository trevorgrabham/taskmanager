package services

import "errors"

/*
	NEED TO GO BACK AND CHANGE THE TESTS FOR THE UserID. I AM ACCESSING IT THROUGH THE SESSION ID, NEED TO ACCOUNT FOR THIS

	Consider creating an error struct for the login/signup pages.
*/

type taskField string
const (
	TaskID taskField = "taskID"
	UserID taskField = "userID"
	Title taskField = "title"
	Category taskField = "category"
	Description taskField = "description"
	DueDate taskField = "dueDate"
	CompletionDate taskField = "completionDate"
	Done taskField = "done"
	RecurringPeriod taskField = "recurringPeriod"
	RecurringID taskField  = "recurringID"
)

var (
	// ErrInvalidUserID = errors.New("invalid userID") // wrap this with the userID '%w %d'
	// ErrInternalRepo = errors.New("service repo error") // wrap around the returned error
	// ErrNotOwner = errors.New("user does not own task") 	// drop the returned error and wrap this one with the user and task ID's
	// ErrInvalidDay = errors.New("invalid day") // drop the returned error and wrap this one with the day
	// ErrEmptySessionID = errors.New("empty sessionID") // drop the returned error
	// ErrSessionNotFound = errors.New("session not found") // wrap with the sessionID '%w %s'
	// ErrSessionExpired = errors.New("session expired") // drop the returned error and wrap this one with '%w for session %s'
	// ErrEmptyUsername = errors.New("empty username") 
	// ErrUsernameNotFound = errors.New("username not found") // drop the returned error and wrap this one with '%w %s'
	// ErrEmptyPassword = errors.New("empty password") 
	// ErrWrongPassword = errors.New("wrong password") // wrap with '%w %s'
)

type ErrInvalidTask interface {
	Error() string
	Field() taskField
}
type errInvalidTask struct {
	err error 
	field taskField
}
func (e *errInvalidTask) Error() string { return e.err.Error() }
func (e *errInvalidTask) Field() taskField { return e.field }
func IsErrInvalidTask(e error) (err ErrInvalidTask) { 
	errors.As(e, &err)
	return
}
