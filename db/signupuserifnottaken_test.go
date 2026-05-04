package db_test

import (
	sqlite "local/taskmanager/db"
	"testing"
	"time"
)

// SignupUserIfNotTaken creates an account using username and password if username is not already taken
//
// Errors:
//   - ErrUsernameTaken
//     username already taken
//   - ErrInternalRepo
//     transient database error
//
// Empty:
//   - username empty
//   - password empty
//
// Happy Path:
//   - returns a user with all fields populated for the new user account
func TestSignupUser(t *testing.T) {
	cases := []struct {
		name string

		paramUsername string
		paramPassword string

		wantUser sqlite.User
		checkErr func(*testing.T, error)
	}{

		// case: Empty username
		// expected: Empty User, nil error
		{
			"Empty Username",

			"",
			"Password",

			sqlite.User{},
			wantNoErr,
		},

		// case: Username already taken
		// expected: ErrUsernameTaken
		{
			"Username Taken",

			"tester1",
			"password",

			sqlite.User{},
			wantErrIs(sqlite.ErrUsernameTaken),
		},

		// case: Empty Password
		// expected: Empty User, nil error
		{
			"Empty Password",

			"tester5",
			"",

			sqlite.User{},
			wantNoErr,
		},

		// case: Happy path
		// expected: User, nil error
		{
			"Happy Path",

			"tester5",
			"password",

			sqlite.User{ID: 5, Username: "tester5", Password: "password", CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotUser, gotErr := r.SignupUserIfNotTaken(tc.paramUsername, tc.paramPassword)

			// Check returned error
			tc.checkErr(t, gotErr)

			// Check returned User
			checkUser(t, gotUser, tc.wantUser)

			newUser, err := r.GetUserByUsername(tc.paramUsername)
			if err != nil {
				t.Errorf("error retrieving inserted User: %v", err)
			}

			if tc.name == "Username Taken" {
				if newUser == (sqlite.User{}) {
					t.Errorf("wanted non-empty User, got empty User")
				}
			} else {
				checkUser(t, newUser, tc.wantUser)
			}
		})
	}
}
