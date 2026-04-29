package services

import (
	sqlite "local/taskmanager/db"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username  string
	Password  string
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// checkPassword compares a plain text password against a hashed version.
// Returns true if the passwords match. Returns false otherwise, including if either the plain text or hashed passwords are empty.
func checkPassword(plainPass, hashedPass string) bool {
	if plainPass == "" || hashedPass == "" {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass)); err != nil {
		return false
	}

	return true
}

// hashPassword generates a hashed password from the plain text string.
//
// Returns an error if the password length is too long.
func hashPassword(plainPass string) (string, error) {
	// TODO: validate the password length
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// parseRepoUserToUser maps repo data to a User object.
func parseRepoUserToUser(repoUser sqlite.User) (user User) {
	user.ID = repoUser.ID
	user.Username = repoUser.Username
	user.Password = repoUser.Password
	user.CreatedAt = time.Unix(repoUser.CreatedAt, 0)
	user.UpdatedAt = time.Unix(repoUser.UpdatedAt, 0)

	return user
}
