package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/account"
)

func GetUserByUsername(db *sql.DB, username string) (user account.User, err error) {
	if err = checkDBConnection(db); err != nil {
		return user, err
	}
	if username == "" {
		return user, fmt.Errorf("GetUserByUsername: no username")
	}

	row := db.QueryRow(`SELECT id, username, hashed_password, created_at, updated_at FROM user WHERE username = ?`, username)

	return parseUser(row)
}

func GetUserByID(db *sql.DB, userID int) (user account.User, err error) {
	if err = checkDBConnection(db); err != nil {
		return user, err
	}
	if userID <= 0 {
		return user, fmt.Errorf("GetUserByID: no userID")
	}

	row := db.QueryRow(`SELECT id, username, hashed_password, created_at, updated_at FROM user WHERE id = ?`, userID)

	return parseUser(row)
}
