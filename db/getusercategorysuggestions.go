package db

import (
	"database/sql"
	"fmt"
)

// GetUserCategorySuggestions returns a list of unique categories for the given userID.
//
// If userID is empty or doesn't exist, an empty list is returned.
// If an error occurs querying the Repo, an ErrInternalRepo is returned.
func (r *Repo) GetUserCategorySuggestions(userID int) (suggestions []string, err error) {
	var (
		rows       *sql.Rows
		suggestion sql.NullString
	)
	rows, err = r.db.Query(`SELECT DISTINCT category FROM task WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	for rows.Next() {
		if err = rows.Scan(&suggestion); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}
		if !suggestion.Valid {
			continue
		}

		suggestions = append(suggestions, suggestion.String)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return suggestions, nil
}
