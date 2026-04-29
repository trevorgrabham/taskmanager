package db

import (
	"database/sql"
	"fmt"
)

// DeleteSession deletes the session identified by sessionID.
//
// If Repo is not initialized, returns an ErrNotConnected.
// If sessionID is empty, returns an ErrEmptySessionID.
// If no matching entry is found, an ErrNoLiveSession is returned.
// If an error occurs in the Repo, returns an ErrInternalRepo.
func (r Repo) DeleteSession(sessionID string) (err error) {
	var (
		res sql.Result
		rowsAffected int64
	)
	if !r.isConnected() { return ErrNotConnected }
	if sessionID == "" { return ErrEmptySessionID }

	if res, err = r.db.Exec(`DELETE FROM session WHERE id = ?`, sessionID); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	if rowsAffected < 1 { 
		return fmt.Errorf("%w %s", ErrNoLiveSession, sessionID)
	}

	return nil
}
