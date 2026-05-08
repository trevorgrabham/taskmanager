package services_test

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// DeleteSession ensures sessionID does not exist
//
// Validation:
//   - ErrEmptySessionID
//     sessionID not set
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - sessionID is removed (if it existed in the first place) from the repo
func TestDeleteSession(t *testing.T) {
	cases := []struct {
		name string

		paramSessionID string

		returnErr error

		checkErr func(*testing.T, error)
	}{

		// case: Empty sessionID
		// expected: ErrEmptySessionID
		{
			"Empty SessionID",

			"",

			nil,

			wantErrIs(services.ErrEmptySessionID),
		},

		// case: SessionID doesn't exist
		// expected: nil error
		{
			"SessionID Not Exist",

			services.GenerateSessionID(),

			nil,

			wantNoErr,
		},

		// case: Repo returned ErrInternalRepo
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			services.GenerateSessionID(),

			sqlite.ErrInternalRepo,

			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: nil
		{
			"Happy Path",

			services.GenerateSessionID(),

			nil,

			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockDeleteSession{
				stubRepo:  stubRepo{callCounts: make(map[string]int)},
				returnErr: tc.returnErr,
			}

			s := services.NewService(mock)
			gotErr := s.DeleteSession(tc.paramSessionID)

			tc.checkErr(t, gotErr)
		})
	}
}
