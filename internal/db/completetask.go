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
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec(`UPDATE task SET done = 1, completion_date = ? WHERE id = ?;`, time.Now().Unix(), taskToComplete.ID)
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	var period string
	err = tx.QueryRow(`SELECT period FROM recurring WHERE task_id = ?`, taskToComplete.ID).Scan(&period)
	if err == sql.ErrNoRows {
		_ = tx.Commit()
		return nil
	}
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	var (
		value    int
		unit     string
		split    []string
		timeDiff time.Duration
	)
	split = strings.Split(period, " ")
	value, err = strconv.Atoi(split[0])
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	unit = split[1]
	switch unit {
	case "days":
		timeDiff = time.Until(time.Now().AddDate(0, 0, value))
	case "weeks":
		timeDiff = time.Until(time.Now().AddDate(0, 0, value*7))
	case "months":
		timeDiff = time.Until(time.Now().AddDate(0, value, 0))
	default:
		return fmt.Errorf("completing task: unknown period format")
	}
	var res sql.Result
	res, err = tx.Exec(`INSERT INTO task(title, category, description, due_date) SELECT title, category, description, strftime('%s', 'now') + ? FROM task WHERE id = ?`, int(timeDiff.Seconds()), taskToComplete.ID)
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	var newID int64
	newID, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	_, err = tx.Exec(`UPDATE recurring SET task_id = ? WHERE task_id = ?`, int(newID), taskToComplete.ID)
	if err != nil {
		return fmt.Errorf("completing task: %s", err)
	}

	_ = tx.Commit()

	return nil
}
