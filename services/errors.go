package services

import "errors"

/*
	NEED TO GO BACK AND CHANGE THE TESTS FOR THE UserID. I AM ACCESSING IT THROUGH THE SESSION ID, NEED TO ACCOUNT FOR THIS

	Consider creating an error struct for the login/signup pages.
*/

type TaskField string

const (
	TaskFieldTaskID          TaskField = "TaskID"
	TaskFieldUserID          TaskField = "UserID"
	TaskFieldTitle           TaskField = "Title"
	TaskFieldCategory        TaskField = "Category"
	TaskFieldDescription     TaskField = "Description"
	TaskFieldDueDate         TaskField = "DueDate"
	TaskFieldCompletionDate  TaskField = "CompletionDate"
	TaskFieldDone            TaskField = "Done"
	TaskFieldRecurringPeriod TaskField = "RecurringPeriod"
	TaskFieldRecurringValue  TaskField = "RecurringValue"
	TaskFieldRecurringUnit   TaskField = "RecurringUnit"
	TaskFieldRecurringID     TaskField = "RecurringID"
)

type UserField string

const (
	UserFieldID        UserField = "ID"
	UserFieldUsername  UserField = "Username"
	UserFieldPassword  UserField = "Password"
	UserFieldCreatedAt UserField = "CreatedAt"
	UserFieldUpdatedAt UserField = "UpdatedAt"
)

var (
	ErrInternalRepo    = errors.New("internal repo error")
	ErrEmptySessionID  = errors.New("empty session ID")
	ErrSessionNotFound = errors.New("session not found")
	ErrEmptyDay        = errors.New("empty day")
)

var (
	ErrInvalidUserID = errors.New("invalid userID")         // wrap this with the userID '%w %d'
	ErrNotOwner      = errors.New("user does not own task") // drop the returned error and wrap this one with the user and task ID's
	ErrInvalidDay    = errors.New("invalid day")            // drop the returned error and wrap this one with the day
	// ErrSessionExpired         = errors.New("session expired")        // drop the returned error and wrap this one with '%w for session %s'
	ErrEmptyUsername          = errors.New("empty username")
	ErrUsernameNotFound       = errors.New("username not found") // drop the returned error and wrap this one with '%w %s'
	ErrEmptyPassword          = errors.New("empty password")
	ErrWrongPassword          = errors.New("wrong password")   // wrap with '%w %s'
	ErrTaskIDNotEmpty         = errors.New("non-empty TaskID") // '%w'
	ErrInvalidID              = errors.New("negative ID")      // '%w'
	ErrEmptyTitle             = errors.New("empty title")
	ErrEmptyTaskID            = errors.New("empty taskID")
	ErrInvalidTaskID          = errors.New("invalid taskID")
	ErrTaskAlreadyComplete    = errors.New("task already complete") // optionally '%w on %s' for CompletionDate
	ErrTaskNotFound           = errors.New("task not found")        // '%w %d'
	ErrUserNotFound           = errors.New("user not found")        // '%w %d'
	ErrInvalidRecurringPeriod = errors.New("unknown RecurringPeriod format")
	ErrInvalidRecurringValue  = errors.New("invalid RecurringValue")
	ErrInvalidRecurringUnit   = errors.New("invalid RecurringUnit")
)

type ErrTaskValidation struct {
	Message string
	Field   TaskField
}

func (e *ErrTaskValidation) Error() string { return e.Message }
func AsErrTaskValidation(e error) (err *ErrTaskValidation) {
	errors.As(e, &err)
	return
}

type ErrUserValidation struct {
	Message string
	Field   UserField
}

func (e *ErrUserValidation) Error() string { return e.Message }
func AsErrUserValidation(e error) (err *ErrUserValidation) {
	errors.As(e, &err)
	return
}

const (
	MessageInvalidTaskID             = "invalid task ID"
	MessageInvalidUserID             = "invalid user ID"
	MessageEmptyTitle                = "title is a required field"
	MessageTaskAlreadyCompleted      = "task already completed"
	MessageTaskNotFound              = "task not found"
	MessageNotOwner                  = "user not the task owner"
	MessageInvalidUsernameOrPassword = "invalid username or password"
	MessageUsernameTaken             = "username already exists"
	MessageInvalidUsername           = "invalid username"
	MessageInvalidPassword           = "invalid password"
)
