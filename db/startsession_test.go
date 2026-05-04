package db_test

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

func TestStartSession(t *testing.T) {
	cases := []struct {
		name string

		paramSessionID string 
		paramUserID int

		checkErr func(*testing.T, error)
}{

// case: UserID doesn't exist
// expected: ErrInternalRepo (FK Constraint fail)
		{
			"UserID Not Exist",

			services.GenerateSessionID(),
			-12,

			wantErrIs(sqlite.ErrInternalRepo),
		},

// case: Empty SessionID
// expected: ErrInternalRepo (CHECK Constraint fail)
		{
			"Empty SessionID",

			"",
			4,

			wantErrIs(sqlite.ErrInternalRepo),
		},

// case: SessionID already exists
// expected: ErrSessionAlreadyExists
		{
			"SessionID Already Exists",

			"56ba76c7856f0aea25fb415bfba50721195efbb19e13ccabd7f384f2dfaa62",
			4,

			wantErrIs(sqlite.ErrInternalRepo),
		},

// case: UserID already has a different session
// expected: old session deleted, new session created, nil error
		{
			"User Has Session",

			services.GenerateSessionID(),
			1,

			wantNoErr,
		},

// case: Happy path
// expected: new session created, nil error
		{
			"Happy Path",

			services.GenerateSessionID(),
			4,

			wantNoErr,
	}}
	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotErr := r.StartSession(tc.paramSessionID, tc.paramUserID)

			tc.checkErr(t, gotErr)

			if gotErr != nil { return }

			gotUser, gotErr := r.GetUserBySessionID(tc.paramSessionID)
			if gotErr != nil {
				t.Errorf("got %v retreiving new session", gotErr)
			}

			if gotUser.ID != tc.paramUserID { t.Errorf("wanted session UserID: %d, got: %d", tc.paramUserID, gotUser.ID) }
		})
	}
}
