package services_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
	"time"
)

// GetTasksForDay retreives the list of tasks with a dueDate specified by day
//
// Validation:
//   - ErrUserValidation {Field: UserFieldID}
//     userID negative or zero valued (< 1)
//   - ErrEmptyDay
//     day.IsZero() is true
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - returns the list of tasks with their fields populated
func TestGetTasksForDay(t *testing.T) {
	cases := []struct {
		name string

		paramUserID int
		paramDay    time.Time

		returnTasks []sqlite.Task
		returnErr   error

		wantedTasks []services.Task
		checkErr    func(*testing.T, error)
	}{

		// case: Empty UserID
		// expected: ErrUserValidation  {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Empty UserID",

			0,
			mustParseDate("2012-01-01"),

			[]sqlite.Task{{
				ID:       12,
				UserID:   2,
				Title:    "Empty UserID",
				Category: sql.NullString{String: "go tests", Valid: true},
			}},
			nil,

			nil,
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Negative UserID
		// expected: ErrUserValidation  {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Negative UserID",

			-1,
			mustParseDate("2012-01-01"),

			[]sqlite.Task{{
				ID:       12,
				UserID:   1,
				Title:    "Negative UserID",
				Category: sql.NullString{String: "go tests", Valid: true},
			}},
			nil,

			nil,
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Empty day
		// expected: ErrEmptyDay
		{
			"Empty Day",

			2,
			time.Time{},

			[]sqlite.Task{{
				ID:          12,
				UserID:      2,
				Title:       "Not for a day",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Shouldn't ever have tasks scheduled for the 70's", Valid: true},
			}},
			nil,

			nil,
			wantErrIs(services.ErrEmptyDay),
		},

		// case: UserID doesn't exist
		// expected: Empty list, nil error
		{
			"UserID Not Exist",

			32,
			mustParseDate("2000-01-01"),

			nil,
			nil,

			nil,
			wantNoErr,
		},

		// case: No tasks due on that day
		// expected: Empty list, nil error
		{
			"No Tasks On Day",

			2,
			mustParseDate("1995-01-01"),

			nil,
			nil,

			nil,
			wantNoErr,
		},

		// case: Repo returned ErrInternalRepo
		// expected: ErrInernalRepo (: err)
		{
			"Repo ErrInternalRepo",

			12,
			mustParseDate("2009-09-09"),

			nil,
			sqlite.ErrInternalRepo,

			nil,
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: List of tasks, nil error
		{
			"Happy Path",

			2,
			mustParseDate("2026-05-05"),

			[]sqlite.Task{{
				ID: 43,
				UserID: 2,
				Title: "Happy Title",
				Category: sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Happy description", Valid: true},
			}},
			nil,

			services.TaskList{{
				ID: 43,
				UserID: 2,
				Title: "Happy Title",
				Category:  "go tests",
				Description:  "Happy description",
			}},
			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockGetTasksForDay{
				stubRepo:    stubRepo{callCounts: make(map[string]int)},
				returnTasks: tc.returnTasks,
				returnErr:   tc.returnErr,
			}
			s := services.NewService(mock)

			gotTasks, gotErr := s.GetTasksForDay(tc.paramUserID, tc.paramDay)

			tc.checkErr(t, gotErr)

			if len(gotTasks) != len(tc.wantedTasks) {
				t.Errorf("wanted %d tasks, got %d", len(tc.wantedTasks), len(gotTasks))
			}
			for i := range gotTasks {
				checkTask(t, gotTasks[i], tc.wantedTasks[i])
			}
		})
	}
}
