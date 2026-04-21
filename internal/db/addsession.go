package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/account"
	"time"
)

func AddSession(db *sql.DB, userID int, sessionID string) (err error) {
	if err = checkDBConnection(db); err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO session(id, user_id, expires_at) VALUES (?, ?, ?)`, sessionID, userID, time.Now().Add(account.SessionTimeoutDuration).Unix())
	if err != nil {
		return fmt.Errorf("AddSession: %s", err)
	}

	return nil
}
