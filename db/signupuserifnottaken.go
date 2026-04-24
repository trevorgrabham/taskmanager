package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (r *Repo) SignupUserIfNotTaken(username, password string) (user User, err error) {
	var (
		caller                = "SignupUseIfNotTaken"
		tx                    *sql.Tx
		row                   *sql.Row
		usernameAlreadyExists int
		now                   time.Time
	)
	if !r.isConnected() {
		return User{}, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	if tx, err = r.db.Begin(); err != nil {
		return User{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRow(`SELECT 1 FROM user WHERE username = ?`, username).Scan(&usernameAlreadyExists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return User{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}
	if usernameAlreadyExists == 1 {
		return User{}, fmt.Errorf("%s: %w", ErrUsernameTaken)
	}

	now = time.Now()
	row = tx.QueryRow(`
		INSERT INTO user(username, password, created_at, updated_at) 
		VALUES (?, ?, ?, ?) 
		RETURNING id, username, password, created_at, updated_at`,
		username, password, now.Unix(), now.Unix())

	if user, err = scanUser(row); err != nil {
		return User{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	if err = tx.Commit(); err != nil {
		return User{}, fmt.Errorf("%s: %w", caller, NewErrTransactionCommit(err))
	}

	return user, nil
}
