package services

import (
	"fmt"
	sqlite "local/taskmanager/db"
)

// SignupUser validates the the username and password and creates a User account on success.
//
// If username is empty, an ErrInvalidUsername is returned.
// If password is empty, an ErrInvalidPassword is returned.
// If the username is already taken, an ErrUsernameExists is returned.
// If an error occurs during encryption, an ErrEncrypt is returned.
// If an error occurrs in the Repo, an ErrInternalRepo is returned.
func (s Service) SignupUser(username, password string) (sessionID string, err error) {
	var (
		hashedPassword string
		repoUser       sqlite.User
	)
	if username == "" {
		return "", &ErrUserValidation{
			Field:   UserFieldUsername,
			Message: MessageInvalidUsername,
		}
	}
	if password == "" {
		return "", &ErrUserValidation{
			Field:   UserFieldPassword,
			Message: MessageInvalidPassword,
		}
	}

	if hashedPassword, err = hashPassword(password); err != nil {
		return "", fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if repoUser, err = s.repo.SignupUserIfNotTaken(username, hashedPassword); err != nil {
		switch sqlite.MapErrorToTypeName(err) {
		case "ErrUsernameTaken":
			return "", &ErrUserValidation{
				Field: UserFieldUsername,
				Message: MessageUsernameTaken,
			}
		default:
			return "", fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}
	}

	sessionID = GenerateSessionID()

	if err = s.repo.StartSession(sessionID, repoUser.ID); err != nil {
		return "", fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return sessionID, nil
}
