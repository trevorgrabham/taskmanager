package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TaskToggleCompleteIfOwned toggles the completion status for the task identified by taskID if it is owned by userID.
//
// If taskID is empty or doesn't exist, returns ErrTaskNotFound.
// If userID is empty or not the owner, returns ErrNotOwner.
// If an error occurs in the Repo, returns an ErrInternalRepo.
func (r *Repo) TaskToggleCompleteIfOwned(taskID, userID int) (err error) {
	var (
		row  *sql.Row
		task Task
	)
	row = r.db.QueryRow(`
		UPDATE task 
		SET done = (done + 1) % 2, completion_date = 
			CASE 
				WHEN completion_date IS NULL THEN ?
				ELSE NULL
			END
		WHERE id = ? AND user_id = ?
		RETURNING id, title, category, description, completion_date + recurring_period AS due_date, completion_date, done, recurring_period, user_id`,
		time.Now().Unix(), taskID, userID)

	if task, err = scanRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if task, err = r.GetTaskByID(taskID, userID); err != nil {
				// GetTaskByID handlers ErrNotOwner, just forward
				return err
			}
			if task == (Task{}) {
				return ErrTaskNotFound
			}

			// Shouldn't ever get here
			return ErrInternalRepo
		}

		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	// Happy Path
	if !task.Done {
		return nil
	}
	if _, err = r.AddTask(task); err != nil {
		return fmt.Errorf("%w: %s", ErrInternalRepo, err)
	}

	return nil
}
