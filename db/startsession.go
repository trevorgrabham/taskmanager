package db

import (
	"database/sql"
	"errors"
	"fmt"
	"local/taskmanager/internal/cookies"
	"time"

	"github.com/mattn/go-sqlite3"
)

// StartSession creates a session for userID under sessionID.
//
// If sessionID is empty, returns ErrInternalRepo.
// If userID is empty or doesn't exist, returns ErrInternalRepo.
// If userID has another session, deletes the old session and creates a new one.
// If an error occurs while querying the database, an ErrInternalRepo is returned.
func (r *Repo) StartSession(sessionID string, userID int) (err error) {
	var (
		targetErr sqlite3.Error
		tx        *sql.Tx
	)
	if tx, err = r.db.Begin(); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	defer func(){_ = tx.Rollback()}()

	_, err = tx.Exec(`INSERT INTO session (id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)`, sessionID, userID, time.Now().Unix(), time.Now().Add(cookies.SessionTimeoutDuration).Unix())
	if err != nil {
		if errors.As(err, &targetErr) {
			if targetErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				if _, err = tx.Exec(`DELETE FROM session WHERE user_id = ?`, userID); err != nil {
					return fmt.Errorf("%w: %s", ErrInternalRepo, err)
				}

				_, err = tx.Exec(`
					INSERT INTO session (id, user_id, created_at, expires_at) 
					VALUES (?, ?, ?, ?)`, 
					sessionID, userID, time.Now().Unix(), time.Now().Add(cookies.SessionTimeoutDuration).Unix())
				if err != nil {
					return fmt.Errorf("%w: %s", ErrInternalRepo, err)
				}

				if err = tx.Commit(); err != nil { return fmt.Errorf("%w: %s", ErrInternalRepo, err) }

				return nil
			}
		}

		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if err = tx.Commit(); err != nil { return fmt.Errorf("%w: %s", ErrInternalRepo, err) }

	return nil
}
