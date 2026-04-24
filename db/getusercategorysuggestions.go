package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repo) GetUserCategorySuggestions(userID int) (suggestions []string, err error) {
	var (
		caller     = "GetUserCategorySuggestions"
		rows       *sql.Rows
		suggestion string
	)
	if !r.isConnected() {
		return nil, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	rows, err = r.db.Query(`SELECT DISTINCT category FROM task WHERE user_id = ?`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]string, 0), nil
		}
		return nil, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	for rows.Next() {
		if err = rows.Scan(&suggestion); err != nil {
			return nil, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
		}

		suggestions = append(suggestions, suggestion)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	return suggestions, nil
}
