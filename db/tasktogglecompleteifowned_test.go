package db_test

import (
	sqlite "local/taskmanager/db"
	"testing"
)

// TaskToggleCompleteIfOwned toggles the completion status of task if userID is the owner
//
// Errors:
//   - ErrTaskNotFound
//     taskID doesn't exist
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

		checkErr func(*testing.T, error)
	}{

		// case: Empty TaskID
		// expected: ErrTaskNotFound
		{
			"Empty TaskID",

			0,
			1,

			wantErrIs(sqlite.ErrTaskNotFound),
		},

		// case: Empty UserID
		// expected: ErrNotOwner
		{
			"Empty UserID",

			3,
			0,

			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: UserID not owner
		// expected: ErrNotOwner
		{
			"UserID Not Owner",

			8,
			3,

			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: Task doesn't exist
		// expected: ErrTaskNotFound
		{
			"TaskID Not Exist",

			35,
			1,

			wantErrIs(sqlite.ErrTaskNotFound),
		},

		// case: Task was completed
		// expected: Task now incomplete. nil error
		{
			"Task Already Completed",

			9,
			3,

			wantNoErr,
		},

		// case: Task was incomplete
		// expected: Task now completed. nil error
		{
			"Task Incomplete",

			6,
			3,

			wantNoErr,
		},

		// case: Completed a recurring task
		// expected: Task now completed, and a new task inserted, nil error
		{
			"Completed Recurring Task",

			4,
			3,

			wantNoErr,
		},

		// case: Uncompleted a recurring task
		// expected: Task now incomplete, recurring task that was previously inserted remains, nil error
		{
			"Uncompleted Recurring Task",

			9,
			3,

			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			beforeToggle, err := r.GetTaskByID(tc.paramTaskID, tc.paramUserID)
			if err != nil {
				tc.checkErr(t, err)
			}

			gotErr := r.TaskToggleCompleteIfOwned(tc.paramTaskID, tc.paramUserID)

			// Check returned error
			tc.checkErr(t, gotErr)

			afterToggle, err := r.GetTaskByID(tc.paramTaskID, tc.paramUserID)
			if err != nil {
				tc.checkErr(t, err)
				// t.Errorf("error retrieving task after toggle: %v", err)
			}

			// Check DB state
			switch gotErr {
			case nil:
				if afterToggle.Done == beforeToggle.Done || afterToggle.CompletionDate.Valid == beforeToggle.CompletionDate.Valid {
					t.Errorf("task not toggled.\nbefore toggle: %v\nafter togge: %v", beforeToggle, afterToggle)
				}
			default:
				if afterToggle.Done != beforeToggle.Done || afterToggle.CompletionDate.Valid != beforeToggle.CompletionDate.Valid {
					t.Errorf("task toggled.\nbefore toggle: %v\nafter togge: %v", beforeToggle, afterToggle)
				}
			}

			// If we completed a Recurring Task, check DB state for the new task
			if tc.name == "Completed Recurring Task" {
				insertedTask, err := r.GetTaskByID(sqlite.NextID, tc.paramUserID)
				if err != nil {
					t.Errorf("retrieving inserted task: %v", err)
					return
				}

				beforeToggle.ID = sqlite.NextID
				checkTask(t, beforeToggle, insertedTask)
			}
		})
	}
}
