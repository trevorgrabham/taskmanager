package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
	"time"
)

// GetWeekOfTasks retrieves the list of tasks with dueDate's in the range [day 00:00, day+6 23:59]
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Empty:
//   - userID doesn't exist
//   - day not initialized
//   - no tasks due
//
// Happy Path:
//   - returns the list of tasks due within a week from day
func TestGetWeekOfTasks(t *testing.T) {
	cases := []struct {
		name string

		paramUserID int
		paramDay    time.Time

		wantTasks []sqlite.Task
		checkErr  func(*testing.T, error)
	}{

		// case: UserID doesn't exist
		// expected: Empty list, nil error
		{
			"UserID Not Exist",

			-12,
			mustParseDate("2020-01-01"),

			nil,
			wantNoErr,
		},

		// case: Empty day
		// expected: Empty list, nil error
		{
			"Empty Day",

			1,
			time.Time{},

			nil,
			wantNoErr,
		},

		// case: No tasks for the week
		// expected: Empty list, nil error
		{
			"No Tasks For Week",

			1,
			mustParseDate("2006-06-14"),

			[]sqlite.Task{},
			wantNoErr,
		},

		// case: Happy path
		// expected: List of tasks, nil error
		{
			"Happy Path",

			1,
			time.Now().AddDate(0, 0, 1),

			[]sqlite.Task{{ID: 1, UserID: 1, Title: "First Task", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "First task's description", Valid: true}, DueDate: sql.NullInt64{Int64: 1, Valid: true}}, {ID: 3, UserID: 1, Title: "Empty Description", Category: sql.NullString{String: "general tests", Valid: true}, DueDate: sql.NullInt64{Int64: 1, Valid: true}}},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotTasks, gotErr := r.GetWeekOfTasks(tc.paramUserID, tc.paramDay)

			tc.checkErr(t, gotErr)

			if len(gotTasks) != len(tc.wantTasks) {
				t.Errorf("wanted tasks: %v, got %v", tc.wantTasks, gotTasks)
				return
			}
			for i := range gotTasks {
				checkTask(t, tc.wantTasks[i], gotTasks[i])
			}

		})
	}
}
