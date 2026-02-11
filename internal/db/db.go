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
var taskListTableName = "tasks_list"

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
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tasks_list (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT UNIQUE NOT NULL, 
	category TEXT, 
	times_completed INTEGER NOT NULL DEFAULT 0,
	last_time_completed INTEGER NOT NULL DEFAULT -1
);`)
	if err != nil {
		return fmt.Errorf("executing setup: %s", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	task_id INTEGER NOT NULL,
	description TEXT,
	due_date INTEGER NOT NULL,
	completion_date INTEGER NOT NULL DEFAULT -1,
	done INTEGER NOT NULL DEFAULT 0
		CHECK (done IN (0, 1)),
	FOREIGN KEY (task_id) REFERENCES tasks_list(id)
);`)
	if err != nil {
		return fmt.Errorf("executing setup: %s", err)
	}

	return nil
}

func TearDown(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("tearing down: cannot tear down a nil db")
	}

	_, err := db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, taskTableName))
	if err != nil {
		return fmt.Errorf("executing tear down: %s", err)
	}

	_, err = db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, taskListTableName))
	if err != nil {
		return fmt.Errorf("executing tear down: %s", err)
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
