package db_test

import (
	sqlite "local/taskmanager/db"
	"testing"
)

// TaskToggleCompleteIfOwned toggles the completion status of task if userID is the owner
//
// Errors:
//   - ErrNotOwner
//     userID doesn't exist
//     userID not the owner
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - toggles the completion status of task. If task was recurring and was just completed, adds the next task to the database
func TestToggleComplete(t *testing.T) {
	cases := []struct {
		name string

		paramTaskID int
		paramUserID int

		wantComplete bool
		checkErr     func(*testing.T, error)
	}{

		// case: Empty TaskID
		// expected: No-op. nil error
		{
			"Empty TaskID",

			0,
			1,

			false,
			wantNoErr,
		},

		// case: Empty UserID
		// expected: ErrNotOwner
		{
			"Empty UserID",

			3,
			0,

			false,
			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: UserID not owner
		// expected: ErrNotOwner
		{
			"UserID Not Owner",

			8,
			3,

			false,
			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: Task doesn't exist
		// expected: No-op, nil error
		{
			"TaskID Not Exist",

			35,
			1,

			false,
			wantNoErr,
		},

		// case: Task was completed
		// expected: Task now incomplete. nil error
		{
			"Task Already Completed",

			9,
			3,

			false,
			wantNoErr,
		},

		// case: Task was incomplete
		// expected: Task now completed. nil error
		{
			"Task Incomplete",

			6,
			3,

			true,
			wantNoErr,
		},

		// case: Completed a recurring task
		// expected: Task now completed, and a new task inserted, nil error
		{
			"Completed Recurring Task",

			4,
			3,

			true,
			wantNoErr,
		},

		// case: Uncompleted a recurring task
		// expected: Task now incomplete, recurring task that was previously inserted remains, nil error
		{
			"Uncompleted Recurring Task",

			9,
			3,

			false,
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotErr := r.TaskToggleCompleteIfOwned(tc.paramTaskID, tc.paramUserID)

			tc.checkErr(t, gotErr)
			if gotErr != nil {
				return
			}

			gotTask, gotErr := r.GetTaskByID(tc.paramTaskID, tc.paramUserID)
			if gotErr != nil {
				t.Errorf("getting task got err %v", gotErr)
				return
			}

			if tc.wantComplete != gotTask.Done {
				t.Errorf("wanted Complete: %t, got %t", tc.wantComplete, gotTask.Done)
			}

			if tc.name == "Completed Recurring Task" {
				gotTask, gotErr = r.GetTaskByID(sqlite.NextID, tc.paramUserID)
				if gotErr != nil {
					t.Errorf("getting task got err %v", gotErr)
					return
				}

				if gotTask == (sqlite.Task{}) {
					t.Errorf("no new recurring task inserted")
				}
			}
		})
	}
}
