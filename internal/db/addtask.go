package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"strings"
)

func AddTask(db *sql.DB, newTask task.Task) (id int, err error) {
	if db == nil {
		return -1, fmt.Errorf("adding task: cannot add task to nil sqlite database")
	}
	if newTask.IsZero() {
		return -1, fmt.Errorf("adding task: cannot add empty task to sqlite database")
	}

	queryColumns := []string{"title"}
	queryPlaceholders := []string{"?"}
	queryArgs := []any{newTask.Title}

	if newTask.Category != "" {
		queryColumns = append(queryColumns, "category")
		queryPlaceholders = append(queryPlaceholders, "?")
		queryArgs = append(queryArgs, newTask.Category)
	}

	if newTask.Description != "" {
		queryColumns = append(queryColumns, "description")
		queryPlaceholders = append(queryPlaceholders, "?")
		queryArgs = append(queryArgs, newTask.Description)
	}

	if !newTask.DueDate.IsZero() {
		queryColumns = append(queryColumns, "due_date")
		queryPlaceholders = append(queryPlaceholders, "?")
		queryArgs = append(queryArgs, newTask.DueDate.Unix())
	}

	var res sql.Result
	res, err = db.Exec(fmt.Sprintf(`INSERT INTO task(%s) VALUES (%s)`, strings.Join(queryColumns, ", "), strings.Join(queryPlaceholders, ", ")), queryArgs...)
	if err != nil {
		return -1, fmt.Errorf("adding task: %s", err)
	}

	var lastID int64
	lastID, err = res.LastInsertId()
	if err != nil { return -1, fmt.Errorf("adding task: %s", err) }

	id = int(lastID)
	return id, nil
}
