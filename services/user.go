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

func checkPassword(plainPass, hashedPass string) bool {
	if plainPass == "" || hashedPass == "" {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass)); err != nil {
		return false
	}

	return true
}

func hashPassword(plainPass string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func parseRepoUserToUser(repoUser sqlite.User) (user User) {
	user.ID = repoUser.ID
	user.Username = repoUser.Username
	user.Password = repoUser.Password
	user.CreatedAt = time.Unix(repoUser.CreatedAt, 0)
	user.UpdatedAt = time.Unix(repoUser.UpdatedAt, 0)

	return user
}
