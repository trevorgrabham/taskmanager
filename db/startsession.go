package db

import (
	"database/sql"
	"fmt"
	"time"
)

// Deletes any other sessions for the user

// StartSession creates a session for userID under sessionID.
//
// If the Repo is not initialized an ErrNotConnected is returned.
// If sessionID is empty, an ErrEmptySessionID is returned.
// If userID is empty, an ErrUserNotExist is returned.
// If an error occurs while querying the database, an ErrInternalRepo is returned.
func (r *Repo) StartSession(sessionID string, userID int) (err error) {
	var (
		rows             *sql.Rows
		tx               *sql.Tx
		idToDelete       string
		sessionsToDelete []string
		now              time.Time
	)
	if !r.isConnected() {
		return ErrNotConnected
	}

	if tx, err = r.db.Begin(); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	// delete any old sessions for the user
	rows, err = tx.Query(`SELECT id FROM session WHERE user_id = ?`, userID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	for rows.Next() {
		if err = rows.Scan(&idToDelete); err != nil {
			return fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}

		sessionsToDelete = append(sessionsToDelete, idToDelete)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if err = r.deleteSessions(sessionsToDelete); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	now = time.Now()
	_, err = tx.Exec(`
		INSERT INTO session (id, user_id, created_at, expires_at) 
		VALUES (?, ?, ?, ?)`,
		sessionID, userID, now.Unix(), now.Add(sessionExpirationDuration).Unix())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%w: %s", ErrTransactionCommit, err)
	}

	return nil
}
