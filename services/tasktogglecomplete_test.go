package services_test

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// TaskToggleComplete toggles the completion status of taskID if userID is the owner. If the task is recurring and has just been completed, a new recurring task is created
//
// Validation:
//   - ErrTaskValidation {Field: TaskFieldTaskID}
//     taskID negative or zero valued (< 1)
//     taskID doesn't belong to userID (repo: ErrNotOwner)
//     taskID doesn't exist
//   - ErrUserValidation {Field: UserFieldID}
//     userID negative or zero valued (< 1)
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//		- task is toggled and a new task is created if we just completed a recurring task
func TestToggleComplete(t *testing.T) {
	cases := []struct {
		name string

		paramTaskID int
		paramUserID int

		returnErr error

		checkErr func(*testing.T, error)
	}{

// case: Empty TaskID 
// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageInvalidTaskID}
		{
			"Empty TaskID",

			0,
			3,

			nil,

			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

// case: Negative TaskID 
// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageInvalidTaskID}
		{
			"Negative TaskID",

			-32,
			2,

			nil,

			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

// case: Empty UserID 
// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Empty UserID",

			2,
			0,

			nil, 

			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

// case: Negative UserID
// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Negative UserID",

			22,
			-32,

			nil, 

			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

// case: UserID not the owner
// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageNotOwner}
		{
			"UserID Not Owner",

			3,
			1,

			sqlite.ErrNotOwner,

			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageNotOwner),
		},

// case: TaskID doesn't exist
// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageTaskNotFound}
		{
			"TaskID Not Exist",

			439,
			5,

			sqlite.ErrTaskNotFound,

			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageTaskNotFound),
		},

// case: Repo returned ErrInternalRepo
// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			15,
			3,

			sqlite.ErrInternalRepo,

			wantErrIs(services.ErrInternalRepo),
		},

// case: Happy path
// expected: Task completion toggled, nil error
		{
			"Happy Path",

			32,
			6,

			nil,

			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockTaskToggleComplete{
				stubRepo:  stubRepo{callCounts: make(map[string]int)},
				returnErr: tc.returnErr,
			}
			s := services.NewService(mock)

			gotErr := s.TaskToggleComplete(tc.paramTaskID, tc.paramUserID)

			tc.checkErr(t, gotErr)
		})
	}
}
