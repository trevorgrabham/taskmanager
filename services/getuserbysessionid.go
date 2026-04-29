package services

import (
	"errors"
	"fmt"
	sqlite "local/taskmanager/db"
)

// GetUserBySessionID returns the User associated with the SessionID.
//
// If SessionID is empty, and empty User is returned.
// If no valid entry for SessionID is found, then an ErrInvalidSessionID is returned.
// If an error occurred from the Repo, an ErrInternalRepo is returned.
func (s Service) GetUserBySessionID(sessionID string) (user User, err error) {
	var (
		caller   = "GetUserBySessionID"
		repoUser sqlite.User
	)
	if sessionID == "" {
		return User{}, nil
	}

	if repoUser, err = s.repo.GetUserBySessionID(sessionID); err != nil {
		if errors.Is(err, sqlite.ErrNoLiveSession) { return User{}, fmt.Errorf("%s: %w for %s", caller, ErrInvalidSessionID, sessionID) }
		return User{}, fmt.Errorf("%s: %w: %s", caller, ErrInternalRepo, err)
	}
	user = parseRepoUserToUser(repoUser)

	return user, nil
}
