package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var dbFileName = "internal/db/backups/taskmanager_20260420_164207.db"

// var dbFileName = "/home/trevorgrabham/.config/taskmanager/tasks.db"

type dbConn interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type scannableRow interface {
	Scan(dest ...any) error
}

func addNewTask(db dbConn, queryColumns, queryPlaceholders []string, queryArgs []any) (id int, err error) {
	var res sql.Result
	res, err = db.Exec(fmt.Sprintf(`INSERT INTO task(%s) VALUES (%s)`, strings.Join(queryColumns, ", "), strings.Join(queryPlaceholders, ", ")), queryArgs...)
	if err != nil {
		return 0, fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	var insertID int64
	insertID, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return int(insertID), nil
}

func checkDateInitialized(day time.Time) (err error) {
	if day.IsZero() {
		return fmt.Errorf("%s: zero date provided", getCallingFunc(2))
	}

	return nil
}

func checkDBConnection(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("%s: cannot add into a nil database", getCallingFunc(2))
	}

	return nil
}

func checkQueryParams(params taskQueryParams) (err error) {
	if params.From.IsZero() && params.To.IsZero() && params.Category == "" && params.ID == 0 && !params.Unscheduled {
		return fmt.Errorf("%s: no query parameters specified", getCallingFunc(3))
	}

	return nil
}

func checkTaskNotEmpty(t task.Task) error {
	if t.IsZero() {
		return fmt.Errorf("%s: zero task provided", getCallingFunc(2))
	}

	return nil
}

func completeTask(db dbConn, taskID int) (completedTask task.Task, err error) {
	row := db.QueryRow(`
		UPDATE task 
		SET done = 1, completion_date = ? 
		WHERE id = ? 
		RETURNING id, title, category, description, due_date, completion_date, done, 
		(SELECT period FROM recurring WHERE recurring.id = task.recurring_id)`, time.Now().Unix(), taskID)
	completedTask, err = parseTask(row)
	if err != nil {
		return completedTask, fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return completedTask, nil
}

func computeNextDueDate(recurringPeriod string) (nextDueDate task.TaskDueDate, err error) {
	var (
		split          []string
		recurringUnit  string
		recurringValue int
		nextDate       time.Time
	)
	split = strings.Split(recurringPeriod, " ")
	recurringUnit = split[1]
	recurringValue, err = strconv.Atoi(split[0])
	if err != nil {
		return nextDueDate, fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	now := time.Now()
	switch recurringUnit {
	case "days":
		nextDate = now.AddDate(0, 0, recurringValue)
		nextDate = time.Date(nextDate.Year(), nextDate.Month(), nextDate.Day(), 0, 0, 0, 0, time.Local)
	case "weeks":
		nextDate = now.AddDate(0, 0, recurringValue*7)
		nextDate = time.Date(nextDate.Year(), nextDate.Month(), nextDate.Day(), 0, 0, 0, 0, time.Local)
	case "months":
		nextDate = now.AddDate(0, recurringValue, 0)
		nextDate = time.Date(nextDate.Year(), nextDate.Month(), nextDate.Day(), 0, 0, 0, 0, time.Local)
	default:
		return task.TaskDueDate(nextDate), fmt.Errorf("%s: unknown period format", getCallingFunc(2))
	}

	return task.TaskDueDate(nextDate), nil
}

func computeStartAndEndDays(startDate time.Time, numDays int) (start, end time.Time) {
	start = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.Local)
	end = time.Date(startDate.Year(), startDate.Month(), startDate.Day()+numDays-1, 23, 59, 0, 0, time.Local)
	return start, end
}

func deleteFutureRecurringTask(db dbConn, t task.Task) (err error) {
	var (
		columns, placeholders []string
		originalDueDate       task.TaskDueDate
		args, zipped          []any
	)
	originalDueDate = t.DueDate
	t.DueDate = task.TaskDueDate{} // so setupQueryParams() doesn't set up a param for due_date
	columns, placeholders, args, err = setupQueryParams(db, t)
	if err != nil {
		return fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	for i := range columns {
		placeholders[i] = "? = ?"
		zipped = append(zipped, columns[i], args[i])
	}
	zipped = append(zipped, originalDueDate)

	_, err = db.Exec(fmt.Sprintf(`DELETE FROM task WHERE %s AND due_date > ?`, strings.Join(placeholders, " AND ")), zipped...)
	if err != nil {
		return fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return nil
}

func deleteTask(db dbConn, taskID int) (err error) {
	_, err = db.Exec(`DELETE FROM task WHERE id = ?`, taskID)
	if err != nil {
		return fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return nil
}

func getCategories(db dbConn) (categories []string, err error) {
	var (
		rows     *sql.Rows
		category string
	)
	rows, err = db.Query(`SELECT DISTINCT category FROM task ORDER BY category`)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("%s: %s", getCallingFunc(2), err)
		}

		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return categories, nil
}

func getRecurringPeriod(db dbConn, taskID int) (recurringPeriod string, err error) {
	err = db.QueryRow(`SELECT period FROM recurring WHERE id = (SELECT recurring_id FROM task WHERE id = ?)`, taskID).Scan(&recurringPeriod)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return recurringPeriod, nil
}

func getRecurringPeriodID(db dbConn, recurringPeriod string) (recurringID int, err error) {
	err = db.QueryRow(`INSERT INTO recurring(period) VALUES (?) ON CONFLICT(period) DO UPDATE SET period=excluded.period RETURNING id`, recurringPeriod).Scan(&recurringID)
	if err != nil {
		return 0, fmt.Errorf("%s: %s", getCallingFunc(3), err)
	}

	return recurringID, nil
}

func getTaskByID(db dbConn, taskID int) (t task.Task, err error) {
	row := db.QueryRow(`
		SELECT task.id, title, category, description, due_date, completion_date, done, period 
		FROM task 
		LEFT JOIN recurring ON task.recurring_id = recurring.id
		WHERE task.id = ?`, taskID)
	if t, err = parseTask(row); err != nil {
		return t, fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return t, nil
}

func parseTask(r scannableRow) (task.Task, error) {
	var (
		id, done                               int64
		unixDueDate, unixCompletionDate        sql.NullInt64
		title                                  string
		category, description, recurringPeriod sql.NullString
		t                                      task.Task
	)

	if err := r.Scan(&id, &title, &category, &description, &unixDueDate, &unixCompletionDate, &done, &recurringPeriod); err != nil {
		return t, fmt.Errorf("parseTask: %s", err)
	}

	t.ID = int(id)
	t.Title = title

	if category.Valid {
		t.Category = category.String
	}

	if description.Valid {
		t.Description = description.String
	}

	if recurringPeriod.Valid {
		t.RecurringPeriod = recurringPeriod.String
	}

	if unixDueDate.Valid {
		t.DueDate = task.TaskDueDate(time.Unix(unixDueDate.Int64, 0))
	}

	if unixCompletionDate.Valid {
		t.CompletionDate = task.TaskDueDate(time.Unix(unixCompletionDate.Int64, 0))
	}

	if done == 1 {
		t.Done = true
	}

	return t, nil
}

func parseTasks(r *sql.Rows) (tasks task.TaskList, err error) {
	var t task.Task
	for r.Next() {
		if t, err = parseTask(r); err != nil {
			return nil, fmt.Errorf("%s: %s", getCallingFunc(3), err)
		}

		tasks = append(tasks, t)
	}
	if err = r.Err(); err != nil {
		return nil, fmt.Errorf("%s: %s", getCallingFunc(3), err)
	}

	return tasks, nil
}

func setupQueryParams(db dbConn, t task.Task) (columns, placeholders []string, args []any, err error) {
	columns = []string{"title"}
	placeholders = []string{"?"}
	args = []any{t.Title}

	if t.Category != "" {
		columns = append(columns, "category")
		placeholders = append(placeholders, "?")
		args = append(args, t.Category)
	}

	if t.Description != "" {
		columns = append(columns, "description")
		placeholders = append(placeholders, "?")
		args = append(args, t.Description)
	}

	columns = append(columns, "due_date")
	placeholders = append(placeholders, "?")
	if !t.DueDate.IsZero() {
		args = append(args, t.DueDate.Unix())
	} else {
		args = append(args, nil)
	}

	columns = append(columns, "recurring_id")
	placeholders = append(placeholders, "?")
	if t.RecurringPeriod != "" {
		var recurringID int
		if recurringID, err = getRecurringPeriodID(db, t.RecurringPeriod); err != nil {
			return nil, nil, nil, err
		}

		args = append(args, recurringID)
	} else {
		args = append(args, nil)
	}

	return columns, placeholders, args, nil
}

func sortWeekOfTasks(tasks task.TaskList) map[int]map[string]task.TaskList {
	var (
		tasksByDay = make(map[int]map[string]task.TaskList)
		dueDate    time.Time
	)

	for i := range 7 {
		tasksByDay[i] = make(map[string]task.TaskList)
	}
	for _, t := range tasks {
		dueDate = time.Time(t.DueDate)
		tasksByDay[int(dueDate.Weekday())][t.Category] = append(tasksByDay[int(dueDate.Weekday())][t.Category], t)
	}

	return tasksByDay
}

func startDBSession(db *sql.DB) (tx *sql.Tx, err error) {
	tx, err = db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return tx, nil
}

func toggleComplete(db dbConn, taskID int) (toggledTask task.Task, err error) {
	row := db.QueryRow(`
		UPDATE task 
		SET done = (done + 1) % 2, completion_date = 
			CASE 
				WHEN completion_date IS NULL THEN ? 
				ELSE NULL
			END
		WHERE id = ?
		RETURNING id, title, category, description, due_date, completion_date, done, (SELECT period FROM recurring WHERE id = task.recurring_id)
	`, time.Now().Unix(), taskID)
	toggledTask, err = parseTask(row)
	if err != nil {
		return toggledTask, fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return toggledTask, nil
}

func updateDueDate(db dbConn, taskID int, unixNewDueDate int64) (err error) {
	if unixNewDueDate == 0 {
		_, err = db.Exec(`UPDATE task SET due_date = NULL WHERE id = ?`, taskID)
	} else {
		_, err = db.Exec(`UPDATE task SET due_date = ? WHERE id = ?`, unixNewDueDate, taskID)
	}
	if err != nil {
		return fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return nil
}

func updateTask(db dbConn, t task.Task) error {
	columns, placeholders, args, err := setupQueryParams(db, t)
	if err != nil {
		return fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	for i := range columns {
		placeholders[i] = fmt.Sprintf("%s = ?", columns[i])
	}
	args = append(args, t.ID)

	_, err = db.Exec(fmt.Sprintf(`
		UPDATE task 
		SET %s
		WHERE id = ?`, strings.Join(placeholders, ", ")),
		args...)
	if err != nil {
		return fmt.Errorf("%s: %s", getCallingFunc(2), err)
	}

	return nil
}

func getCallingFunc(skip int) string {
	if skip < 0 {
		return ""
	}

	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}

	split := strings.Split(fn.Name(), ".")
	return split[len(split)-1]
}
