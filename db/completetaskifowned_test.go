package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
)

// CompleteTask marks the task identified by taskID as completed at the current date and time, if it is owned by userID.
//
// Errors:
//   - ErrTaskNotFound
//     taskID doesn't exist
//   - ErrNotOwner
//     task exists, but userID is not the owner (or doesn't exist)
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - updates the task as completed at time.Now()
func TestCompleteTask(t *testing.T) {
	cases := []struct {
		name string

		paramTaskID int
		paramUserID int

		checkErr func(*testing.T, error)
	}{

		// case: Empty UserID
		// expected: ErrNotOwner
		{
			"EmptyUserID",

			1,
			0,

			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: UserID not the owner of taskID
		// expected: ErrNotOwner
		{
			"UserID Not Owner",

			1,
			2,

			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: UserID doesn't exist
		// expected: ErrNotOwner
		{
			"UserID Not Exist",

			1,
			40,

			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: Empty TaskID
		// expected: ErrInvalidTaskID
		{
			"Empty TaskID",

			0,
			1,

			wantErrIs(sqlite.ErrTaskNotFound),
		},

		// case: TaskID doesn't exist
		// expected: ErrTaskNotFound
		{
			"TaskID Not Exist",

			25,
			1,

			wantErrIs(sqlite.ErrTaskNotFound),
		},

		// case: TaskID already completed
		// expected: No-op. nil error
		{
			"Task Already Completed",

			9,
			3,

			wantNoErr,
		},

		// case: TaskID is recurring
		// expected: Task completed, new task inserted, nil error
		{
			"Recurring Task Completed",

			4,
			3,

			wantNoErr,
		},

		// case: Happy path
		// expected: nil error (task completed)
		{
			"Happy Path",

			4,
			3,

			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			beforeCompleted, err := r.GetTaskByID(tc.paramTaskID, tc.paramUserID)
			if err == nil && beforeCompleted == (sqlite.Task{}) {
				err = sqlite.ErrTaskNotFound
			}
			if err != nil {
				tc.checkErr(t, err)
			}

			gotErr := r.CompleteTaskIfOwned(tc.paramTaskID, tc.paramUserID)

			// Check returned error
			tc.checkErr(t, gotErr)

			completedTask, err := r.GetTaskByID(tc.paramTaskID, tc.paramUserID)
			if err != nil && err != gotErr {
				t.Errorf("error retrieving completed task: %v", err)
			}

			// Check DB state
			switch gotErr {
			case nil:
				if !completedTask.Done || !completedTask.CompletionDate.Valid {
					t.Errorf("database state not updated. Not completed: %v", completedTask)
				}
			default:
				checkTask(t, beforeCompleted, completedTask)
			}

			if tc.name == "Recurring Task Completed" {
				// Check next Recurring Task inserted
				insertedTask, err := r.GetTaskByID(sqlite.NextID, tc.paramUserID)
				if err != nil {
					t.Errorf("error retrieving inserted task: %v", err)
				}

				completedTask.Done = false
				completedTask.CompletionDate = sql.NullInt64{}
				completedTask.ID = sqlite.NextID

				checkTask(t, completedTask, insertedTask)
			}
		})
	}
}
