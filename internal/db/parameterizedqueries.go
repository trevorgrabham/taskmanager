package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"strings"
	"time"
)

type taskQueryParams struct {
	WhichTasks  completionStatus
	From        time.Time
	To          time.Time
	Unscheduled bool
	Category    string
	ID          int
}

type completionStatus int

const (
	dontCare completionStatus = iota
	incomplete
	complete
)

func queryTasks(db *sql.DB, params taskQueryParams) (task.TaskList, error) {
	var (
		err        error
		conditions []string
		args       []any
		rows       *sql.Rows
	)
	if err = checkDBConnection(db); err != nil {
		return nil, err
	}

	if err = checkQueryParams(params); err != nil {
		return nil, err
	}

	conditions, args = prepareQueryParams(params)

	if len(conditions) > 0 {
		rows, err = db.Query(fmt.Sprintf(`
			SELECT task.id, title, category, description, due_date, completion_date, done, period 
			FROM task 
			LEFT JOIN recurring ON recurring.id = task.recurring_id
			WHERE %s 
			ORDER BY category, done DESC, completion_date, title DESC`, strings.Join(conditions, " AND ")), args...)
	} else {
		rows, err = db.Query(`
			SELECT task.id, title, category, description, due_date, completion_date, done, period 
			FROM task 
			LEFT JOIN recurring ON recurring.id = task.recurring_id
			ORDER BY category, done DESC, completion_date, title DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("queryTasks: %s", err)
	}
	defer rows.Close()

	return parseTasks(rows)
}

func prepareQueryParams(params taskQueryParams) (conditions []string, args []any) {
	if params.ID > 0 {
		conditions = append(conditions, "task.id = ?")
		args = append(args, params.ID)
	}

	if params.Category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, params.Category)
	}

	switch params.WhichTasks {
	case incomplete:
		conditions = append(conditions, "done = ?")
		args = append(args, 0)
	case complete:
		conditions = append(conditions, "done = ?")
		args = append(args, 1)
	}

	if params.Unscheduled {
		conditions = append(conditions, "due_date IS ?")
		args = append(args, nil)
		return conditions, args
	}

	if !params.From.IsZero() {
		switch params.WhichTasks {
		case incomplete:
			conditions = append(conditions, "due_date >= ?")
			args = append(args, params.From.Unix())
		case complete:
			conditions = append(conditions, "completion_date >= ?")
			args = append(args, params.From.Unix())
		case dontCare:
			conditions = append(conditions, "((done = 0 AND due_date >= ?) OR (done = 1 AND completion_date >= ?))")
			args = append(args, params.From.Unix(), params.From.Unix())
		}
	}

	if !params.To.IsZero() {
		switch params.WhichTasks {
		case incomplete:
			conditions = append(conditions, "due_date <= ?")
			args = append(args, params.To.Unix())
		case complete:
			conditions = append(conditions, "completion_date <= ?")
			args = append(args, params.To.Unix())
		case dontCare:
			conditions = append(conditions, "((done = 0 AND due_date <= ?) OR (done = 1 AND completion_date <= ?))")
			args = append(args, params.To.Unix(), params.To.Unix())
		}
	}

	return conditions, args
}
