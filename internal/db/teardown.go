package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func TearDown(db *sql.DB) error {
	if db == nil {
		return errors.New("tearing down: cannot tear down a nil db")
	}

	_, err := db.Exec(`DROP TABLE IF EXISTS tasks;`)
	if err != nil {
		return fmt.Errorf("executing tear down: %s", err)
	}

	return nil
}
