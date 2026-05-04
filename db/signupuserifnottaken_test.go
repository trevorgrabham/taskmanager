package db_test

import (
	sqlite "local/taskmanager/db"
	"testing"
	"time"
)

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

			tc.checkErr(t, gotErr)

			checkUser(t, gotUser, tc.wantUser)

		})
	}
}
