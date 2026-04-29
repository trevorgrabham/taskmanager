// Package db implements functions for interacting with the business data Repo.
package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var dbFileName = "/home/trevorgrabham/.config/taskmanager/tasks.db"

var sessionExpirationDuration = 30 * 24 * time.Hour

var defaultTaskSelect = "SELECT task.id, title, category, description, due_date, completion_date, done, recurring_id, user_id, period FROM task LEFT JOIN recurring ON task.recurring_id = recurring.id"
var defaultUserSelect = "SELECT id, username, hashed_password, created_at, updated_at FROM user"

type Repo struct {
	db *sql.DB
}

func NewRepo() (r Repo, err error) {
	err = r.Connect()
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
	RecurringID     sql.NullInt64
	RecurringPeriod sql.NullString
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

func (r *Repo) isConnected() bool {
	return r.db != nil
}

// parseRecurringPeriod parses the recurringPeriod and adds the parsed period to the current time.
//
// If recurringPeriod is empty, an empty time.Time object is returned.
// If recurringPeriod is not in a recognized format, ErrInvalidRecurringPeriod is returned.
// If recurring value is not an integer, ErrInvalidRecurringValue is returned.
// If recurring unit is not one of 'days', 'weeks', 'months', ErrInvalidRecurringUnit is returned. 
func parseRecurringPeriod(recurringPeriod string) (nextDueDate time.Time, err error) {
	var (
		split []string
		recurringValue int
		recurringUnit string
		today time.Time
	)
	split = strings.Split(recurringPeriod, " ")
	if len(split) != 2 { return time.Time{}, ErrInvalidRecurringPeriod }
	if recurringValue, err = strconv.Atoi(split[0]); err != nil { return time.Time{}, ErrInvalidRecurringValue }

	today = time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	recurringUnit = split[1]
	switch recurringUnit {
	case "days":
		nextDueDate = today.AddDate(0, 0, recurringValue)
	case "weeks":
		nextDueDate = today.AddDate(0, 0, 7*recurringValue)
	case "months":
		nextDueDate = today.AddDate(0, recurringValue, 0)
	default:
		return time.Time{}, ErrInvalidRecurringUnit
	}

	return nextDueDate, nil
}

// insertRecurringPeriodIfNotExists returns the id of the record matching the 'recurringPeriod'. If no record exists, then one is created and its id is returned.
//
// An error is passed through if one occurs while scanning the matched record.
func insertRecurringPeriodIfNotExists(db dbConn, recurringPeriod string) (recurringID int, err error) {
	err = db.QueryRow(`INSERT INTO recurring(period) VALUES (?) ON CONFLICT(period) DO UPDATE SET period=excluded.period RETURNING id`, recurringPeriod).Scan(&recurringID)
	if err != nil {
		return -1, err
	}

	return recurringID, nil
}

// scanRow maps a *sql.Row object to a Task.
//
// Assumes the form of task.id, title, category, description, due_date, completion_date, done, recurring_id, user_id, period.
// If an error occurs while scanning the row, it is returned to the caller.
func scanRow(row scannable) (t Task, err error) {
	err = row.Scan(&t.ID, &t.Title, &t.Category, &t.Description, &t.DueDate, &t.CompletionDate, &t.Done, &t.RecurringID, &t.UserID, &t.RecurringPeriod)
	if err != nil {
		return Task{}, err
	}

	return t, nil
}

// scanRows maps a *sql.Rows object to []Task.
//
// Assumes the form of task.id, title, category, description, due_date, completion_date, done, recurring_id, user_id, period.
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

// deleteSessions removes all sessions with a matching SessionID in sessionsToDelete.
//
// If the Repo is not initialized an ErrNotConnected is returned.
// If sessionsToDelete is empty, nothing happens.
func (r *Repo) deleteSessions(sessionsToDelete []string) (err error) {
	if !r.isConnected() {
		return ErrNotConnected
	}
	if len(sessionsToDelete) < 1 {
		return nil
	}

	placeholder := strings.Repeat("?,", len(sessionsToDelete))
	placeholder = placeholder[:len(placeholder)-1] // trim the trailing comma
	args := make([]any, len(sessionsToDelete))
	for i, sessionID := range sessionsToDelete {
		args[i] = sessionID
	}

	if _, err = r.db.Exec(fmt.Sprintf(`DELETE FROM session WHERE id IN (%s)`, placeholder), args...); err != nil {
		return err
	}

	return nil
}
