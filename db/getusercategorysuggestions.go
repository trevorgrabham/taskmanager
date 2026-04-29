package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetUserCategorySuggestions returns a list of unique categories for the given userID.
//
// If the Repo was not initialized an ErrNotConnected is returned.
// If userID is empty, an empty list is returned.
// If an error occurs querying the Repo, an ErrInternalRepo is returned.
func (r *Repo) GetUserCategorySuggestions(userID int) (suggestions []string, err error) {
	var (
		caller     = "GetUserCategorySuggestions"
		rows       *sql.Rows
		suggestion string
	)
	if !r.isConnected() {
		return nil, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}
	if userID < 1 {
		return make([]string, 0), nil
	}

	rows, err = r.db.Query(`SELECT DISTINCT category FROM task WHERE user_id = ?`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]string, 0), nil
		}
		return nil, fmt.Errorf("%s: %w: %s", caller, ErrInternalRepo, err)
	}

	for rows.Next() {
		if err = rows.Scan(&suggestion); err != nil {
			return nil, fmt.Errorf("%s: %w: %s", caller, ErrInternalRepo, err)
		}

		suggestions = append(suggestions, suggestion)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", caller, ErrInternalRepo, err)
	}

	return suggestions, nil
}
