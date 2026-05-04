package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
)

// SignupUserIfNotTaken inserts a record for the given username and password in the Repo. It is the callers responsibility to ensure that the passed password has been hashed.
//
// If username or password is empty, treated as a no-op.
// If username already exists in the Repo, an ErrUsernameExists is returned.
// If an error occurs while querying the database, an ErrInternalRepo is returned.
func (r *Repo) SignupUserIfNotTaken(username, hashedPassword string) (user User, err error) {
	var (
		tx  *sql.Tx
		row *sql.Row
		now time.Time
	)
	if tx, err = r.db.Begin(); err != nil {
		return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	defer func() { _ = tx.Rollback() }()

	now = time.Now()
	row = tx.QueryRow(`
		INSERT INTO user(username, hashed_password, created_at, updated_at) 
		VALUES (?, ?, ?, ?) 
		RETURNING id, username, hashed_password, created_at, updated_at`,
		username, hashedPassword, now.Unix(), now.Unix())

	if user, err = scanUser(row); err != nil {
		var target sqlite3.Error
		if errors.As(err, &target) {
			if target.ExtendedCode == sqlite3.ErrConstraintUnique {
				return User{}, ErrUsernameTaken
			}
			if target.ExtendedCode == sqlite3.ErrConstraintCheck {
				_ = tx.Rollback()
				return User{}, nil
			}
		}
		return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if err = tx.Commit(); err != nil {
		return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return user, nil
}
