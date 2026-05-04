package db_test

import (
	sqlite "local/taskmanager/db"
	"slices"
	"testing"
)

// GetUserCategorySuggestions retrieves the list of unique categories for userID
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Empty:
//   - userID doesn't exist
//   - no categories
//
// Happy Path:
//   - returns the list of unique categories for userID
func TestGetUserCategorySuggestions(t *testing.T) {
	cases := []struct {
		name string

		paramUserID int

		wantCategories []string
		checkErr       func(*testing.T, error)
	}{

		// case: Empty UserID
		// expected: No categories, nil error
		{
			"Empty UserID",

			0,

			nil,
			wantNoErr,
		},

		// case: UserID doesn't exist
		// expected: No categories, nil error
		{
			"UserID Not Exist",

			199,

			nil,
			wantNoErr,
		},

		// case: No category suggestions
		// expected: Empty list, nil error
		{
			"No Categories",

			4,

			nil,
			wantNoErr,
		},

		// case: Happy path
		// expected: List of categories, nil error
		{
			"Happy Path",

			1,

			[]string{"go tests", "general tests"},
			wantNoErr,
		}}

	for _, tc := range cases {
		sqlite.ResetDB(t, r)
		t.Run(tc.name, func(t *testing.T) {
			gotCategories, gotErr := r.GetUserCategorySuggestions(tc.paramUserID)

			tc.checkErr(t, gotErr)

			if len(gotCategories) != len(tc.wantCategories) {
				t.Errorf("wanted categories: %v\ngot: %v", tc.wantCategories, gotCategories)
				return
			}

			slices.Sort(gotCategories)
			slices.Sort(tc.wantCategories)

			for i := range gotCategories {
				if gotCategories[i] != tc.wantCategories[i] {
					t.Errorf("wanted categories: %v\ngot: %v", tc.wantCategories, gotCategories)
					break
				}
			}

		})
	}
}
