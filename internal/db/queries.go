package db

import (
	"database/sql"
	"fmt"
	"time"

	"local/taskmanager2.0/internal/task"
)

type dbScannable interface {
	Scan(dest ...any) error
}

func scanTaskRow(r dbScannable) (task.Task, error) {
	var (
		id, dueDateUnix, completionDateUnix, done int64
		title                                     string
		category, description                     sql.NullString
	)

	if err := r.Scan(&id, &title, &category, &description, &dueDateUnix, &completionDateUnix, &done); err != nil {
		return task.Task{}, fmt.Errorf("scanning row: %s", err)
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

	return t, nil
}

func QueryTasks(db *sql.DB, params QueryParams) (task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("querying tasks: cannot query a nil database")
	}

	var rows *sql.Rows
	var err error
	if params.WhichTasks == task.All && params.Category == "" {
		rows, err = db.Query(fmt.Sprintf(`SELECT id, title, category, description, due_date, completion_date, done FROM %s WHERE due_date > ? AND due_date < ? ORDER BY category, due_date ASC`, taskTableName), params.From.Unix(), params.To.Unix())
		if err != nil {
			return nil, fmt.Errorf("querying tasks: %s", err)
		}
	} else if params.WhichTasks == task.All {
		rows, err = db.Query(fmt.Sprintf(`SELECT id, title, category, description, due_date, completion_date, done FROM %s WHERE category = ? AND due_date > ? AND due_date < ? ORDER BY category, due_date ASC`, taskTableName), params.Category, params.From.Unix(), params.To.Unix())
		if err != nil {
			return nil, fmt.Errorf("querying tasks: %s", err)
		}
	} else if params.Category == "" {
		rows, err = db.Query(fmt.Sprintf(`SELECT id, title, category, description, due_date, completion_date, done FROM %s WHERE done = ? AND due_date > ? AND due_date < ? ORDER BY category, due_date ASC`, taskTableName), params.WhichTasks-1, params.From.Unix(), params.To.Unix())
		if err != nil {
			return nil, fmt.Errorf("querying tasks: %s", err)
		}
	} else {
		rows, err = db.Query(fmt.Sprintf(`SELECT id, title, category, description, due_date, completion_date, done FROM %s WHERE category = ? AND done = ? AND due_date > ? AND due_date < ? ORDER BY category, due_date ASC`, taskTableName), params.Category, params.WhichTasks-1, params.From.Unix(), params.To.Unix())
		if err != nil {
			return nil, fmt.Errorf("querying tasks: %s", err)
		}
	}
	defer rows.Close()

	var results task.TaskList
	for rows.Next() {
		var t task.Task
		t, err = scanTaskRow(rows)
		if err != nil {
			return nil, fmt.Errorf("querying tasks: %s", err)
		}

		results = append(results, &t)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("processing task rows: %s", rows.Err())
	}

	return results, nil
}

func QueryAllTasks(db *sql.DB) (task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("querying all tasks: cannot query a nil database")
	}

	rows, err := db.Query(fmt.Sprintf(`SELECT id, title, category, description, due_date, completion_date, done FROM %s ORDER BY done DESC, category, completion_date DESC, due_date ASC`, taskTableName))
	if err != nil {
		return nil, fmt.Errorf("querying: %s", err)
	}
	defer rows.Close()

	var results task.TaskList
	for rows.Next() {
		var t task.Task
		t, err = scanTaskRow(rows)
		if err != nil {
			return nil, fmt.Errorf("querying all tasks: %s", err)
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

	rows, err := db.Query(fmt.Sprintf(`SELECT id, title, category, description, due_date, completion_date, done FROM %s WHERE done = 0 ORDER BY category, due_date ASC`, taskTableName))
	if err != nil {
		return nil, fmt.Errorf("querying incomplete tasks: %s", err)
	}
	defer rows.Close()

	var results task.TaskList
	for rows.Next() {
		var t task.Task
		t, err = scanTaskRow(rows)
		if err != nil {
			return nil, fmt.Errorf("querying incomplete tasks: %s", err)
		}

		results = append(results, &t)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("processing incomplete task rows: %s", rows.Err())
	}

	return results, nil
}

func QueryCompleteTasks(db *sql.DB) (task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("querying complete tasks: cannot query a nil database")
	}

	rows, err := db.Query(fmt.Sprintf(`SELECT id, title, category, description, due_date, completion_date, done FROM %s WHERE done = 1 ORDER BY category, due_date ASC`, taskTableName))
	if err != nil {
		return nil, fmt.Errorf("querying complete tasks: %s", err)
	}
	defer rows.Close()

	var results task.TaskList
	for rows.Next() {
		var t task.Task
		t, err = scanTaskRow(rows)
		if err != nil {
			return nil, fmt.Errorf("querying complete tasks: %s", err)
		}

		results = append(results, &t)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("processing complete task rows: %s", rows.Err())
	}

	return results, nil
}

func QueryTaskByID(db *sql.DB, id int) (task.Task, error) {
	if db == nil {
		return task.Task{}, fmt.Errorf("querying task by id: cannot query a nil database")
	}
	if id < 0 {
		return task.Task{}, fmt.Errorf("querying task by id: cannot query for id %d", id)
	}

	row := db.QueryRow(fmt.Sprintf(`SELECT id, title, category, description, due_date, completion_date, done FROM %s WHERE id = ?`, taskTableName), id)

	t, err := scanTaskRow(row)
	if err != nil {
		return task.Task{}, fmt.Errorf("querying task by id: %s", err)
	}
	if row.Err() != nil {
		return task.Task{}, fmt.Errorf("querying task by id: scanning row: %s", err)
	}

	return t, nil
}
