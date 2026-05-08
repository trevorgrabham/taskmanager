package services

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"time"
)

type Task struct {
	ID              int
	UserID          int
	Title           string
	Category        string
	Description     string
	DueDate         time.Time
	CompletionDate  time.Time
	Done            bool
	RecurringPeriod int
}
type TaskList []Task

func (t Task) IsZero() bool {
	return t.ID == 0 && t.Title == "" && t.Category == "" && t.Description == "" && t.DueDate.IsZero() && t.CompletionDate.IsZero() && !t.Done && t.RecurringPeriod == 0 && t.UserID == 0
}

// parseRepoTaskToTask maps a repo Task to a Task object.
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
	if repoTask.RecurringPeriod.Valid {
		t.RecurringPeriod = int(repoTask.RecurringPeriod.Int64)
	}

	return t
}

func parseRepoTaskListToTaskList(repoTasks []sqlite.Task) (tasks TaskList) {
	for _, rt := range repoTasks {
		tasks = append(tasks, parseRepoTaskToTask(rt))
	}

	return tasks
}

// parseTaskToRepoTask maps a Task object to a Repo Task.
func parseTaskToRepoTask(t Task) (repoTask sqlite.Task) {
	repoTask.ID = t.ID
	repoTask.UserID = t.UserID
	repoTask.Title = t.Title
	repoTask.Done = t.Done
	repoTask.Category = toNullString(t.Category)
	repoTask.Description = toNullString(t.Description)
	repoTask.DueDate = unixDateToNullInt64(t.DueDate.Unix())
	repoTask.CompletionDate = unixDateToNullInt64(t.CompletionDate.Unix())
	repoTask.RecurringPeriod = unixDateToNullInt64(int64(t.RecurringPeriod))

	return repoTask
}

// toNullString populates a sql.NullString using a string. It sets the Valid field to true if the string is not empty.
func toNullString(s string) sql.NullString {
	return sql.NullString{Valid: s != "", String: s}
}

// unixDateToNullInt64 populates a sql.NullInt64 using an int64 representation of a time.Time obejct. It sets the Valid field to true if the int64 > 0.
//
// For times on or before Jan 1, 1970, this will set them to non-valid sql.NullInt64 objects. They must be parsed by hand.
func unixDateToNullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{Valid: n > 0, Int64: n}
}
