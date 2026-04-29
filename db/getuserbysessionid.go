package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetUserBySessionID returns the User data associated with SessionID.
//
// If the Repo is not initialized an ErrNotConnected is returned.
// If SessionID is empty, an empty User is returned (User.ID = 0).
// If no session is found, an ErrNoLiveSession is returned.
// If an error occurs while querying the database, an ErrInternalRepo is returned.
func (r *Repo) GetUserBySessionID(sessionID string) (user User, err error) {
	var row *sql.Row
	if !r.isConnected() {
		return User{}, ErrNotConnected
	}
	if sessionID == "" {
		return User{}, nil
	}

	row = r.db.QueryRow(fmt.Sprintf(`
		%s
		WHERE id = (
			SELECT user_id FROM session WHERE id = ?
		)`, defaultUserSelect),
		sessionID)

	if user, err = scanUser(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("%w for %s", ErrNoLiveSession, sessionID)
		}
		return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return user, nil
}
