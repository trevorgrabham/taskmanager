package db_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"testing"
)

func TestAddTask(t *testing.T) {
	cases := []struct {
		name string

		paramTask sqlite.Task

		wantTask sqlite.Task
		checkErr func(*testing.T, error)
	}{

		// case: Empty UserID
		// expected: ErrInternalRepo (FK Constraint failed)
		{
			"Empty UserID",

			sqlite.Task{Title: "Empty UserID"},

			sqlite.Task{},
			wantErrIs(sqlite.ErrInternalRepo),
		},

		// case: Has negative UserID
		// expected: ErrInternalRepo (FK Constraint failed)
		{
			"Negative UserID",

			sqlite.Task{UserID: -12, Title: "Negative UserID"},

			sqlite.Task{},
			wantErrIs(sqlite.ErrInternalRepo),
		},

		// case: UserID that doens't exist
		// expected: ErrInternalRepo (FK Constraint failed)
		{
			"UserID Not Exist",

			sqlite.Task{UserID: 12345, Title: "UserID doesn't exist"},

			sqlite.Task{},
			wantErrIs(sqlite.ErrInternalRepo),
		},

		// case: Empty title
		// expected: ErrInternalRepo (NOT NULL constraint failed)
		{
			"Empty Title",

			sqlite.Task{UserID: 1, Description: sql.NullString{String: "Empty Title", Valid: true}},

			sqlite.Task{},
			wantErrIs(sqlite.ErrInternalRepo),
		},

		// case: Empty category
		// expected: Inserted task, nil error
		{
			"Empty Category",

			sqlite.Task{UserID: 1, Title: "Empty cateogry"},

			sqlite.Task{ID: sqlite.NextID, UserID: 1, Title: "Empty cateogry"},
			wantNoErr,
		},

		// case: Recurring period does not exist yet
		// expected: added task (new RecurringID), nil error.
		{
			"New Recurring Period",

			sqlite.Task{UserID: 2, Title: "New recurring period", RecurringPeriod: sql.NullString{String: "1 days", Valid: true}},

			sqlite.Task{ID: sqlite.NextID, UserID: 2, Title: "New recurring period", RecurringPeriod: sql.NullString{String: "1 days", Valid: true}, RecurringID: sql.NullInt64{Int64: 3, Valid: true}},
			wantNoErr,
		},

		// case: Recurring period already exists
		// expected: added task (same RecurringID as previous), nil error.
		{
			"Recurring Period Exists",

			sqlite.Task{UserID: 1, Title: "Existing recurring period", RecurringPeriod: sql.NullString{String: "2 days", Valid: true}},

			sqlite.Task{ID: sqlite.NextID, UserID: 1, Title: "Existing recurring period", RecurringPeriod: sql.NullString{String: "2 days", Valid: true}, RecurringID: sql.NullInt64{Int64: 1, Valid: true}},
			wantNoErr,
		},

		// case: Already Done
		// expected: Field ignored, nil error
		{
			"Task Already Done",

			sqlite.Task{UserID: 2, Title: "Already Done", Done: true},

			sqlite.Task{ID: sqlite.NextID, UserID: 2, Title: "Already Done", Done: false},
			wantNoErr,
		},

		// case: Has CompletionDate
		// expected: Field ignored, nil error
		{
			"Task Already Completed",

			sqlite.Task{UserID: 2, Title: "Already Completed", CompletionDate: sql.NullInt64{Int64: mustParseDate("2026-06-06").Unix(), Valid: true}},

			sqlite.Task{ID: sqlite.NextID, UserID: 2, Title: "Already Completed", CompletionDate: sql.NullInt64{}},
			wantNoErr,
		},

		// case: Has ID
		// expected: Field ignored, nil error
		{
			"Task Has ID",

			sqlite.Task{ID: 18, UserID: 1, Title: "With ID set"},

			sqlite.Task{ID: sqlite.NextID, UserID: 1, Title: "With ID set"},
			wantNoErr,
		},

		// case: DueDate is already passed
		// expected: Inserted as is, nil error
		{
			"DueDate Past Due",

			sqlite.Task{UserID: 2, Title: "DueDate already due", DueDate: sql.NullInt64{Int64: mustParseDate("2025-05-05").Unix(), Valid: true}},

			sqlite.Task{ID: sqlite.NextID, UserID: 2, Title: "DueDate already due", DueDate: sql.NullInt64{Int64: mustParseDate("2025-05-05").Unix(), Valid: true}},
			wantNoErr,
		},

		// case: Happy path
		// expected: added task, nil error
		{
			"Happy Path",

			sqlite.Task{UserID: 1, Title: "Happy path", Category: sql.NullString{String: "happy tests", Valid: true}, Description: sql.NullString{String: "So so happy", Valid: true}, RecurringID: sql.NullInt64{Int64: 1, Valid: true}},

			sqlite.Task{ID: sqlite.NextID, UserID: 1, Title: "Happy path", Category: sql.NullString{String: "happy tests", Valid: true}, Description: sql.NullString{String: "So so happy", Valid: true}, RecurringID: sql.NullInt64{Int64: 1, Valid: true}, RecurringPeriod: sql.NullString{String: "2 days", Valid: true}},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotTask, gotErr := r.AddTask(tc.paramTask)

			checkTask(t, tc.wantTask, gotTask)
			tc.checkErr(t, gotErr)
		})
	}
}
