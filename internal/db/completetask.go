package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"strconv"
	"strings"
	"time"
)

func CompleteTask(db *sql.DB, taskToComplete task.Task) error {
	if db == nil {
		return fmt.Errorf("completing task: cannot complete task for nil database")
	}
	var defaultTask task.Task
	if taskToComplete == defaultTask {
		return fmt.Errorf("completing task: cannot complete an empty task")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.Exec(`UPDATE task SET done = 1, completion_date = ? WHERE id = ?;`, time.Now().Unix(), taskToComplete.ID)
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	var (
		period      string
		recurringID int
	)
	err = tx.QueryRow(`SELECT id, period FROM recurring WHERE id = (SELECT recurring_id FROM task WHERE id = ?)`, taskToComplete.ID).Scan(&recurringID, &period)
	if err == sql.ErrNoRows {
		_ = tx.Commit()
		return nil
	}
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	var (
		value           int
		unit            string
		split           []string
		now, newDueDate time.Time
	)
	split = strings.Split(period, " ")
	value, err = strconv.Atoi(split[0])
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	unit = split[1]
	now = time.Now()
	switch unit {
	case "days":
		newDueDate = now.AddDate(0, 0, value)
		newDueDate = time.Date(newDueDate.Year(), newDueDate.Month(), newDueDate.Day(), 23, 59, 0, 0, time.Local)
	case "weeks":
		newDueDate = now.AddDate(0, 0, value*7)
		newDueDate = time.Date(newDueDate.Year(), newDueDate.Month(), newDueDate.Day(), 23, 59, 0, 0, time.Local)
	case "months":
		newDueDate = now.AddDate(0, value, 0)
		newDueDate = time.Date(newDueDate.Year(), newDueDate.Month(), newDueDate.Day(), 23, 59, 0, 0, time.Local)
	default:
		return fmt.Errorf("completing task: unknown period format")
	}
	_, err = tx.Exec(`INSERT INTO task(title, category, description, due_date, recurring_id) SELECT title, category, description, ?, ? FROM task WHERE id = ?`, newDueDate.Unix(), recurringID, taskToComplete.ID)
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	_ = tx.Commit()

	return nil
}
