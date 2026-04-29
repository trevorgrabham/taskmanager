package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetUserByUsername returns the User data associated with username.
// 
// If the Repo is not initialized an ErrNotConnected is returned.
// If username is empty, an empty User is returned (User.ID = 0).
// If no User is found, an ErrUserNotExist is returned.
// If an error occurs while querying the database, an ErrInternalRepo is returned.
func (r *Repo) GetUserByUsername(username string) (user User, err error) {
	var row    *sql.Row
	if !r.isConnected() {
		return User{}, ErrNotConnected
	}
	if username == "" {
		return User{}, nil
	}

	row = r.db.QueryRow(fmt.Sprintf(`
		%s 
		WHERE username = ?`, defaultUserSelect),
		username)

	if user, err = scanUser(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("%w: for %s", ErrUserNotExist, username)
		}
		return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return user, nil
}
