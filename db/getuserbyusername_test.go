package db_test

import (
	sqlite "local/taskmanager/db"
	"testing"
)

// GetUserByUsername retrieves the user with username
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Empty:
//   - username doesn't exist
//
// Happy Path:
//   - returns the user with all fields populated
func TestGetUserByUsername(t *testing.T) {
	cases := []struct {
		name string

		paramUsername string

		wantUser sqlite.User
		checkErr func(*testing.T, error)
	}{

		// case: Empty username
		// expected: No user, nil error
		{
			"Empty Username",

			"",

			sqlite.User{},
			wantNoErr,
		},

		// case: Username doesn't exist
		// expected: No user, nil error
		{
			"Username Not Exist",

			"tester12",

			sqlite.User{},
			wantNoErr,
		},

		// case: Happy path
		// expected: User, nil error
		{
			"Happy Path",

			"tester1",

			sqlite.User{ID: 1, Username: "tester1", Password: "password1", CreatedAt: mustParseDate("1997-10-19").Unix(), UpdatedAt: mustParseDate("1997-10-19").Unix()},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotUser, gotErr := r.GetUserByUsername(tc.paramUsername)

			checkUser(t, gotUser, tc.wantUser)

			tc.checkErr(t, gotErr)
		})
	}
}
