package services

import (
	"fmt"
)

// DeleteSession removes the session identified by sessionID from live sessions.
//
// If sessionID is empty, an ErrInvalidSessionID is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) DeleteSession(sessionID string) (err error) {
	if sessionID == "" {
		return ErrEmptySessionID
	}

	if err = s.repo.DeleteSession(sessionID); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return nil
}
