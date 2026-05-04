package db_test

import (
	sqlite "local/taskmanager/db"
	"testing"
)

// GetUserBySessionID retrieves the user associated with sessionID
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Empty:
//   - sessionID doesn't exist
//   - sessionID expired
//
// Happy Path:
//   - returns the user with all fields populated
func TestGetUserBySessionID(t *testing.T) {
	cases := []struct {
		name string

		sessionID string

		wantUser sqlite.User
		checkErr func(*testing.T, error)
	}{

		// case: SessionID doesn't exist
		// expected: no User, nil error
		{
			"SessionID Not Exist",

			"",

			sqlite.User{},
			wantNoErr,
		},

		// case: SessionID matches, but it is expired
		// expected: no User, nil error
		{
			"SessionID Expired",

			"36ba76c7856f0aea25fb415bfba50721195efbb19e13ccabd7f384f2dfaa62",

			sqlite.User{},
			wantNoErr,
		},

		// case: Happy path
		// expected: User, nil error
		{
			"Happy Path",

			"66ba76c7856f0aea25fb415bfba50721195efbb19e13ccabd7f384f2dfaa62",

			sqlite.User{ID: 1, Username: "tester1", Password: "password1", CreatedAt: mustParseDate("1997-10-19").Unix(), UpdatedAt: mustParseDate("1997-10-19").Unix()},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotUser, gotErr := r.GetUserBySessionID(tc.sessionID)

			checkUser(t, gotUser, tc.wantUser)

			tc.checkErr(t, gotErr)
		})
	}
}
