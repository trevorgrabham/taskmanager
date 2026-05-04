package db_test

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

func TestDeleteSession(t *testing.T) {
	cases := []struct {
		name string

		paramSessionID string

		checkErr func(*testing.T, error)
	}{

		// case: sessionID is empty
		// expected: nil error (no-op)
		{
			"Empty SessionID",

			"",

			wantNoErr,
		},

		// case: sessionID doesn't exist
		// expected: nil error (no-op)
		{
			"SessionID Not Exist",

			services.GenerateSessionID(),

			wantNoErr,
		},

		// case: Happy path
		// expected: nil
		{
			"Happy Path",

			"66ba76c7856f0aea25fb415bfba50721195efbb19e13ccabd7f384f2dfaa62",

			wantNoErr,
		}}
	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotErr := r.DeleteSession(tc.paramSessionID)

			tc.checkErr(t, gotErr)
		})
	}
}
