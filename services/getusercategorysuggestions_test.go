package services_test

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// GetUserCategorySuggestions retrieves a list of unique categories for userID
//
// Validation:
//   - ErrUserValidation {Field: UserFieldID}
//     userID negative or zero valued (< 1)
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - returns a list of unique categories for userID
func TestGetUserCategorySuggestions(t *testing.T) {
	cases := []struct {
		name string

		paramUserID int

		returnSuggestions []string
		returnErr         error

		wantSuggestions []string
		checkErr        func(*testing.T, error)
	}{

		// case: Empty UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Empty UserID",

			0,

			[]string{"go tests", "stop tests", "continue tests"},
			nil,

			nil,
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Negative UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Negative UserID",

			-3,

			[]string{"go tests", "stop tests", "continue tests"},
			nil,

			nil,
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: UserID doesn't exist
		// expected: Empty list, nil error
		{
			"UserID Not Exist",

			1234,

			nil,
			nil,

			nil,
			wantNoErr,
		},

		// case: No category suggestions
		// expected: Empty list, nil error
		{
			"No suggestions",

			1,

			nil,
			nil,

			nil,
			wantNoErr,
		},

		// case: Repo returned ErrInternalRepo
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			1,

			nil,
			sqlite.ErrInternalRepo,

			nil,
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: List of catgories, nil error
		{
			"Happy Path",

			1,

			[]string{"testing", "failing", "passing", "testing again"},
			nil,

			[]string{"testing", "failing", "passing", "testing again"},
			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockUserCategorySuggestions{
				stubRepo:          stubRepo{callCounts: make(map[string]int)},
				returnSuggestions: tc.returnSuggestions,
				returnErr:         tc.returnErr,
			}
			s := services.NewService(mock)

			gotSuggestions, gotErr := s.GetUserCategorySuggestions(tc.paramUserID)

			tc.checkErr(t, gotErr)

			if len(gotSuggestions) != len(tc.wantSuggestions) {
				t.Errorf("UserSuggestions wanted: %v\ngot: %v", tc.wantSuggestions, gotSuggestions)
			}
			for i := range gotSuggestions {
				if gotSuggestions[i] != tc.wantSuggestions[i] {
					t.Errorf("UserSuggestions wanted: %v\ngot: %v", tc.wantSuggestions, gotSuggestions)
				}
			}
		})
	}
}
