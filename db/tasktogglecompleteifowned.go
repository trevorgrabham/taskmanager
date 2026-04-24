package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (r *Repo) TaskToggleCompleteIfOwned(taskID, userID int) (task Task, err error) {
	var (
		caller = "TaskToggleCompleteIfOwned"
		tx     *sql.Tx
		row    *sql.Row
    split []string
    recurringValue int
    recurringUnit string 
    today time.Time
    nextTask Task
	)
	if !r.isConnected() {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	if tx, err = r.db.Begin(); err != nil {
		return Task{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}
	defer func() { _ = tx.Rollback() }()

	row = tx.QueryRow(`
    UPDATE task 
    SET done = (done + 1) % 2, completion_date = 
      CASE 
        WHEN completion_date IS NULL THEN ?
        ELSE NULL
      END
    WHERE id = ? AND user_id = ?
    RETURNING id, title, category, description, due_date, completion_date, done, recurring_id, user_id,  (
      SELECT period FROM recurring WHERE recurring.id = task.recurring_id
    )`, taskID, userID)

  if task, err = scanRow(row); err != nil { 
    if errors.Is(err, sql.ErrNoRows) { return Task{}, NewErrNotOwner(userID, taskID) }
    return Task{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err)) 
  }

  // Not recurring, nothing more to do
  if !task.RecurringID.Valid || task.RecurringID.Int64 == 0 {
    if err = tx.Commit(); err != nil { return Task{}, fmt.Errorf("%s: %w", caller, NewErrTransactionCommit(err)) }

    return task, nil
  }
  // Just uncompleted a task. There is still a task in the database that was added when the recurring task was initially completed, but attempting to remove it may cause errors.
  // We may consider adding a "recurring_parent" field that references the taskID for recurring tasks, but until then, we don't want to accidentally remove the wrong task
  if !task.Done {
    if err = tx.Commit(); err != nil { return Task{}, fmt.Errorf("%s: %w", caller, NewErrTransactionCommit(err)) }

    return task, nil
  }

  // If we completed a recurring task, add the next one
  nextTask = task
  nextTask.Done = false
  nextTask.CompletionDate.Valid = false
  split = strings.Split(task.RecurringPeriod.String, " ")
  if len(split) != 2 { return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskBadRecurringPeriod) }
  if recurringValue, err = strconv.Atoi(split[0]); err != nil { return Task{}, fmt.Errorf("%s: %w", caller, err) }

  today = time.Now()
  today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
  recurringUnit = split[1]
  switch recurringUnit {
  case "days":
    nextTask.DueDate.Valid = true
    nextTask.DueDate.Int64 = today.AddDate(0, 0, recurringValue).Unix()
  case "weeks":
    nextTask.DueDate.Valid = true
    nextTask.DueDate.Int64 = today.AddDate(0, 0, 7 * recurringValue).Unix()
  case "months":
    nextTask.DueDate.Valid = true
    nextTask.DueDate.Int64 = today.AddDate(0, recurringValue, 0).Unix()
  default:
    return Task{}, fmt.Errorf("%s: %w", caller, ErrTaskBadRecurringPeriod)
  }

  row = tx.QueryRow(`
    INSERT INTO task(title, category, description, due_date, recurring_id, user_id) 
    VALUES (?, ?, ?, ?, ?, ?)
    RETURNING task.id, title, category, description, due_date, completion_date, done, recurring_id, user_id, (
      SELECT period FROM recurring WHERE recurring.id = task.recurring_id)`, 
    nextTask.Title, nextTask.Category, nextTask.Description, nextTask.DueDate, nextTask.RecurringID, nextTask.UserID)

  if nextTask, err = scanRow(row); err != nil { return Task{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err)) }

  if err = tx.Commit(); err != nil { return Task{}, fmt.Errorf("%s: %w", caller, NewErrTransactionCommit(err)) }

  return task, nil
}
