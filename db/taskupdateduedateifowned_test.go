package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
	"time"
)

// TaskUpdateDueDateIfOwned changes the dueDate for the task and returns the updated task
//
// Errors:
//		- ErrNotOwner 
//			userID doesn't exist 
//			userID not the owner
//		- ErrInternalRepo
//			transient database error
//
// Empty:
//   - taskID doesn't exist
//
// Happy Path:
//   - returns the updated task with its dueDate updated (or removed if day was empty)
func TestUpdateDueDate(t *testing.T) {
	cases := []struct {
		name string

		paramTaskID int
		paramUserID int
		paramDay    time.Time

		wantTask sqlite.Task
		checkErr func(*testing.T, error)
	}{

		// case: Task not exist
		// expected: No-op, nil error
		{
			"TaskID Not Exist",

			0,
			1,
			mustParseDate("2027-02-07"),

			sqlite.Task{},
			wantNoErr,
		},

		// case: Task not exist
		// expected: No-op, nil error
		{
			"Unscheduled TaskID Not Exist",

			0,
			1,
			time.Time{},

			sqlite.Task{},
			wantNoErr,
		},

		// case: UserID doesn't exist
		// expected: ErrNotOwner
		{
			"UserID Not Exist",

			1,
			12,
			mustParseDate("2027-02-07"),

			sqlite.Task{},
			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: Day is empty
		// expected: Due date removed, nil error
		{
			"Empty Day",

			1,
			1,
			time.Time{},

			sqlite.Task{ID: 1, UserID: 1, Title: "First Task", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "First task's description", Valid: true}},
			wantNoErr,
		},

		// case: UserID not the owner
		// expected: ErrNotOwner
		{
			"UserID Not Owner",

			1,
			2,
			mustParseDate("2027-02-07"),

			sqlite.Task{},
			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: Happy path
		// expected: Due date updated, nil error
		{
			"Happy Path",

			6,
			3,
			mustParseDate("2012-12-12"),

			sqlite.Task{ID: 6, UserID: 3, Title: "Unscheduled Task", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "For checking unscheduled", Valid: true}, DueDate: sql.NullInt64{Int64: 1, Valid: true}},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotTask, gotErr := r.TaskUpdateDueDateIfOwned(tc.paramTaskID, tc.paramUserID, tc.paramDay)

			checkTask(t, tc.wantTask, gotTask)

			tc.checkErr(t, gotErr)
		})
	}
}
