package services_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// GetOverdueTasks retrieves the list of tasks that have a dueDate before time.Now()
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
//   - returns the list of overdue tasks for userID
func TestGetOverdueTasks(t *testing.T) {
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
				ID:          12,
				UserID:      2,
				Title:       "Shouldn't see this",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Really shouldn't see this", Valid: true},
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
				ID:          12,
				UserID:      2,
				Title:       "Shouldn't see this",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Really shouldn't see this", Valid: true},
			}},
			nil,

			nil,
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: UserID exists, but no overdue tasks
		// expected: Empty list, nil error
		{
			"UserID Not Exist",

			3,

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
		// expected: List of overdue tasks, nil error
		{
			"Happy Path",

			2,

			[]sqlite.Task{{
				ID:              12,
				UserID:          2,
				Title:           "Now we should see this",
				Category:        sql.NullString{String: "go tests", Valid: true},
				Description:     sql.NullString{String: "We should be seeing all of this", Valid: true},
				DueDate:         sql.NullInt64{Int64: mustParseDate("2020-02-02").Unix(), Valid: true},
				RecurringPeriod: sql.NullInt64{Int64: 3600 * 24, Valid: true},
			}},
			nil,

			services.TaskList{{
				ID:              12,
				UserID:          2,
				Title:           "Now we should see this",
				Category:        "go tests",
				Description:     "We should be seeing all of this",
				DueDate:         mustParseDate("2020-02-02"),
				RecurringPeriod: 3600 * 24,
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
			}
			for i := range gotTasks {
				checkTask(t, gotTasks[i], tc.wantedTasks[i])
			}
		})
	}
}
