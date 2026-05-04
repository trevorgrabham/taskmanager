package db_test

import (
	sqlite "local/taskmanager/db"
	"testing"
)

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
			gotErr := r.CompleteTaskIfOwned(tc.paramTaskID, tc.paramUserID)

			tc.checkErr(t, gotErr)
		})
	}
}
