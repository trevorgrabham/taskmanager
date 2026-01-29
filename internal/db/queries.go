package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"local/taskmanager2.0/internal/task"
)

func QueryAll(db *sql.DB) (task.TaskList, error) {
	if db == nil {
		return nil, errors.New("querying: cannot query a nil database")
	}
	rows, err := db.Query(`SELECT id, title, category, description, due_date, completion_date, done FROM tasks`)
	if err != nil {
		return nil, fmt.Errorf("querying: %s", err)
	}
	defer rows.Close()

	var (
		id, dueDateUnix, completionDateUnix, done int64
		title                                     string
		category, description                     sql.NullString
		results                                   task.TaskList
	)

	for rows.Next() {
		if err = rows.Scan(&id, &title, &category, &description, &dueDateUnix, &completionDateUnix, &done); err != nil {
			return nil, fmt.Errorf("scanning row: %s", err)
		}

		dueDate := task.TaskDueDate(time.Unix(dueDateUnix, 0))
		t := task.Task{
			Title:   title,
			DueDate: &dueDate,
		}

		if category.Valid {
			t.Category = category.String
		}

		if description.Valid {
			t.Description = description.String
		}

		if completionDateUnix != -1 {
			competionDate := task.TaskDueDate(time.Unix(completionDateUnix, 0))
			t.CompletionDate = &competionDate
		}

		if !time.Now().After(time.Time(dueDate)) {
			t.Type = task.Upcoming
		} else {
			t.Type = task.Due
		}

		if done == 1 {
			t.Done = true
			t.Type = task.Done
		}

		results = append(results, &t)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return results, nil
}
