package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// LoginUser validates the the username and password against the account Repo.
//
// If username is empty, an ErrInvalidUsername is returned.
// If password is empty, an ErrInvalidPassword is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) LoginUser(username, password string) (sessionID string, err error) {
	var (
		caller        = "LoginUser"
		repoUser      sqlite.User
		passwordMatch bool
	)
	if username == "" {
		return "", fmt.Errorf("%s: %w", caller, ErrInvalidUsername)
	}
	if password == "" {
		return "", fmt.Errorf("%s: %w", caller, ErrInvalidPassword)
	}

	if repoUser, err = s.repo.GetUserByUsername(username); err != nil {
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	if passwordMatch = checkPassword(password, repoUser.Password); !passwordMatch {
		return "", fmt.Errorf("%s: %w", caller, ErrInvalidPassword)
	}

	sessionID = GenerateSessionID()

	if err = s.repo.StartSession(sessionID, repoUser.ID); err != nil {
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	return sessionID, nil
}
