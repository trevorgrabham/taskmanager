package services_test

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// CompleteTask completes the task if userID is the owner
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
//   - task is completed at the current date and time
func TestCompleteTask(t *testing.T) {
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
			2,

			nil,

			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

		// case: Negative TaskID
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageInvalidTaskID}
		{
			"Negative TaskID",

			-12,
			1,

			nil,

			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

		// case: Empty UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Empty UserID",

			12,
			0,

			nil,

			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Negative UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Negative UserID",

			3,
			-54,

			nil,

			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Repo returned ErrNotOwner
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageNotOwner}
		{
			"Repo ErrNotOwner",

			1,
			1,

			sqlite.ErrNotOwner,

			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageNotOwner),
		},

		// case: Repo returned ErrTaskNotFound
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageTaskNotFound}
		{
			"Repo ErrTaskNotFound",

			123,
			1,

			sqlite.ErrTaskNotFound,

			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageTaskNotFound),
		},

		// case: Happy path
		// expected: nil error
		{
			"Happy Path",

			11,
			2,

			nil,

			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockCompleteTask{
				stubRepo:  stubRepo{callCounts: make(map[string]int)},
				returnErr: tc.returnErr,
			}
			s := services.NewService(mock)

			gotErr := s.CompleteTask(tc.paramTaskID, tc.paramUserID)

			tc.checkErr(t, gotErr)
		})
	}
}
