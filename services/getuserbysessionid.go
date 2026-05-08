package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// GetUserBySessionID returns the User associated with the SessionID.
//
// If SessionID is empty, and empty User is returned.
// If no valid entry for SessionID is found, then an ErrInvalidSessionID is returned.
// If an error occurred from the Repo, an ErrInternalRepo is returned.
func (s Service) GetUserBySessionID(sessionID string) (user User, err error) {
	var repoUser sqlite.User
	if sessionID == "" {
		return User{}, ErrEmptySessionID
	}

	if repoUser, err = s.repo.GetUserBySessionID(sessionID); err != nil {
		return User{}, fmt.Errorf("%w: %s", ErrInternalRepo, err) 
	}
	if repoUser == (sqlite.User{}) { return User{}, ErrSessionNotFound }
	user = parseRepoUserToUser(repoUser)

	return user, nil
}
