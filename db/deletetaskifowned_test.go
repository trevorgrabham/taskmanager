package db_test

import (
	sqlite "local/taskmanager/db"
	"testing"
)

// DeleteTask deletes the task identified by taskID if it's owned by userID
//
// Errors:
//   - ErrNotOwner
//     userID doesn't exist
//     userID isn't the owner
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - task is deleted. If it didn't exist, no-op.
func TestDeleteTask(t *testing.T) {
	cases := []struct {
		name string

		paramTaskID int
		paramUserID int

		checkErr func(*testing.T, error)
	}{

		// case: TaskID doesn't exist
		// expected: no-op. nil error
		{
			"TaskID Not Exist",

			0,
			1,

			wantNoErr,
		},

		// case: UserID doesn't exist
		// expected: ErrNotOwner
		{
			"UserID Not Exist",

			1,
			5,

			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: UserID not the owner
		// expected: ErrNotOwner
		{
			"UserID Not Owner",

			1,
			3,

			wantErrIs(sqlite.ErrNotOwner),
		},

		// case: Happy path
		// expected: nil error (task no longer there)
		{
			"Happy Path",

			3,
			1,

			wantNoErr,
		}}
	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotErr := r.DeleteTaskIfOwned(tc.paramTaskID, tc.paramUserID)

			tc.checkErr(t, gotErr)
		})
	}
}
