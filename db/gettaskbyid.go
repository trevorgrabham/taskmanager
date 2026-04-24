package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repo) GetTaskByID(taskID int) (t Task, err error) {
	var (
		caller = "GetTaskByID"
		row    *sql.Row
	)
	if !r.isConnected() {
		return Task{}, fmt.Errorf("%s: %w", caller, ErrNotConnected)
	}

	row = r.db.QueryRow(fmt.Sprintf(`
		%s
		WHERE task.id = ?`, defaultTaskSelect), 
		taskID)
	if t, err = scanRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return Task{}, err }
		return Task{}, fmt.Errorf("%s: %w", caller, NewErrRepo(err))
	}

	return t, nil
}
