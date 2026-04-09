package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"strconv"
	"strings"
	"time"
)

func ToggleTaskComplete(db *sql.DB, taskToToggle task.Task) (updatedTask, newTask task.Task, err error) {
	if db == nil {
		return task.Task{}, task.Task{}, fmt.Errorf("toggle task: cannot toggle task completion for a nil database")
	}
	if taskToToggle.IsZero() {
		return task.Task{}, task.Task{}, fmt.Errorf("toggle task: cannot toggle task for an empty task")
	}

	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		return task.Task{}, task.Task{}, fmt.Errorf("toggle task: %s", err)
	}
	defer func() { _ = tx.Rollback() }()

	row := tx.QueryRow(`
		UPDATE task 
		SET done = (done + 1) % 2, completion_date = 
			CASE 
				WHEN completion_date IS NULL THEN ? 
				ELSE NULL
			END
		WHERE id = ?
		RETURNING id, title, category, description, due_date, completion_date, done, (SELECT period FROM recurring WHERE id = task.recurring_id)
	`, time.Now().Unix(), taskToToggle.ID)
	updatedTask, err = scanTaskRow(row)
	if err != nil {
		return task.Task{}, task.Task{}, fmt.Errorf("toggle task: %s", err)
	}

	if updatedTask.RecurringPeriod != "" {
		queryColumns := []string{"title"}
		queryPlaceholders := []string{"?"}
		queryArgs := []any{updatedTask.Title}

		if updatedTask.Category != "" {
			queryColumns = append(queryColumns, "category")
			queryPlaceholders = append(queryPlaceholders, "?")
			queryArgs = append(queryArgs, updatedTask.Category)
		}

		if updatedTask.Description != "" {
			queryColumns = append(queryColumns, "description")
			queryPlaceholders = append(queryPlaceholders, "?")
			queryArgs = append(queryArgs, updatedTask.Description)
		}

		switch updatedTask.Done {
		case true: // Just completed a recurring task. Add a new one to task table
			var (
				newDueDate         time.Time
				value, recurringID int
			)
			err = tx.QueryRow(`SELECT recurring_id FROM task WHERE id = ?`, updatedTask.ID).Scan(&recurringID)
			if err != nil {
				return task.Task{}, task.Task{}, fmt.Errorf("toggle task: %s", err)
			}
			queryColumns = append(queryColumns, "recurring_id")
			queryPlaceholders = append(queryPlaceholders, "?")
			queryArgs = append(queryArgs, recurringID)

			split := strings.Split(updatedTask.RecurringPeriod, " ")
			valueString, unit := split[0], split[1]
			value, err = strconv.Atoi(valueString)
			if err != nil {
				return task.Task{}, task.Task{}, fmt.Errorf("toggle task: %s", err)
			}

			switch unit {
			case "days":
				newDueDate = time.Now().AddDate(0, 0, value)
			case "weeks":
				newDueDate = time.Now().AddDate(0, 0, value*7)
			case "months":
				newDueDate = time.Now().AddDate(0, value, 0)
			}
			newDueDate = time.Date(newDueDate.Year(), newDueDate.Month(), newDueDate.Day(), 23, 59, 0, 0, time.Local)

			if !updatedTask.DueDate.IsZero() {
				queryColumns = append(queryColumns, "due_date")
				queryPlaceholders = append(queryPlaceholders, "?")
				queryArgs = append(queryArgs, newDueDate.Unix())
			}

			var res sql.Result
			res, err = tx.Exec(fmt.Sprintf(`INSERT INTO task(%s) VALUES (%s)`, strings.Join(queryColumns, ", "), strings.Join(queryPlaceholders, ", ")), queryArgs...)
			if err != nil {
				return task.Task{}, task.Task{}, fmt.Errorf("toggle task: %s", err)
			}

			var newID int64
			newID, err = res.LastInsertId()
			if err != nil {
				return task.Task{}, task.Task{}, fmt.Errorf("toggle task: %s", err)
			}

			newTask = task.Task{
				ID:              int(newID),
				Title:           updatedTask.Title,
				Category:        updatedTask.Category,
				Description:     updatedTask.Description,
				DueDate:         task.TaskDueDate(newDueDate),
				RecurringPeriod: updatedTask.RecurringPeriod,
			}

			_ = tx.Commit()

			return updatedTask, newTask, nil
		case false: // Just uncompleted a recurring task. Need to update the original, and then remove any tasks that match the title, category, description, recurring_id
			zipped := make([]any, 0, len(queryColumns)*2)
			for i := range queryColumns {
				queryPlaceholders[i] = "? = ?"
				zipped = append(zipped, queryColumns[i], queryArgs[i])
			}
			zipped = append(zipped, updatedTask.RecurringPeriod)

			_, err = tx.Exec(fmt.Sprintf(`DELETE FROM task WHERE %s AND recurring_id = (SELECT id FROM recurring WHERE period = ?`, strings.Join(queryPlaceholders, " AND ")), zipped...)
		}
	}

	_ = tx.Commit()

	return updatedTask, task.Task{}, nil
}
