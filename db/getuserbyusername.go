package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repo) GetUserByUsername(username string) (user User, err error) {
	var (
		caller = "GetUserByUsername"
		row    *sql.Row
	)
	if !r.isConnected() {
		return User{}, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	row = r.db.QueryRow(fmt.Sprint(`
		%s 
		WHERE username = ?`, defaultUserSelect),
		username)

	if user, err = scanUser(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}
		return User{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	return user, nil
}
