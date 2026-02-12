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

	var err error
	if newTask.Category == "" {
		_, err = db.Exec(fmt.Sprintf(`INSERT INTO %s(title) VALUES (?) ON CONFLICT(title) DO NOTHING`, taskListTableName), newTask.Title)
	} else {
		_, err = db.Exec(fmt.Sprintf(`INSERT INTO %s(title, category) VALUES (?, ?) ON CONFLICT(title) DO NOTHING`, taskListTableName), newTask.Title, newTask.Category)
	}
	if err != nil {
		return fmt.Errorf("adding task: %s", err)
	}

	var taskID int
	row := db.QueryRow(`SELECT id FROM tasks_list WHERE title = ?`, newTask.Title)
	err = row.Scan(&taskID)
	if err != nil {
		return fmt.Errorf("adding task: %s", err)
	}

	// insertStatement, err := db.Prepare(fmt.Sprintf(`INSERT INTO %s(title, category, description, due_date, completion_date, done) VALUES (?, ?, ?, ?, ?, ?)`, taskTableName))
	// if err != nil {
	// return fmt.Errorf("preparing insert stmnt: %s", err)
	// }

	if newTask.Done {
		if newTask.Description == "" {
			_, err = db.Exec(fmt.Sprintf(`INSERT INTO %s(task_id, due_date, completion_date, done) VALUES (?, ?, ?, 1)`, taskListTableName), taskID, time.Time(newTask.DueDate).Unix(), time.Time(newTask.CompletionDate).Unix())
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}
		} else {
			_, err = db.Exec(fmt.Sprintf(`INSERT INTO %s(task_id, description, due_date, completion_date, done) VALUES (?, ?, ?, 1)`, taskListTableName), taskID, newTask.Description, time.Time(newTask.DueDate).Unix(), time.Time(newTask.CompletionDate).Unix())
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}
		}
	} else {
		if newTask.Description == "" {
			_, err = db.Exec(fmt.Sprintf(`INSERT INTO %s(task_id, due_date) VALUES (?, ?)`, taskListTableName), taskID, time.Time(newTask.DueDate).Unix())
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}
		} else {
			_, err = db.Exec(fmt.Sprintf(`INSERT INTO %s(task_id, description, due_date) VALUES (?, ?, ?)`, taskListTableName), taskID, newTask.Description, time.Time(newTask.DueDate).Unix())
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}
		}
	}

	return nil
}
