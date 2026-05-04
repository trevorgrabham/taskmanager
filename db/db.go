// Package db implements functions for interacting with the business data Repo.
package db

import (
	"database/sql"
	"time"
)

var dbFileName = "/home/trevorgrabham/.config/taskmanager/tasks.db"

var sessionExpirationDuration = 30 * 24 * time.Hour

var defaultTaskSelect = "SELECT task.id, title, category, description, due_date, completion_date, done, recurring_period, user_id FROM task"
var defaultUserSelect = "SELECT id, username, hashed_password, created_at, updated_at FROM user"

type Repo struct {
	db *sql.DB
}

func NewRepo(file string) (r Repo, err error) {
	if file == "" {
		file = dbFileName
	}
	err = r.Connect(file)
	return r, err
}

type Task struct {
	ID              int
	UserID          int
	Title           string
	Category        sql.NullString
	Description     sql.NullString
	DueDate         sql.NullInt64
	CompletionDate  sql.NullInt64
	Done            bool
	RecurringPeriod sql.NullInt64
}

type User struct {
	ID        int
	Username  string
	Password  string
	CreatedAt int64
	UpdatedAt int64
}

type dbConn interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type scannable interface {
	Scan(...any) error
}

// scanRow maps a *sql.Row object to a Task.
//
// Assumes the form of task.id, title, category, description, due_date, completion_date, done, recurring_period, user_id.
// If an error occurs while scanning the row, it is returned to the caller.
func scanRow(row scannable) (t Task, err error) {
	err = row.Scan(&t.ID, &t.Title, &t.Category, &t.Description, &t.DueDate, &t.CompletionDate, &t.Done, &t.RecurringPeriod, &t.UserID)
	if err != nil {
		return Task{}, err
	}

	return t, nil
}

// scanRows maps a *sql.Rows object to []Task.
//
// Assumes the form of task.id, title, category, description, due_date, completion_date, done, recurring_period, user_id.
// If an error occurs while scanning the row, it is returned to the caller.
func scanRows(rows *sql.Rows) (tasks []Task, err error) {
	var t Task
	for rows.Next() {
		if t, err = scanRow(rows); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// scanUser maps an *sql.Row object to a User.
//
// Assumes the form of username, hashed_password, created_at, updated_at.
// If an error occurs while scanning the row, it is returned to the caller.
func scanUser(row *sql.Row) (user User, err error) {
	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}
