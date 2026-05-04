package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
)

// UpdateTaskIfOwned updates taskID if it is owned by userID
//
// Errors:
//   - ErrNotOwner
//     userID doesn't exist
//     userID not the owner
//   - ErrInternalRepo
//     transient database error
//
// Empty:
//   - taskID doesn't exist
//
// Happy Path:
//   - returns the updated task with all its fields populated
func TestUpdateTask(t *testing.T) {
	cases := []struct {
		name string

		paramTask sqlite.Task

		wantTask sqlite.Task
		checkErr func(*testing.T, error)
	}{

		// case: TaskID doesn't exist
		// expected: No-op. Empty task, nil error
		{
			"TaskID Not Exist",

			sqlite.Task{ID: 0, UserID: 1, Title: "Shouldn't see this"},

			sqlite.Task{},
			wantNoErr,
		},

		// case: UserID doesn't exist
		// expected: ErrNotOwner
		{
			"UserID Not Exist",

			sqlite.Task{ID: 1, UserID: 97, Title: "Still shouldn't see this"},

			sqlite.Task{},
			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: UserID not owner
		// expected: ErrNotOwner
		{
			"UserID Not Owner",

			sqlite.Task{ID: 1, UserID: 3, Title: "If you're seeing this something went wrong"},

			sqlite.Task{},
			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: Empty title
		// expected: ErrConstraintFailure (CHECK Constraint Failure)
		{
			"Empty Title",

			sqlite.Task{ID: 1, UserID: 1, Category: sql.NullString{String: "This shouldn't work", Valid: true}},

			sqlite.Task{},
			wantErrIs(sqlite.ErrConstraintFailure),
		},

		// case: Done is true
		// expected: Done ignored, nil error
		{
			"Already Done",

			sqlite.Task{ID: 6, UserID: 3, Title: "New Title", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "New description", Valid: true}, Done: true},

			sqlite.Task{ID: 6, UserID: 3, Title: "New Title", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "New description", Valid: true}},
			wantNoErr,
		},

		// case: Completion date is not empty
		// expected: Nil error
		{
			"Already Completed",

			sqlite.Task{ID: 6, UserID: 3, Title: "Task Title", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "New description", Valid: true}, CompletionDate: sql.NullInt64{Int64: 1, Valid: true}},

			sqlite.Task{ID: 6, UserID: 3, Title: "Task Title", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "New description", Valid: true}},
			wantNoErr,
		},

		// case: Recurring period does not exist yet
		// expected: New recurring_id value, nil error
		{
			"Has RecurringPeriod",

			sqlite.Task{ID: 2, UserID: 2, Title: "Second Task", Category: sql.NullString{String: "go tests", Valid: true}, RecurringPeriod: sql.NullInt64{Int64: 3600 * 24 * 5, Valid: true}},

			sqlite.Task{ID: 2, UserID: 2, Title: "Second Task", Category: sql.NullString{String: "go tests", Valid: true}, RecurringPeriod: sql.NullInt64{Int64: 3600 * 24 * 5, Valid: true}},
			wantNoErr,
		},

		// case: Recurring period is removed
		// expected: No recurring_id value, nil error
		{
			"RecurringPeriod Removed",

			sqlite.Task{ID: 4, UserID: 3, Title: "With RecurringPeriod", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "Seeing if the recurring_id works", Valid: true}, DueDate: sql.NullInt64{Int64: mustParseDate("2010-10-10").Unix(), Valid: true}},

			sqlite.Task{ID: 4, UserID: 3, Title: "With RecurringPeriod", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "Seeing if the recurring_id works", Valid: true}, DueDate: sql.NullInt64{Int64: mustParseDate("2010-10-10").Unix(), Valid: true}},
			wantNoErr,
		},

		// case: Happy path
		// expected: Updated task, nil error

		{
			"Happy Path",

			sqlite.Task{ID: 1, UserID: 1, Title: "New Title", Category: sql.NullString{String: "New Category", Valid: true}, Description: sql.NullString{String: "New Description", Valid: true}, RecurringPeriod: sql.NullInt64{Int64: 3600 * 24, Valid: true}},

			sqlite.Task{ID: 1, UserID: 1, Title: "New Title", Category: sql.NullString{String: "New Category", Valid: true}, Description: sql.NullString{String: "New Description", Valid: true}, RecurringPeriod: sql.NullInt64{Int64: 3600 * 24, Valid: true}},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			beforeUpdate, err := r.GetTaskByID(tc.paramTask.ID, tc.paramTask.UserID)
			if err != nil {
				tc.checkErr(t, err)
			}

			gotTask, gotErr := r.UpdateTaskIfOwned(tc.paramTask)

			// Check returned error
			tc.checkErr(t, gotErr)

			// Check returned task
			checkTask(t, tc.wantTask, gotTask)

			updatedTask, err := r.GetTaskByID(tc.paramTask.ID, tc.paramTask.UserID)
			if err != nil {
				tc.checkErr(t, err)
			}

			// Check DB state
			switch gotErr {
			case nil:
				checkTask(t, tc.wantTask, updatedTask)
			default:
				checkTask(t, beforeUpdate, updatedTask)
			}
		})
	}
}
