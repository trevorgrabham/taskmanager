package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
	"time"
)

func TestGetTasksForDay(t *testing.T) {
	cases := []struct {
		name string

		paramUserID int
		paramDay    time.Time

		wantTasks []sqlite.Task
		checkErr  func(*testing.T, error)
	}{

		// case: UserID doesn't exist
		// expected: nil tasks, nil error 
		{
			"UserID Not Exist",

			-34,
			mustParseDate("2000-01-01"),

			nil,
			wantNoErr,
		},

		// case: Empty day
		// expected: nil tasks, nil error
		{
			"Empty Day",

			2, 
			time.Time{},

			nil,
			wantNoErr,
		},

		// case: No tasks due on that day
		// expected: Empty list, nil error
		{
			"No Tasks Due",

			2, 
			time.Now().AddDate(0, 0, 1),

			nil,
			wantNoErr,
		},

		// case: Happy path
		// expected: List of tasks, nil error
		{
			"Happy Path",

			2,
			time.Now().AddDate(0, 0, 3),

			[]sqlite.Task{{ID: 2, UserID: 2, Title: "Second Task", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "Another task description", Valid: true}, DueDate: sql.NullInt64{Int64: time.Now().AddDate(0,0,3).Unix(), Valid: true}}},
			wantNoErr,
		}}
	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotTasks, gotErr := r.GetTasksForDay(tc.paramUserID, tc.paramDay)

			for i := range gotTasks {
				checkTask(t, tc.wantTasks[i], gotTasks[i])
			}

			tc.checkErr(t, gotErr)
		})
	}
}
