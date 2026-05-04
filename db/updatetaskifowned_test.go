package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
)

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
// expected: ErrInternalRepo (CHECK Constraint Failure)
		{
			"Empty Title",

			sqlite.Task{ID: 1, UserID: 1, Category: sql.NullString{String: "This shouldn't work", Valid: true}},

			sqlite.Task{},
			wantErrIs(sqlite.ErrInternalRepo),
		},

// case: Done is true
// expected: Done ignored, nil error
		{
			"Already Done",

			sqlite.Task{ID: 6, UserID: 3, Title: "New Title",Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "New description", Valid: true}, Done: true},

			sqlite.Task{ID: 6, UserID: 3, Title: "New Title", Category: sql.NullString{String: "go tests", Valid: true}, Description: sql.NullString{String: "New description", Valid: true} },
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
			"New RecurringPeriod",

			sqlite.Task{ID: 2, UserID: 2, Title: "Second Task", Category: sql.NullString{String: "go tests", Valid: true}, RecurringPeriod: sql.NullString{String: "5 days", Valid: true}},

			sqlite.Task{ID: 2, UserID: 2, Title: "Second Task", Category: sql.NullString{String: "go tests", Valid: true}, RecurringPeriod: sql.NullString{String: "5 days", Valid: true}, RecurringID: sql.NullInt64{Int64: 3, Valid: true}},
			wantNoErr,
		},

// case: Recurring period does exist already
// expected: Old recurring_id value, nil error
		{
			"Old RecurringPeriod",

			sqlite.Task{ID: 2, UserID: 2, Title: "Second Task", Category: sql.NullString{String: "go tests", Valid: true}, RecurringPeriod: sql.NullString{String: "2 days", Valid: true}},
			
			sqlite.Task{ID: 2, UserID: 2, Title: "Second Task", Category: sql.NullString{String: "go tests", Valid: true}, RecurringPeriod: sql.NullString{String: "2 days", Valid: true}, RecurringID: sql.NullInt64{Int64: 1, Valid: true}},
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

			sqlite.Task{ID: 1, UserID: 1, Title: "New Title", Category: sql.NullString{String: "New Category", Valid: true}, Description: sql.NullString{String: "New Description", Valid: true}, RecurringPeriod: sql.NullString{String: "1 days", Valid: true}},

			sqlite.Task{ID: 1, UserID: 1, Title: "New Title", Category: sql.NullString{String: "New Category", Valid: true}, Description: sql.NullString{String: "New Description", Valid: true}, RecurringPeriod: sql.NullString{String: "1 days", Valid: true}, RecurringID: sql.NullInt64{Int64: 3, Valid: true}},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotTask, gotErr := r.UpdateTaskIfOwned(tc.paramTask)

			checkTask(t, tc.wantTask, gotTask)
			tc.checkErr(t, gotErr)
		})
	}
}
