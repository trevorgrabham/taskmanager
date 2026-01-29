package db

import (
	"database/sql"
	"fmt"
	"time"

	"local/taskmanager2.0/internal/task"
)

func QueryAllTasks(db *sql.DB) (task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("querying all tasks: cannot query a nil database")
	}
	rows, err := db.Query(`SELECT id, title, category, description, due_date, completion_date, done FROM tasks ORDER BY done DESC, category, completion_date DESC, due_date ASC`)
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
			ID:      int(id),
			Title:   title,
			DueDate: &dueDate,
		}

		if category.Valid {
			t.Category = category.String
		}

		if description.Valid {
			t.Description = description.String
		}

		if completionDateUnix > -1 {
			competionDate := task.TaskDueDate(time.Unix(completionDateUnix, 0))
			t.CompletionDate = &competionDate
		}

		if time.Now().Before(time.Time(dueDate)) {
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
		return nil, fmt.Errorf("processing task rows: %s", rows.Err())
	}

	return results, nil
}

func QueryIncompleteTasks(db *sql.DB) (task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("querying incomplete tasks: cannot query a nil database")
	}

	rows, err := db.Query(`SELECT id, title, category, description, due_date, completion_date, done FROM tasks WHERE done = 0 ORDER BY category, due_date ASC`)
	if err != nil {
		return nil, fmt.Errorf("querying incomplete tasks: %s", err)
	}
	defer rows.Close()

	var (
		id, dueDateUnix, completionDateUnix, done int64
		title                                     string
		category, description                     sql.NullString
		results                                   task.TaskList
	)

	for rows.Next() {
		err = rows.Scan(&id, &title, &category, &description, &dueDateUnix, &completionDateUnix, &done)
		if err != nil {
			return nil, fmt.Errorf("scanning incomplete tasks: %s", err)
		}

		dueDate := task.TaskDueDate(time.Unix(dueDateUnix, 0))
		t := task.Task{ID: int(id), Title: title, DueDate: &dueDate}

		if category.Valid {
			t.Category = category.String
		}

		if description.Valid {
			t.Description = description.String
		}

		if completionDateUnix > -1 {
			completionDate := task.TaskDueDate(time.Unix(completionDateUnix, 0))
			t.CompletionDate = &completionDate
		}

		if time.Now().Before(time.Time(dueDate)) {
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
		return nil, fmt.Errorf("processing incomplete task rows: %s", rows.Err())
	}

	return results, nil
}
