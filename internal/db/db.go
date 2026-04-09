package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var dbFileName = "/home/trevorgrabham/.config/taskmanager/tasks.db"

func Setup(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("setting up: cannot setup a nil db")
	}

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS recurring (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	period TEXT UNIQUE NOT NULL
	);`)
	if err != nil {
		return fmt.Errorf("setting up: %s", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS task (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	category TEXT,
	description TEXT,
	due_date INTEGER,
	completion_date INTEGER,
	done INTEGER NOT NULL DEFAULT 0
		CHECK (done IN (0, 1)),
	recurring_id INTEGER,
	FOREIGN KEY (recurring_id) REFERENCES recurring(id)
);`)
	if err != nil {
		return fmt.Errorf("setting up: %s", err)
	}

	return nil
}

func TearDown(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("tearing down: cannot tear down a nil db")
	}

	_, err := db.Exec(`DROP TABLE IF EXISTS task;`)
	if err != nil {
		return fmt.Errorf("tearing down: %s", err)
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS recurring;`)
	if err != nil {
		return fmt.Errorf("tearing down: %s", err)
	}

	return nil
}

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, fmt.Errorf("connecting: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("connecting: %s", err)
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return nil, fmt.Errorf("connecting: %s", err)
	}

	return db, nil
}
