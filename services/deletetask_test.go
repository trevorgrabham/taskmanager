package services_test

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// DeleteTask removes the task if userID is the owner
//
// Validation:
//   - ErrTaskValidation {Field: TaskFieldTaskID}
//     taskID negative or zero valued (< 1)
//     taskID doesn't belong to userID (repo: ErrNotOwner)
//   - ErrUserValidation {Field: UserFieldID}
//     userID negative or zero valued (< 1)
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - taskID is removed from the repo
func TestDeleteTask(t *testing.T) {
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

			-29,
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

			2,
			-32,

			nil,

			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Repo returned ErrNotOwner
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageNotOwner}
		{
			"Repo ErrNotOwner",

			12,
			2,

			sqlite.ErrNotOwner,

			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageNotOwner),
		},

		// case: Repo returned ErrInternalRepo
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			12,
			12,

			sqlite.ErrInternalRepo,

			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: nil
		{
			"Happy Path",

			12,
			3,

			nil,

			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockDeleteTask{
				stubRepo:  stubRepo{callCounts: make(map[string]int)},
				returnErr: tc.returnErr,
			}

			s := services.NewService(mock)
			gotErr := s.DeleteTask(tc.paramTaskID, tc.paramUserID)

			tc.checkErr(t, gotErr)
		})
	}
}
