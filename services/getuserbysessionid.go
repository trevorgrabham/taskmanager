package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

func (s Service) GetUserBySessionID(sessionID string) (user User, err error) {
	var (
		caller   = "GetUserBySessionID"
		repoUser sqlite.User
	)
	if sessionID == "" {
		return User{}, fmt.Errorf("%s: %w", caller, ErrSessionNoID)
	}

	if repoUser, err = s.repo.GetUserBySessionID(sessionID); err != nil {
		return User{}, fmt.Errorf("%s: %w", caller, err)
	}
	user = parseRepoUserToUser(repoUser)

	return user, nil
}
