package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
)

// GetOverdueTasks retrieves all incomplete tasks with a dueDate before time.Now()
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Empty:
//   - no overdue tasks
//   - userID doesn't exist
//
// Happy Path:
//   - returns all overdue tasks, or nothing if no tasks are overdue
func TestGetOverdueTasks(t *testing.T) {
	cases := []struct {
		name string

		paramUserID int

		wantTasks []sqlite.Task
		checkErr  func(*testing.T, error)
	}{

		// case: UserID that doesn't exist
		// expected: Empty list, nil error
		{
			"UserID Not Exist",

			89,

			[]sqlite.Task{},
			wantNoErr,
		},

		// case: UserID does exist, but no overdue tasks
		// expected: Empty list, nil error
		{
			"UserID Exists No Tasks",

			1,

			[]sqlite.Task{},
			wantNoErr,
		},

		// case: Happy path
		// expected: List of tasks, nil error
		{
			"Happy Path",

			3,

			[]sqlite.Task{{ID: 7, UserID: 3, Title: "Overdue Task", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "For checking overdue", Valid: true}, DueDate: sql.NullInt64{Int64: 1, Valid: true}}},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotTasks, gotErr := r.GetOverdueTasks(tc.paramUserID)

			// Check returned error
			tc.checkErr(t, gotErr)

			// Check returned tasks
			for i := range gotTasks {
				checkTask(t, tc.wantTasks[i], gotTasks[i])
			}
		})
	}
}
