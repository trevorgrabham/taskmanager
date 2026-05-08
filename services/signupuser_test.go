package services_test

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
	"time"
)

// SignupUser validates the user's username and password. If valid, an account is created for the user, and a new session is created
//
// Validation:
//   - ErrUserValidation {Field: UserFieldUsername, Message: MessageUsernameTaken}
//     username already registered to another user
//   - ErrUserValidation {Field: UserFieldPassword, Message: MessageInvalidPassword}
//     password empty
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - returns the newly created sessionID for user associated with the username
func TestSignupUser(t *testing.T) {
	cases := []struct {
		name string

		paramUsername string
		paramPassword string

		returnUser       sqlite.User
		returnUserErr    error
		returnSessionErr error

		checkErr func(*testing.T, error)
	}{

		// case: Empty username
		// expected: ErrUserValidation {Field: UserFieldUsername, Message: MessageUsernameTaken}
		{
			"Empty Username",

			"",
			"ValidPassword",

			sqlite.User{
				ID: 3,
				Username: "ValidUsername",
				Password: "ValidPassword",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			nil,
			nil,

			wantErrUserValidation(services.UserFieldUsername, services.MessageInvalidUsername),
		},

		// case: Empty password
		// expected: ErrUserValidation {Field: UserFieldPassword, Message: MessageInvalidPassword}
		{
			"Empty Password",

			"ValidUsername",
			"",

			sqlite.User{
				ID: 2,
				Username: "ValidUsername",
				Password: "ValidPassword",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			nil,
			nil,

			wantErrUserValidation(services.UserFieldPassword, services.MessageInvalidPassword),
		},

		// case: Username already taken (returned ErrUsernameTaken)
		// expected: ErrUserValidation {Field: UserFieldUsername, Message: MessageInvalidUsername}
		{
			"Username Taken",

			"NonUniqueUsername",
			"TotallyUniquePassword",

			sqlite.User{},
			sqlite.ErrUsernameTaken,
			nil,

			wantErrUserValidation(services.UserFieldUsername, services.MessageUsernameTaken),
		},

		// case: Repo returned ErrInternalRepo from SignupUser()
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo SignupUser",

			"UniqueUsername",
			"Password",

			sqlite.User{},
			sqlite.ErrInternalRepo,
			nil,

			wantErrIs(services.ErrInternalRepo),
		},

		// case: Repo returned ErrInternalRepo from StartSession()
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo StartSession",

			"UniqueUsername",
			"Password",

			sqlite.User{
				ID: 3,
				Username: "UniqueUsername",
				Password: "Password",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			nil,
			sqlite.ErrInternalRepo,

			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: User exists, session started, nil error
		{
			"Happy Path",

			"MyUsername",
			"test",

			sqlite.User{
				ID: 3,
				Username: "MyUsername",
				Password: "test",
				CreatedAt: mustParseDate("2020-02-02").Unix(),
				UpdatedAt: mustParseDate("2020-02-02").Unix(),
			},
			nil,
			nil,

			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockSignupUser{
				stubRepo:         stubRepo{callCounts: make(map[string]int)},
				returnUser:       tc.returnUser,
				returnUserErr:    tc.returnUserErr,
				returnSessionErr: tc.returnSessionErr,
			}
			s := services.NewService(mock)

			gotSessionID, gotErr := s.SignupUser(tc.paramUsername, tc.paramPassword)

			tc.checkErr(t, gotErr)

			if gotErr == nil && gotSessionID == "" {
				t.Errorf("wanted sessionID, got empty string")
			}
		})
	}
}
