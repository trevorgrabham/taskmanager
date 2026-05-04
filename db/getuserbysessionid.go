package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// GetUserBySessionID returns the User data associated with SessionID.
//
// If SessionID is empty, doesn't exist, or expired, acts as a no-op.
// If an error occurs while querying the database, an ErrInternalRepo is returned.
func (r *Repo) GetUserBySessionID(sessionID string) (user User, err error) {
	row := r.db.QueryRow(fmt.Sprintf(`
		%s
		WHERE id = (
			SELECT user_id FROM session WHERE id = ? AND expires_at > ?
		)`, defaultUserSelect),
		sessionID, time.Now().Unix())

	if user, err = scanUser(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}
		return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return user, nil
}
