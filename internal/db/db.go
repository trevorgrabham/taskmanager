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
