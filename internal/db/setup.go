package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func Setup(db *sql.DB) error {
	if db == nil {
		return errors.New("setting up: cannot setup a nil db")
	}
	_, err := db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL, 
	category TEXT, 
	description TEXT,
	due_date INTEGER NOT NULL,
	completion_date INTEGER NOT NULL DEFAULT -1,
	done INTEGER NOT NULL DEFAULT 0
		CHECK (done IN (0, 1))
);`, taskTableName))
	if err != nil {
		return fmt.Errorf("executing setup: %s", err)
	}

	return nil
}
