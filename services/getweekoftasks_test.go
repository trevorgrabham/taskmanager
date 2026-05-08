package services_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
	"time"
)

// GetWeekOfTasks retrieves the list of tasks due in the range [day 00:00, day+6 23:59]
//
// Validation:
//   - ErrUserValidation {Field: UserFieldID}
//     userID negative or zero valued (< 1)
//		- ErrEmptyDay
//			day.IsZero() is true
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//		- returns the list of tasks for userID for the week
func TestGetWeekOfTasks(t *testing.T) {
	cases := []struct {
		name string

		paramUserID int
		paramDay time.Time

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
			time.Now(),

			[]sqlite.Task{{
				ID: 1,
				UserID: 1,
				Title: "First task",
				Category: sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "No userID", Valid: true},
				DueDate: sql.NullInt64{Int64: mustParseDate("2000-01-01").Unix(), Valid: true},
			}},
			nil,

			nil,
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

// case: Negative UserID
// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Negative UserID",

			0,
			time.Now(),

			[]sqlite.Task{{
				ID: 1,
				UserID: 1,
				Title: "The first task",
				Category: sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Negative userID", Valid: true},
				DueDate: sql.NullInt64{Int64: mustParseDate("2008-08-08").Unix(), Valid: true},
			}},
			nil,

			nil,
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

// case: Empty day
// expected: ErrEmptyDay 
		{
			"Empty Day",

			3,
			time.Time{},

			[]sqlite.Task{{
				ID: 15,
				UserID: 3,
				Title: "Task without a day",
				Category: sql.NullString{String: "go tests", Valid: true},
			}},
			nil,

			nil,
			wantErrIs(services.ErrEmptyDay),
		},

// case: UserID doesn't exist
// expected: Empty list, nil error
		{
			"UserID Not Exist",

			321,
			time.Now(),

			nil,
			nil,

			nil,
			wantNoErr,
		},

// case: No tasks for the week
// expected: Empty list, nil error
		{
			"No Tasks For Week",

			2,
			time.Now(),

			nil,
			nil,

			nil,
			wantNoErr,
		},

// case: Repo returned ErrInternalRepo
// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			3,
			time.Now(),

			nil,
			sqlite.ErrInternalRepo,

			nil,
			wantErrIs(services.ErrInternalRepo),
		},

// case: Happy path
// expected: List of tasks, nil error
		{
			"Happy Path",

			4,
			mustParseDate("2026-05-05"),

			[]sqlite.Task{{
				ID: 43,
				UserID: 4,
				Title: "First task this week",
				DueDate: sql.NullInt64{Int64: mustParseDate("2026-05-05").AddDate(0, 0, 3).Unix(), Valid: true},
			},{
				ID: 51,
				UserID: 4,
				Title: "Second task this week",
				DueDate: sql.NullInt64{Int64: mustParseDate("2026-05-05").AddDate(0, 0, 5).Unix(), Valid: true},
			}},
			nil,

			services.TaskList{{
				ID: 43,
				UserID: 4,
				Title: "First task this week",
				DueDate: mustParseDate("2026-05-05").AddDate(0, 0, 3),
			},{
				ID: 51,
				UserID: 4,
				Title: "Second task this week",
				DueDate:  mustParseDate("2026-05-05").AddDate(0, 0, 5),
			}},
			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockWeekOfTasks{
				stubRepo:    stubRepo{callCounts: make(map[string]int)},
				returnTasks: tc.returnTasks,
				returnErr:   tc.returnErr,
			}
			s := services.NewService(mock)

			gotTasks, gotErr := s.GetWeekOfTasks(tc.paramUserID, tc.paramDay)

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
