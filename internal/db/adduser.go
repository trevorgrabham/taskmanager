package db

import (
	"database/sql"
)

func AddUser(db *sql.DB, username, hashedPassword string) (err error) {
	if err = checkDBConnection(db); err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO user(username, hashed_password) VALUES (?, ?)`, username, hashedPassword)
	return err
}
