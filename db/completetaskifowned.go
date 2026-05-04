package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// CompleteTaskIfOwned completes the task identified by taskID if it is owned by userID.
//
// If the task identified by TaskID is already completed, it's treated as a no-op.
//
// If taskID is empty or no task matches, returns ErrTaskNotFound.
// If userID is empty, not the owner, or does not exist, returns ErrNotOwner.
// If a transient error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) CompleteTaskIfOwned(taskID, userID int) (err error) {
	var (
		row  *sql.Row
		task Task
	)
	row = r.db.QueryRow(`
		UPDATE task 
		SET completion_date = ?, done = 1 
		WHERE done = 0 AND completion_date IS NULL AND id = ? AND user_id = ?
		RETURNING id, title, category, description, completion_date + recurring_period AS due_date, completion_date, done, recurring_period, user_id`,
		time.Now().Unix(), taskID, userID)
	if task, err = scanRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if task, err = r.GetTaskByID(taskID, userID); err != nil {
				return err // GetTaskByID can return ErrNotOwner, so just forward the err
			}
			if task == (Task{}) {
				return ErrTaskNotFound
			}

			// Task was already completed
			return nil
		}
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	// Happy Path
	if _, err = r.AddTask(task); err != nil {
		return err
	}

	return nil
}
