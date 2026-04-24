package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repo) GetUserBySessionID(sessionID string) (user User, err error) {
	var (
		caller = "GetUserBySessionID"
		row *sql.Row
	)
	if !r.isConnected() { return User{}, fmt.Errorf("%s: %w", caller, ErrNotConnected) }

	row = r.db.QueryRow(fmt.Sprintf(`
		%s
		WHERE id = (
			SELECT user_id FROM session WHERE id = ?
		)`, defaultUserSelect),
		sessionID)

	if user, err = scanUser(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return User{}, nil }
		return User{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	return user, nil
}
