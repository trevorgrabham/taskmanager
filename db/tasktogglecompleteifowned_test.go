package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
)

func TestToggleComplete(t *testing.T) {
	cases := []struct {
		name string

		paramTaskID int
		paramUserID int

		wantTask sqlite.Task
		checkErr func(*testing.T, error)
	}{

		// case: Empty TaskID
		// expected: No-op. nil error
		{
			"Empty TaskID",

			0,
			1,

			sqlite.Task{},
			wantNoErr,
		},

		// case: Empty UserID
		// expected: ErrNotOwner
		{
			"Empty UserID",

			3,
			0,

			sqlite.Task{},
			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: UserID not owner
		// expected: ErrNotOwner
		{
			"UserID Not Owner",

			8,
			3,

			sqlite.Task{},
			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: Task doesn't exist
		// expected: ErrTaskNotFound
		{
			"TaskID Not Exist",

			35,
			1,

			sqlite.Task{},
			wantNoErr,
		},

		// case: Task was completed
		// expected: Task now incomplete. nil error
		{
			"Task Already Completed",

			9,
			3,

			sqlite.Task{ID: 9, UserID: 3, Title: "A completed task", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "For testing with completed tasks", Valid: true}, DueDate: sql.NullInt64{Int64: 1, Valid: true}},
			wantNoErr,
		},

		// case: Task was incomplete
		// expected: Task now completed. nil error
		{
			"Task Incomplete",

			6,
			3,

			sqlite.Task{ID: 6, UserID: 3, Title: "Unscheduled Task", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "For checking unscheduled", Valid: true}, CompletionDate: sql.NullInt64{Int64: 1, Valid: true}, Done: true},
			wantNoErr,
		},

		// case: Uncompleted a recurring task
		// expected: Uncompleted returned task, recurring task that was previously inserted remains, nil error
		{
			"Uncompleted Recurring Task",

			4,
			3,

			sqlite.Task{ID: 4, UserID: 3, Title: "With RecurringPeriod", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "Seeing if the recurring_id works", Valid: true}, DueDate: sql.NullInt64{Int64: 1, Valid: true}, RecurringID: sql.NullInt64{Int64: 1, Valid: true}, RecurringPeriod: sql.NullString{String: "2 days", Valid: true}},
			wantNoErr,
		},

		// case: Completed a recurring task
		// expected: Completed task returned, and a new task inserted, nil error
		{
			"Completed Recurring Task",

			4,
			3,

			sqlite.Task{ID: 4, UserID: 3, Title: "With RecurringPeriod", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "Seeing if the recurring_id works", Valid: true}, DueDate: sql.NullInt64{Int64: 1, Valid: true}, RecurringID: sql.NullInt64{Int64: 1, Valid: true}, RecurringPeriod: sql.NullString{String: "2 days", Valid: true}, CompletionDate: sql.NullInt64{Int64: 1, Valid: true}, Done: true},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotTask, gotErr := r.TaskToggleCompleteIfOwned(tc.paramTaskID, tc.paramUserID)

			switch tc.name {
			case "Completed Recurring Task":
				tc.checkErr(t, gotErr)
				checkTask(t, tc.wantTask, gotTask)

				gotTask, gotErr = r.GetTaskByID(sqlite.NextID, tc.paramUserID)
				if gotErr != nil {
					t.Errorf("got error checking newly insterted task: %v", gotErr)
				}

				if gotTask == (sqlite.Task{}) {
					t.Errorf("no newly inserted task")
				}
			case "Uncompleted Recurring Task":
				gotTask, gotErr = r.TaskToggleCompleteIfOwned(gotTask.ID, gotTask.UserID)

				tc.checkErr(t, gotErr)
				checkTask(t, tc.wantTask, gotTask)
			default:
				tc.checkErr(t, gotErr)
				checkTask(t, tc.wantTask, gotTask)
			}
		})
	}
}
