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

	return db, nil
}
