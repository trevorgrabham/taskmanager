package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager2.0/internal/task"
	"time"
)

func AddTask(db *sql.DB, newTask task.Task) error {
	if db == nil {
		return fmt.Errorf("adding task: cannot add task to nil sqlite database")
	}

	var emptyTask task.Task
	if newTask == emptyTask {
		return fmt.Errorf("adding task: cannot add empty task to sqlite database")
	}

	insertStatement, err := db.Prepare(`INSERT INTO tasks(title, category, description, due_date, completion_date, done) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("preparing insert stmnt: %s", err)
	}

	switch {
	case newTask.Category == "" && newTask.Description == "" && !newTask.Done:
		_, err = insertStatement.Exec(newTask.Title, nil, nil, time.Time(*newTask.DueDate).Unix(), -1, 0)
		if err != nil {
			return fmt.Errorf("executing insert stmnt: %s", err)
		}
	case newTask.Category == "" && newTask.Description == "":
		_, err = insertStatement.Exec(newTask.Title, nil, nil, time.Time(*newTask.DueDate).Unix(), time.Time(*newTask.CompletionDate).Unix(), 1)
		if err != nil {
			return fmt.Errorf("executing insert stmnt: %s", err)
		}
	case newTask.Category == "" && !newTask.Done:
		_, err = insertStatement.Exec(newTask.Title, nil, newTask.Description, time.Time(*newTask.DueDate).Unix(), -1, 0)
		if err != nil {
			return fmt.Errorf("executing insert stmnt: %s", err)
		}
	case newTask.Description == "" && !newTask.Done:
		_, err = insertStatement.Exec(newTask.Title, newTask.Category, nil, time.Time(*newTask.DueDate).Unix(), -1, 0)
		if err != nil {
			return fmt.Errorf("executing insert stmnt: %s", err)
		}
	case newTask.Category == "":
		_, err = insertStatement.Exec(newTask.Title, nil, newTask.Description, time.Time(*newTask.DueDate).Unix(), time.Time(*newTask.CompletionDate).Unix(), 1)
		if err != nil {
			return fmt.Errorf("executing insert stmnt: %s", err)
		}
	case newTask.Description == "":
		_, err = insertStatement.Exec(newTask.Title, newTask.Category, nil, time.Time(*newTask.DueDate).Unix(), time.Time(*newTask.CompletionDate).Unix(), 1)
		if err != nil {
			return fmt.Errorf("executing insert stmnt: %s", err)
		}
	case !newTask.Done:
		_, err = insertStatement.Exec(newTask.Title, newTask.Category, newTask.Description, time.Time(*newTask.DueDate).Unix(), -1, 0)
		if err != nil {
			return fmt.Errorf("executing insert stmnt: %s", err)
		}
	default:
		_, err = insertStatement.Exec(newTask.Title, newTask.Category, newTask.Description, time.Time(*newTask.DueDate).Unix(), time.Time(*newTask.CompletionDate).Unix(), 1)
		if err != nil {
			return fmt.Errorf("executing insert stmnt: %s", err)
		}
	}

	return nil
}
