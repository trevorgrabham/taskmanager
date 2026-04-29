package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// SignupUserIfNotTaken inserts a record for the given username and password in the Repo. It is the callers responsibility to ensure that the passed password has been hashed. 
//
// If the Repo is not initialized an ErrNotConnected is returned.
// If username is empty, an ErrInvalidUsername.
// If password is empty, an ErrInvalidPassword.
// If username already exists in the Repo, an ErrUsernameExists is returned.
// If an error occurs while querying the database, an ErrInternalRepo is returned.
func (r *Repo) SignupUserIfNotTaken(username, hashedPassword string) (user User, err error) {
	var (
		tx                    *sql.Tx
		row                   *sql.Row
		usernameAlreadyExists int
		now                   time.Time
	)
	if !r.isConnected() { return User{}, ErrNotConnected }
	if username == "" { return User{}, ErrInvalidUsername }
	if hashedPassword == "" { return User{}, ErrInvalidPassword }

	if tx, err = r.db.Begin(); err != nil { return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err) }
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRow(`SELECT 1 FROM user WHERE username = ?`, username).Scan(&usernameAlreadyExists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { 
 		return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err) 
	}
	if usernameAlreadyExists == 1 { return User{}, ErrUsernameExists }

	now = time.Now()
	row = tx.QueryRow(`
		INSERT INTO user(username, hashed_password, created_at, updated_at) 
		VALUES (?, ?, ?, ?) 
		RETURNING id, username, hashed_password, created_at, updated_at`,
		username, hashedPassword, now.Unix(), now.Unix())

	if user, err = scanUser(row); err != nil { return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err) }

	if err = tx.Commit(); err != nil { return User{}, fmt.Errorf("%w: %s", ErrTransactionCommit, err) }

	return user, nil
}
