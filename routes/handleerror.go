package routes

import (
	"fmt"
	"net/http"
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
}

func (e *UnknownPathError) Error() string {
	return fmt.Sprintf("%s: uknown path", e.Caller)
}
func NewUnknownPathError(callingFunc string) error {
	return &UnknownPathError{Caller: callingFunc}
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
	Caller string
	Password string 
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
	Value string
	Err error
}

func (e *DateTimeParseError) Error() string {
	return fmt.Sprintf("%s parsing %s: %v", e.Caller, e.Value, e.Err)
}
func NewDateTimeParseError(callingFunc, value string, err error) error {
	return &DateTimeParseError{Caller: callingFunc, Value: value, Err: err}
}

type RecurringParseError struct {
	Caller string
	Err error
}

func (e *RecurringParseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Caller, e.Err)
}
func NewRecurringParseError(callingFunc string, err error) error {
	return &RecurringParseError{Caller: callingFunc, Err: err}
}

type TemplateRenderError struct {
	Caller string
	Template string
	Err error
}

func (e *TemplateRenderError) Error() string {
	return fmt.Sprintf("%s rendering %s: %v", e.Caller, e.Template, e.Err)
}
func NewTemplateRenderError(callingFunc, template string, err error) error {
	return &TemplateRenderError{Caller: callingFunc, Template: template, Err: err}
}

func (h Handlers) HandleError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *ParseFormError:
	// Unrecoverable. Log error and let user know
	case *SessionIDCookieParseError:
	// Recoverable. Signal user to delete sessionID cookie and redirect to login page
	case *UnauthenticatedError:
		// Recoverable. Redirect to login page
	case *DayParseError:
	// Depending on the caller, we can do different things.
	// If it was from a form, use an oob swap to return the date input with .error attached.
	// If it was from a drag and drop, we need to think about wether the element was removed from the list, and if it needs to be re-added. If the javascript waits for a response then we can just send an error
	case *RecurringParseError:
	// If it was from a form, use an oob swap to return the date input with .error attached.
	case *IDParseError:
	// Unrecoverable. Log error and let user know
	case *NoIDError:
	// Unrecoverable. Log error and let user know
	case *DatabaseError:
	// Unrecoverable. Log error and let user know
	case *UnknownPathError:
	// Unrecoverable. Log error and let user know
	case *WrongMethodError:
	// Unrecoverable. Log error and let user know
	case *NoConfirmPasswordError:
	// Recoverable. Should be from the login or signup page, so we can use an oob-swap with a .error class to highlight the error
	case *NoPasswordError:
	// Recoverable. Should be from the login or signup page, so we can use an oob-swap with a .error class to highlight the error
	case *MismatchPasswordError:
	// Recoverable. Should be from the login or signup page, so we can use an oob-swap with a .error class to highlight the error
	case *NoUsernameError:
		// Recoverable. Should be from the login or signup page, so we can use an oob-swap with a .error class to highlight the error
	case *NoTitleError:
	// Unrecoverable. Log error and let user know
	case *DateTimeParseError:
	// Unrecoverable. Log error and let user know
	}
}
