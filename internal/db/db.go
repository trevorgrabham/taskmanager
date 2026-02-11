package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"local/taskmanager2.0/internal/task"
)

var dbFileName = "/home/trevorgrabham/.config/taskmanager/tasks.db"
var taskTableName = "tasks"

type QueryParams struct {
	WhichTasks	task.WhichTasks
	From 	time.Time
	To 	time.Time
	Category string 
}

func Setup(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("setting up: cannot setup a nil db")
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
	if err != nil {
		return nil, fmt.Errorf("setting up foreign keys: %s", err)
	}

	return db, nil
}
