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

	_, err := db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, taskTableName))
	if err != nil {
		return fmt.Errorf("executing tear down: %s", err)
	}

	return nil
}
