package services_test

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// GetUserBySessionID retrieves the User data associated with sessionID
//
// Validation:
//   - ErrEmptySessionID
//     sessionID not set
//   - ErrSessionNotFound
//     sessionID doesn't exist (or expired)
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
func TestGetUserBySessionID(t *testing.T) {
	cases := []struct {
		name string

		paramSessionID string

		returnUser sqlite.User
		returnErr  error

		wantedUser services.User
		checkErr   func(*testing.T, error)
	}{

		// case: Empty SessionID
		// expected: ErrEmptySessionID
		{
			"Empty SessionID",

			"",

			sqlite.User{
				ID:        1,
				Username:  "testUser",
				Password:  "testPassword",
				CreatedAt: mustParseDate("2026-06-06").Unix(),
				UpdatedAt: mustParseDate("2026-06-06").Unix(),
			},
			nil,

			services.User{},
			wantErrIs(services.ErrEmptySessionID),
		},

		// case: Empty return from repo (SessionID not exist or expired)
		// expected: ErrSessionNotFound
		{
			"SessionID Not Exist Or Expired",

			services.GenerateSessionID(),

			sqlite.User{},
			nil,

			services.User{},
			wantErrIs(services.ErrSessionNotFound),
		},

		// case: Repo returns ErrInternalRepo
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			services.GenerateSessionID(),

			sqlite.User{},
			sqlite.ErrInternalRepo,

			services.User{},
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: User, nil error
		{
			"Happy Path",

			"2904kdjsf089723kj",

			sqlite.User{
				ID:        2,
				Username:  "HappyUser",
				Password:  "HappyPassword",
				CreatedAt: mustParseDate("2026-05-05").Unix(),
				UpdatedAt: mustParseDate("2026-05-05").Unix(),
			},
			nil,

			services.User{
				ID:        2,
				Username:  "HappyUser",
				Password:  "HappyPassword",
				CreatedAt: mustParseDate("2026-05-05"),
				UpdatedAt: mustParseDate("2026-05-05"),
			},
			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockUserBySessionID{
				stubRepo:   stubRepo{callCounts: make(map[string]int)},
				returnUser: tc.returnUser,
				returnErr:  tc.returnErr,
			}

			s := services.NewService(mock)
			gotUser, gotErr := s.GetUserBySessionID(tc.paramSessionID)

			checkUser(t, gotUser, tc.wantedUser)
			tc.checkErr(t, gotErr)
		})
	}
}
