package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Deletes any other sessions for the user
func (r *Repo) StartSession(userID int) (sessionID string, err error) {
	var (
		caller           = "StartSession"
		rows             *sql.Rows
		tx               *sql.Tx
		idToDelete string
		sessionsToDelete []string
		now time.Time
	)
	if !r.isConnected() {
		return "", fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	if tx, err = r.db.Begin(); err != nil {
		return "", fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	rows, err = tx.Query(`SELECT id FROM session WHERE user_id = ?`, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	if err == nil {
		for rows.Next() {
			if err = rows.Scan(&idToDelete); err != nil {
				return "", fmt.Errorf("%s: %w", caller, NewErrRepo(err))
			}

			sessionsToDelete = append(sessionsToDelete, idToDelete)
		}
		if err = rows.Err(); err != nil {
			return "", fmt.Errorf("%s: %w", caller, NewErrRepo(err))
		}

		if err = r.deleteSessions(sessionsToDelete); err != nil {
			return "", fmt.Errorf("%s: %w", caller, NewErrRepo(err))
		}
	}

	now = time.Now()
	err = tx.QueryRow(`
		INSERT INTO session (user_id, created_at, expires_at) 
		VALUES (?, ?, ?)
		RETURNING id`,
		userID, now.Unix(), now.Add(sessionExpirationDuration).Unix()).Scan(&sessionID)
	if err != nil { return "", fmt.Errorf("%s: %w", caller, NewErrRepo(err)) }

	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("%s: %w", caller, NewErrTransactionCommit(err)) 
	}

	return sessionID, nil
}
