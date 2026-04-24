package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

func (s Service) LoginUser(username, password string) (sessionID string, err error) {
	var (
		caller        = "LoginUser"
		repoUser      sqlite.User
		passwordMatch bool
	)
	if username == "" {
		return "", fmt.Errorf("%s: %w", caller, ErrUserNoUsername)
	}
	if password == "" {
		return "", fmt.Errorf("%s: %w", caller, ErrUserNoPassword)
	}

	if repoUser, err = s.repo.GetUserByUsername(username); err != nil {
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	if passwordMatch = checkPassword(password, repoUser.Password); !passwordMatch {
		return "", ErrUserWrongPassword
	}

	if sessionID, err = s.repo.StartSession(repoUser.ID); err != nil {
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	return sessionID, nil
}
