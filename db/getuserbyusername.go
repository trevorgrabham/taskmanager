package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetUserByUsername returns the User data associated with username.
//
// If username is empty or not found, treated as a no-op.
// If an error occurs while querying the database, an ErrInternalRepo is returned.
func (r *Repo) GetUserByUsername(username string) (user User, err error) {
	row := r.db.QueryRow(fmt.Sprintf(`
		%s 
		WHERE username = ?`, defaultUserSelect),
		username)

	if user, err = scanUser(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}
		return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return user, nil
}
