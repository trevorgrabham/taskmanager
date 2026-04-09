package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
)

func TaskByID(db *sql.DB, id int) (task.Task, error) {
	if db == nil {
		return task.Task{}, fmt.Errorf("task by id: cannot get a task from a nil database")
	}
	if id < 1 {
		return task.Task{}, fmt.Errorf("task by id: cannot get a task without an id")
	}

	row := db.QueryRow(`
		SELECT task.id, title, category, description, due_date, completion_date, done, period
		FROM task 
		LEFT JOIN recurring ON task.recurring_id = recurring.id
		WHERE task.id = ?`, id)

	t, err := scanTaskRow(row)
	if err != nil {
		return task.Task{}, fmt.Errorf("task by id: %s", err)
	}

	return t, nil
}
