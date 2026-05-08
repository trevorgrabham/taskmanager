package services_test

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
	"time"
)

// LoginUser validates the user's login credentials. If valid, a session is created for the user
//
// Validation:
//   - ErrUserValidation {Field: UserFieldUsername, Message: MessageInvalidUsernameOrPassword}
//		username not registered to a user
//   - ErrUserValidation {Field: UserFieldPassword, Message: MessageInvalidUsernameOrPassword}
//		password does not match (or empty)
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
// 		- returns the newly created sessionID for user associated with the username
func TestLoginUser(t *testing.T) {
	cases := []struct{
		name string 

		paramUsername string
		paramPassword string

		returnUser sqlite.User
		returnUserErr error 
		returnSessionErr error 

		checkErr func(*testing.T, error)
	}{

// case: Empty username 
// expected: ErrUserValidation {Field: UserFieldUsername, Message: MessageInvalidUsernameOrPassword}
		{
			"Empty Username",

			"",
			"Password",

			sqlite.User{
				ID: 3,
				Username: "TestUser",
				Password: "TestPassword",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			nil,
			nil,

			wantErrUserValidation(services.UserFieldUsername, services.MessageInvalidUsernameOrPassword),
		},

// case: Username doesn't exist 
// expected: ErrUserValidation {Field: UserFieldUsername, Message: MessageInvalidUsernameOrPassword}
		{
			"Username Not Exist",

			"ThisDontExist",
			"ThisDoesntMatter",

			sqlite.User{},
			nil,
			nil,

			wantErrUserValidation(services.UserFieldUsername, services.MessageInvalidUsernameOrPassword),
		},

// case: Empty password 
// expected: ErrUserValidation {Field: UserFieldPassword, Message: MessageInvalidUsernameOrPassword}
		{
			"Empty Password",

			"ThisIsMyAccount",
			"",

			sqlite.User{
				ID: 1,
				Username: "ThisIsMyAccount",
				Password: "ThisIsMyPassword",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			nil,
			nil,

			wantErrUserValidation(services.UserFieldPassword, services.MessageInvalidUsernameOrPassword),
		},

// case: Password doesn't match
// expected: ErrUserValidation {Field: UserFieldPassword, Message: MessageInvalidUsernameOrPassword}
		{
			"Wrong Password",

			"ThisIsMyAccount",
			"ThisIsNotMyPassword",

			sqlite.User{
				ID: 3,
				Username: "ThisIsMyAccount",
				Password: "ThisIsMyPassword",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			nil,
			nil,

			wantErrUserValidation(services.UserFieldPassword, services.MessageInvalidUsernameOrPassword),
		},

// case: Repo returned empty User
// expected: ErrUserValidation {Field: UserFieldUsername, Message: MessageInvalidUsernameOrPassword}
		{
			"No User Returned",

			"WrongUsername",
			"GoodPassword",

			sqlite.User{},
			nil,
			nil,

			wantErrUserValidation(services.UserFieldUsername, services.MessageInvalidUsernameOrPassword),
		},

// case: Repo returned ErrInternalRepo from GetUserByUsername
// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo GetUserByUsername",

			"ThisIsMyUsername",
			"ThisIsMyPassword",

			sqlite.User{},
			sqlite.ErrInternalRepo,
			nil,

			wantErrIs(services.ErrInternalRepo),
		},

// case: Repo returned ErrInternalRepo from StartSession
// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo StartSession",

			"ThisIsMyUsername",
			"test",

			sqlite.User{
				ID: 2,
				Username: "ThisIsMyUsername",
				Password: "$2a$10$L7hSmZ9zeDgAAs76FVHPpOy0AWkPSMHF7j85iLBiSTyafSWye8Xwq$10$L7hSmZ9zeDgAAs76FVHPpOy0AWkPSMHF7j85iLBiSTyafSWye8Xwq",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			nil,
			sqlite.ErrInternalRepo,

			wantErrIs(services.ErrInternalRepo),
		},

// case: Happy path
// expected: Session started, nil error
		{
			"Happy Path",

			"MyAccount",
			"test",

			sqlite.User{
				ID: 2,
				Username: "MyAccount",
				Password: "$2a$10$L7hSmZ9zeDgAAs76FVHPpOy0AWkPSMHF7j85iLBiSTyafSWye8Xwq$10$L7hSmZ9zeDgAAs76FVHPpOy0AWkPSMHF7j85iLBiSTyafSWye8Xwq",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			nil, 
			nil,

			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockLoginUser{
				stubRepo: stubRepo{callCounts: make(map[string]int)},
				returnUser: tc.returnUser,
				returnUserErr: tc.returnUserErr,
				returnSessionErr: tc.returnSessionErr,
			}
			s := services.NewService(mock)

			gotSessionID, gotErr := s.LoginUser(tc.paramUsername, tc.paramPassword)

			tc.checkErr(t, gotErr)

			if gotErr == nil && gotSessionID == "" { t.Errorf("wanted sessionID, got empty string") }
		})
	}
}
