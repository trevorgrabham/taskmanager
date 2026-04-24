package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

func (s Service) SignupUser(username, password string) (sessionID string, err error) {
	var (
		caller         = "SignupUser"
		hashedPassword string
		repoUser       sqlite.User
	)
	if username == "" {
		return "", fmt.Errorf("%s: %w", caller, ErrUserNoUsername)
	}
	if password == "" {
		return "", fmt.Errorf("%s: %w", caller, ErrUserNoPassword)
	}

	if hashedPassword, err = hashPassword(password); err != nil {
		return "", NewErrBcrypt(err)
	}

	if repoUser, err = s.repo.SignupUserIfNotTaken(username, hashedPassword); err != nil {
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	if sessionID, err = s.repo.StartSession(repoUser.ID); err != nil {
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	return sessionID, nil
}
