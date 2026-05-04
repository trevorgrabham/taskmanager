package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
)

func TestGetUnscheduledTasks(t *testing.T) {
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

			0,

			nil,
			wantNoErr,
		},

// case: No unscheduled tasks
// expected: Empty list, nil error
		{
			"No Unscheduled Tasks",

			2,

			nil, 
			wantNoErr,
		},

// case: Happy path
// expected: List of tasks, nil error
		{
			"Happy Path",

			3,

			[]sqlite.Task{{ID: 6, UserID: 3, Title: "Unscheduled Task", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "For checking unscheduled", Valid: true}}},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotTasks, gotErr := r.GetUnscheduledTasks(tc.paramUserID)

			for i := range gotTasks {
				checkTask(t, tc.wantTasks[i], gotTasks[i])
			}

			tc.checkErr(t, gotErr)
		})
	}
}
