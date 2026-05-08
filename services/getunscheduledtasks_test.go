package services_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// GetUnscheduledTasks retrieves the list of tasks without a dueDate
//
// Validation:
//   - ErrUserValidation {Field: UserFieldID}
//     userID negative or zero valued (< 1)
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - returns the list of unscheduled tasks for userID
func TestGetUnscheduledTasks(t *testing.T) {
	cases := []struct {
		name string

		paramUserID int

		returnTasks []sqlite.Task
		returnErr   error

		wantedTasks []services.Task
		checkErr    func(*testing.T, error)
	}{
		// case: Empty UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Empty UserID",

			0,

			[]sqlite.Task{{
				ID:       12,
				UserID:   2,
				Title:    "Test title",
				Category: sql.NullString{String: "go tests", Valid: true},
			}},
			nil,

			nil,
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Negative UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Negative UserID",

			-2,

			[]sqlite.Task{{
				ID:       12,
				UserID:   2,
				Title:    "Test title",
				Category: sql.NullString{String: "go tests", Valid: true},
			}},
			nil,

			nil,
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: UserID doesn't exist
		// expected: Empty list, nil error
		{
			"UserID Not Exist",

			1234,

			nil,
			nil,

			nil,
			wantNoErr,
		},

		// case: No unscheduled tasks
		// expected: Empty list, nil error
		{
			"No Unscheduled",

			1,

			nil,
			nil,

			nil,
			wantNoErr,
		},

		// case: Repo returned ErrInternalRepo
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			2,

			nil,
			sqlite.ErrInternalRepo,

			nil,
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: List of tasks, nil error
		{
			"Happy Path",

			3,

			[]sqlite.Task{{
				ID:     43,
				UserID: 3,
				Title:  "First Task",
			}, {
				ID:     57,
				UserID: 3,
				Title:  "Second Task",
			}},
			nil,

			services.TaskList{{
				ID:     43,
				UserID: 3,
				Title:  "First Task",
			}, {
				ID:     57,
				UserID: 3,
				Title:  "Second Task",
			}},
			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockGetOverdueTasks{
				stubRepo:    stubRepo{callCounts: make(map[string]int)},
				returnTasks: tc.returnTasks,
				returnErr:   tc.returnErr,
			}
			s := services.NewService(mock)

			gotTasks, gotErr := s.GetOverdueTasks(tc.paramUserID)

			tc.checkErr(t, gotErr)

			if len(gotTasks) != len(tc.wantedTasks) {
				t.Errorf("got %d tasks, wanted %d", len(gotTasks), len(tc.wantedTasks))
				return
			}
			for i := range gotTasks {
				checkTask(t, gotTasks[i], tc.wantedTasks[i])
			}
		})
	}
}
