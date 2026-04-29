package routes

import (
	"errors"
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"log"
	"net/http"
)

var (
	ErrInvalidSessionID  = errors.New("invalid session_id")
	ErrEmptySessionID    = errors.New("empty session_id")
	ErrParsingDay        = errors.New("day parsing error")
	ErrParsingDate       = errors.New("date parsing error")
	ErrParsingInt        = errors.New("int parsing error")
	ErrRenderingTemplate = errors.New("error rendering template")
	ErrInvalidDueDate    = errors.New("invalid due-date")
	ErrUnknownPath       = errors.New("unknown path")
	ErrParsingForm       = errors.New("request form parsing error")
	ErrEmptyUsername     = errors.New("empty username")
	ErrEmptyPassword     = errors.New("empty password")
	ErrEmptyConfirmPassword     = errors.New("empty confirm password")
	ErrUsernameTaken = errors.New("error username taken")
	ErrPasswordsNotMatch = errors.New("password and confirm password do not match")
	ErrEmptyTitle = errors.New("empty title")
	ErrWrongMethod       = errors.New("unsupported method")
	ErrInvalidRecurringPeriod = errors.New("invalid recurring period")
)

type ParseFormError struct {
	Caller string
	Err    error
}

func (e *ParseFormError) Error() string {
	return fmt.Sprintf("%s could not parse form: %v", e.Caller, e.Err)
}
func NewParseFormError(callingFunc string, e error) error {
	return &ParseFormError{Caller: callingFunc, Err: e}
}

type SessionIDCookieParseError struct {
	Caller string
	Err    error
}

func (e *SessionIDCookieParseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Caller, e.Err)
}
func NewSessionIDCookieParseError(callingFunc string, err error) error {
	return &SessionIDCookieParseError{Caller: callingFunc, Err: err}
}

type UnauthenticatedError struct {
	Caller string
}

func (e *UnauthenticatedError) Error() string {
	return fmt.Sprintf("%s: user unauthenticated", e.Caller)
}
func NewUnauthenticatedError(callingFunc string) error {
	return &UnauthenticatedError{Caller: callingFunc}
}

type DayParseError struct {
	Caller string
	Err    error
}

func (e *DayParseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Caller, e.Err)
}
func NewDayParseError(callingFunc string, err error) error {
	return &DayParseError{Caller: callingFunc, Err: err}
}

type IDParseError struct {
	Caller string
	Err    error
}

func (e *IDParseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Caller, e.Err)
}
func NewIDParseError(callingFunc string, err error) error {
	return &IDParseError{Caller: callingFunc, Err: err}
}

type NoIDError struct {
	Caller string
}

func (e *NoIDError) Error() string {
	return fmt.Sprintf("%s: no id", e.Caller)
}
func NewNoIDError(callingFunc string) error {
	return &NoIDError{Caller: callingFunc}
}

type DatabaseError struct {
	Caller string
	Err    error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("%s: database error: %v", e.Caller, e.Err)
}
func NewDatabaseError(callingFunc string, err error) error {
	return &IDParseError{Caller: callingFunc, Err: err}
}

type UnknownPathError struct {
	Caller string
	Path   string
}

func (e *UnknownPathError) Error() string {
	return fmt.Sprintf("%s: unknown path %s", e.Caller, e.Path)
}
func NewUnknownPathError(callingFunc, path string) error {
	return &UnknownPathError{Caller: callingFunc, Path: path}
}

type WrongMethodError struct {
	Caller   string
	Expected string
	Got      string
}

func (e *WrongMethodError) Error() string {
	return fmt.Sprintf("%s: expected: %s  got: %s", e.Caller, e.Expected, e.Got)
}
func NewWrongMethodError(callingFunc, expected, got string) error {
	return &WrongMethodError{Caller: callingFunc, Expected: expected, Got: got}
}

type NoUsernameError struct {
	Caller string
}

func (e *NoUsernameError) Error() string {
	return fmt.Sprintf("%s: no username", e.Caller)
}
func NewNoUsernameError(callingFunc string) error {
	return &NoUsernameError{Caller: callingFunc}
}

type NoConfirmPasswordError struct {
	Caller string
}

func (e *NoConfirmPasswordError) Error() string {
	return fmt.Sprintf("%s: no confirm password", e.Caller)
}
func NewNoConfirmPasswordError(callingFunc string) error {
	return &NoConfirmPasswordError{Caller: callingFunc}
}

type NoPasswordError struct {
	Caller string
}

func (e *NoPasswordError) Error() string {
	return fmt.Sprintf("%s: no password", e.Caller)
}
func NewNoPasswordError(callingFunc string) error {
	return &NoPasswordError{Caller: callingFunc}
}

type MismatchPasswordError struct {
	Caller          string
	Password        string
	ConfirmPassword string
}

func (e *MismatchPasswordError) Error() string {
	return fmt.Sprintf("%s: password mismatch. %s != %s", e.Caller, e.Password, e.ConfirmPassword)
}
func NewMismatchPasswordError(callingFunc, pass, confirm string) error {
	return &MismatchPasswordError{Caller: callingFunc, Password: pass, ConfirmPassword: confirm}
}

type NoTitleError struct {
	Caller string
}

func (e *NoTitleError) Error() string {
	return fmt.Sprintf("%s: no title", e.Caller)
}
func NewNoTitleError(callingFunc string) error {
	return &NoTitleError{Caller: callingFunc}
}

type DateTimeParseError struct {
	Caller string
	Value  string
	Err    error
}

func (e *DateTimeParseError) Error() string {
	return fmt.Sprintf("%s parsing %s: %v", e.Caller, e.Value, e.Err)
}
func NewDateTimeParseError(callingFunc, value string, err error) error {
	return &DateTimeParseError{Caller: callingFunc, Value: value, Err: err}
}

type RecurringParseError struct {
	Caller string
	Err    error
}

func (e *RecurringParseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Caller, e.Err)
}
func NewRecurringParseError(callingFunc string, err error) error {
	return &RecurringParseError{Caller: callingFunc, Err: err}
}

type TemplateRenderError struct {
	Caller   string
	Template string
	Err      error
}

func (e *TemplateRenderError) Error() string {
	return fmt.Sprintf("%s rendering %s: %v", e.Caller, e.Template, e.Err)
}
func NewTemplateRenderError(callingFunc, template string, err error) error {
	return &TemplateRenderError{Caller: callingFunc, Template: template, Err: err}
}

func (h Handler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, context.ErrWrongContextType) {
		// SessionID cookie exists, but is not a string
	} else if errors.Is(err, services.ErrInvalidSessionID) || errors.Is(err, ErrInvalidSessionID) {
		// User is unauthenticated
	} else if errors.Is(err, services.ErrInternalRepo) {
	} else if errors.Is(err, ErrParsingInt) {
	} else if errors.Is(err, services.ErrInvalidTaskID) {
		// taskID < 1
	} else if errors.Is(err, services.ErrNotOwner) {
		// task.UserID != userID
	}
	log.Println(err)
}
