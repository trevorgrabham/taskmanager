package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var dbFileName = "/home/trevorgrabham/.config/taskmanager/tasks.db"

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite db: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging sqlite db: %s", err)
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil { return nil, fmt.Errorf("setting up foreign keys: %s", err) }

	return db, nil
}
