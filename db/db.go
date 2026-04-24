package db

import (
	"database/sql"
	"fmt"
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

func insertRecurringPeriodIfNotExists(db dbConn, recurringPeriod string) (recurringID int, err error) {
	err = db.QueryRow(`INSERT INTO recurring(period) VALUES (?) ON CONFLICT(period) DO UPDATE SET period=excluded.period RETURNING id`, recurringPeriod).Scan(&recurringID)
	if err != nil {
		return -1, fmt.Errorf("for recurring period %s: %w", recurringPeriod, err)
	}

	return recurringID, nil
}

// Assumes the form of task.id, title, category, description, due_date, completion_date, done, recurring_id, user_id, period
func scanRow(row scannable) (t Task, err error) {
	if row == nil {
		return Task{}, ErrNilSQLRow
	}

	err = row.Scan(&t.ID, &t.Title, &t.Category, &t.Description, &t.DueDate, &t.CompletionDate, &t.Done, &t.RecurringID, &t.UserID, &t.RecurringPeriod)
	if err != nil {
		return Task{}, err
	}

	return t, nil
}

func scanRows(rows *sql.Rows) (tasks []Task, err error) {
	if rows == nil {
		return nil, ErrNilSQLRow
	}

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

// Assumes the form of username, hashed_password, created_at, updated_at
func scanUser(row *sql.Row) (user User, err error) {
	if row == nil {
		return User{}, ErrNilSQLRow
	}

	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r *Repo) deleteSessions(sessionsToDelete []string) (err error) {
	if !r.isConnected() {
		return ErrNotConnected
	}

	for _, sessionID := range sessionsToDelete {
		if _, err = r.db.Exec(`DELETE FROM session WHERE id = ?`, sessionID); err != nil {
			return err
		}
	}

	return nil
}
