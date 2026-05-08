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
		repoUser      sqlite.User
		passwordMatch bool
	)
	if username == "" {
		return "", &ErrUserValidation{
			Field:   UserFieldUsername,
			Message: MessageInvalidUsernameOrPassword,
		}
	}
	if password == "" {
		return "", &ErrUserValidation{
			Field:   UserFieldPassword,
			Message: MessageInvalidUsernameOrPassword,
		}
	}

	if repoUser, err = s.repo.GetUserByUsername(username); err != nil {
		return "", fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	if repoUser == (sqlite.User{}) {
		return "", &ErrUserValidation{
			Field:   UserFieldUsername,
			Message: MessageInvalidUsernameOrPassword,
		}
	}

	if passwordMatch = checkPassword(password, repoUser.Password); !passwordMatch {
		return "", &ErrUserValidation{
			Field:   UserFieldPassword,
			Message: MessageInvalidUsernameOrPassword,
		}
	}

	sessionID = GenerateSessionID()

	if err = s.repo.StartSession(sessionID, repoUser.ID); err != nil {
		return "", fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return sessionID, nil
}
