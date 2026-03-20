package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"local/taskmanager/internal/task"
)

type dbScannable interface {
	Scan(dest ...any) error
}

func scanTaskRow(r dbScannable) (task.Task, error) {
	var (
		id, done                        int64
		dueDateUnix, completionDateUnix sql.NullInt64
		title                           string
		category, description           sql.NullString
	)

	if err := r.Scan(&id, &title, &category, &description, &dueDateUnix, &completionDateUnix, &done); err != nil {
		return task.Task{}, fmt.Errorf("scanning row: %s", err)
	}

	t := task.Task{
		ID:    int(id),
		Title: title,
	}

	if category.Valid {
		t.Category = category.String
	}

	if description.Valid {
		t.Description = description.String
	}

	if dueDateUnix.Valid {
		t.DueDate = task.TaskDueDate(time.Unix(dueDateUnix.Int64, 0))
	}

	if completionDateUnix.Valid {
		t.CompletionDate = task.TaskDueDate(time.Unix(completionDateUnix.Int64, 0))
	}

	if done == 1 {
		t.Done = true
	}

	return t, nil
}

func QueryTasks(db *sql.DB, params TaskQueryParams) (task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("querying tasks: cannot query a nil database")
	}

	var (
		queryFilters []string
		queryArgs    []any
	)

	switch params.WhichTasks {
	case IncTasks:
		queryFilters = append(queryFilters, "done = ?")
		queryArgs = append(queryArgs, 0)
	case CompTasks:
		queryFilters = append(queryFilters, "done = ?")
		queryArgs = append(queryArgs, 1)
	}

	if !params.From.IsZero() {
		queryFilters = append(queryFilters, "due_date >= ?")
		queryArgs = append(queryArgs, params.From.Unix())
	}

	if !params.To.IsZero() {
		queryFilters = append(queryFilters, "due_date <= ?")
		queryArgs = append(queryArgs, params.To.Unix())
	}

	if params.Category != "" {
		queryFilters = append(queryFilters, "category = ?")
		queryArgs = append(queryArgs, params.Category)
	}

	var (
		rows *sql.Rows
		err  error
	)
	if len(queryFilters) > 0 {
		rows, err = db.Query(fmt.Sprintf(`SELECT id, title, category, description, due_date, completion_date, done FROM task WHERE %s ORDER BY category, done DESC, completion_date, due_date `, strings.Join(queryFilters, " AND ")), queryArgs...)
	} else {
		rows, err = db.Query(`SELECT id, title, category, description, due_date, completion_date, done FROM task`)
	}
	if err != nil {
		return nil, fmt.Errorf("querying tasks: %s", err)
	}
	defer rows.Close()

	var results task.TaskList
	for rows.Next() {
		var t task.Task
		t, err = scanTaskRow(rows)
		if err != nil {
			return nil, fmt.Errorf("querying tasks: %s", err)
		}

		results = append(results, t)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("querying tasks: %s", rows.Err())
	}

	return results, nil
}
