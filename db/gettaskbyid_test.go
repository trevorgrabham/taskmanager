package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
)

func TestGetTaskByID(t *testing.T) {
	cases := []struct {
		name string

		paramTaskID int
		paramUserID int

		wantTask sqlite.Task
		checkErr func(*testing.T, error)
	}{

		// case: UserID not owner
		// expected: ErrNotOwner
		{
			"UserID Not Owner",

			1,
			3,

			sqlite.Task{},
			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: UserID doesn't exist
		// expected: ErrNotOwner
		{
			"UserID Not Exist",

			3,
			12,

			sqlite.Task{},
			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: Task doesn't exist
		// expected: Empty task, nil error
		{
			"TaskID Not Exist",

			sqlite.NextID + 12,
			1,

			sqlite.Task{},
			wantNoErr,
		},

		// case: Happy path
		// expected: task with a matching id, non-empty title. Nil error
		{
			"Happy Path",

			6,
			3,

			sqlite.Task{ID: 6, UserID: 3, Title: "Unscheduled Task", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "For checking unscheduled", Valid: true}},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotTask, gotErr := r.GetTaskByID(tc.paramTaskID, tc.paramUserID)

			checkTask(t, tc.wantTask, gotTask)

			tc.checkErr(t, gotErr)
		})
	}
}
