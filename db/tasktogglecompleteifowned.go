package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TaskToggleCompleteIfOwned toggles the completion status for the task identified by taskID if it is owned by userID.
//
// If Repo is not initialized, returns an ErrNotConnected.
// If taskID is empty or not task matches, returns an ErrTaskNotExist.
// If userID is empty, returns an ErrUserNotExist.
// If an error occurs in the Repo, returns an ErrInternalRepo.
// If the task is not owned by userID, returns an ErrNotOwner.
// If no task was not toggled, but the task exists and is owned by userID, returns an ErrorUnknown. This should never happen.
func (r *Repo) TaskToggleCompleteIfOwned(taskID, userID int) (task Task, err error) {
	var (
		tx           *sql.Tx
		res          sql.Result
		rowsAffected int64
		row          *sql.Row
		nextDueDate  time.Time
		nextTask     Task
	)
	if !r.isConnected() {
		return Task{}, ErrNotConnected
	}
	if taskID < 1 {
		return Task{}, fmt.Errorf("%w for id %d", ErrTaskNotExist, taskID)
	}
	if userID < 1 {
		return Task{}, fmt.Errorf("%w for id %d", ErrUserNotExist, userID)
	}

	if tx, err = r.db.Begin(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err = tx.Exec(`
    UPDATE task 
    SET done = (done + 1) % 2, completion_date = 
      CASE 
        WHEN completion_date IS NULL THEN ?
        ELSE NULL
      END
    WHERE id = ? AND user_id = ?`, time.Now().Unix(), taskID, userID)
	if err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	row = tx.QueryRow(fmt.Sprintf(`
		%s 
		WHERE task.id = ?`, defaultTaskSelect),
		taskID)

	// Nothing was updated
	if rowsAffected != 1 {
		if task, err = scanRow(row); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return Task{}, fmt.Errorf("%w for %d", ErrTaskNotExist, taskID)
			}
			return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
		}

		if task.UserID != userID {
			return Task{}, ErrNotOwner
		}

		return Task{}, ErrUnknown
	}

	if task, err = scanRow(row); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	// Not recurring, nothing more to do
	if !task.RecurringID.Valid || task.RecurringID.Int64 == 0 {
		if err = tx.Commit(); err != nil {
			return Task{}, fmt.Errorf("%w: %s", ErrTransactionCommit, err)
		}

		return task, nil
	}

	// TODO: either figure out how we want to manage this, or merge the behaviour together with the non-recurring ones
	// Just uncompleted a task. There is still a task in the database that was added when the recurring task was initially completed, but attempting to remove it may cause errors.
	// We may consider adding a "recurring_parent" field that references the taskID for recurring tasks, but until then, we don't want to accidentally remove the wrong task
	if !task.Done {
		if err = tx.Commit(); err != nil {
			return Task{}, fmt.Errorf("%w: %s", ErrTransactionCommit, err)
		}

		return task, nil
	}

	// If we completed a recurring task, add the next one
	nextTask = task

	if nextDueDate, err = parseRecurringPeriod(nextTask.RecurringPeriod.String); err != nil {
		return Task{}, err
	}

	nextTask.DueDate.Valid = true
	nextTask.DueDate.Int64 = nextDueDate.Unix()

	row = tx.QueryRow(`
    INSERT INTO task(title, category, description, due_date, recurring_id, user_id) 
    VALUES (?, ?, ?, ?, ?, ?)
    RETURNING task.id, title, category, description, due_date, completion_date, done, recurring_id, user_id, (
      SELECT period FROM recurring WHERE recurring.id = task.recurring_id)`,
		nextTask.Title, nextTask.Category, nextTask.Description, nextTask.DueDate, nextTask.RecurringID, nextTask.UserID)

	if nextTask, err = scanRow(row); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	if err = tx.Commit(); err != nil {
		return Task{}, fmt.Errorf("%w: %s", ErrTransactionCommit, err)
	}

	return task, nil
}
