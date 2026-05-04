package db

import (
	"fmt"
)

// DeleteSession deletes the session identified by sessionID.
//
// If sessionID is empty or doesn't exist, it is treated as a no-op.
// If a transient error occurs in the Repo, returns ErrInternalRepo.
func (r Repo) DeleteSession(sessionID string) (err error) {
	if _, err = r.db.Exec(`DELETE FROM session WHERE id = ?`, sessionID); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return nil
}
