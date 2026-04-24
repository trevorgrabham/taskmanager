package services

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"time"
)

type Task struct {
	ID              int
	UserID          int
	RecurringID     int
	Title           string
	Category        string
	Description     string
	DueDate         time.Time
	CompletionDate  time.Time
	Done            bool
	RecurringPeriod string // n days || n weeks || n months
}
type TaskList []Task

func (t Task) IsZero() bool {
	return t.ID == 0 && t.Title == "" && t.Category == "" && t.Description == "" && t.DueDate.IsZero() && t.CompletionDate.IsZero() && !t.Done && t.RecurringPeriod == "" && t.UserID == 0
}

func parseRepoTaskToTask(repoTask sqlite.Task) (t Task) {
	t.ID = repoTask.ID
	t.UserID = repoTask.UserID
	t.Title = repoTask.Title
	t.Done = repoTask.Done
	if repoTask.Category.Valid {
		t.Category = repoTask.Category.String
	}
	if repoTask.Description.Valid {
		t.Description = repoTask.Description.String
	}
	if repoTask.DueDate.Valid {
		t.DueDate = time.Unix(repoTask.DueDate.Int64, 0)
	}
	if repoTask.CompletionDate.Valid {
		t.CompletionDate = time.Unix(repoTask.CompletionDate.Int64, 0)
	}
	if repoTask.RecurringID.Valid {
		t.RecurringID = int(repoTask.RecurringID.Int64)
	}
	if repoTask.RecurringPeriod.Valid {
		t.RecurringPeriod = repoTask.RecurringPeriod.String
	}

	return t
}

func parseRepoTaskListToTaskList(repoTasks []sqlite.Task) (tasks TaskList) {
	for _, rt := range repoTasks {
		tasks = append(tasks, parseRepoTaskToTask(rt))
	}

	return tasks
}

func parseTaskToRepoTask(t Task) (repoTask sqlite.Task) {
	repoTask.ID = t.ID
	repoTask.UserID = t.UserID
	repoTask.Title = t.Title
	repoTask.Done = t.Done
	repoTask.Category = toNullString(t.Category)
	repoTask.Description = toNullString(t.Description)
	repoTask.DueDate = unixDateToNullInt64(t.DueDate.Unix())
	repoTask.CompletionDate = unixDateToNullInt64(t.CompletionDate.Unix())
	repoTask.RecurringID = toNullInt64(int64(t.RecurringID))
	repoTask.RecurringPeriod = toNullString(t.RecurringPeriod)

	return repoTask
}

func toNullString(s string) sql.NullString {
	return sql.NullString{Valid: s != "", String: s}
}

func toNullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{Valid: n != 0, Int64: n}
}

func unixDateToNullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{Valid: n > 0, Int64: n}
}

