package services_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// GetTaskByID retreives the task identified by taskID if userID is the owner
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
//   - returns the task with all fields populated
func TestGetTaskByID(t *testing.T) {
	cases := []struct {
		name string

		paramTaskID int
		paramUserID int

		returnTask sqlite.Task
		returnErr  error

		wantedTask services.Task
		checkErr   func(*testing.T, error)
	}{
		// case: Empty TaskID
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageInvalidTaskID}
		{
			"Empty TaskID",

			0,
			1,

			sqlite.Task{
				ID:       4,
				UserID:   1,
				Title:    "Empty TaskID",
				Category: sql.NullString{String: "go tests", Valid: true},
				DueDate:  sql.NullInt64{Int64: mustParseDate("2012-12-12").Unix(), Valid: true},
			},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

		// case: Negative TaskID
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageInvalidTaskID}
		{
			"Negative TaskID",

			-2,
			1,

			sqlite.Task{
				ID:       4,
				UserID:   1,
				Title:    "Negative TaskID",
				Category: sql.NullString{String: "go tests", Valid: true},
				DueDate:  sql.NullInt64{Int64: mustParseDate("2012-12-12").Unix(), Valid: true},
			},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

		// case: Empty UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Empty UserID",

			12,
			0,

			sqlite.Task{
				ID:          12,
				UserID:      3,
				Title:       "Empty UserID",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "No user id at all", Valid: true},
			},
			nil,

			services.Task{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Negative UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Negative UserID",

			12,
			-3,

			sqlite.Task{
				ID:          12,
				UserID:      3,
				Title:       "Negative UserID",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "super negative", Valid: true},
			},
			nil,

			services.Task{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: TaskID doesn't exist
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageTaskNotFound}
		{
			"TaskID Not Exist",

			123,
			2,

			sqlite.Task{},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageTaskNotFound),
		},

		// case: UserID not the owner of TaskID
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageNotOwner}
		{
			"UserID Not Owner",

			12,
			3,

			sqlite.Task{},
			sqlite.ErrNotOwner, 

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageNotOwner),
		},

		// case: Repo returned ErrInternalRepo
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			1,
			1,

			sqlite.Task{},
			sqlite.ErrInternalRepo,

			services.Task{},
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: task with a matching id, non-empty title. Nil error
		{
			"Happy Path",

			12,
			3,

			sqlite.Task{
				ID: 12,
				UserID: 3,
				Title: "Happy happy path",
				Category: sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "The happy path should be happy", Valid: true},
			},
			nil,

			services.Task{
				ID: 12,
				UserID: 3,
				Title: "Happy happy path",
				Category:  "go tests",
				Description:  "The happy path should be happy",
			},
			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockGetTaskByID{
				stubRepo:   stubRepo{callCounts: make(map[string]int)},
				returnTask: tc.returnTask,
				returnErr:  tc.returnErr,
			}
			s := services.NewService(mock)

			gotTask, gotErr := s.GetTaskByID(tc.paramTaskID, tc.paramUserID)
			tc.checkErr(t, gotErr)

			checkTask(t, gotTask, tc.wantedTask)
		})
	}
}
