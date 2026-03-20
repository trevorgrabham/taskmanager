package db

import (
	"time"
)

/*
func TaskParams(opts ...TaskQueryParamsFunc) TaskQueryParams {
	t := TaskQueryParams{
		WhichTasks: AllTasks,
		From:       time.Time{},
		To:         time.Date(3000, 1, 1, 0, 0, 0, 0, time.Now().Location()),
		Category:   "",
	}
	for _, fn := range opts {
		fn(&t)
	}
	return t
}
*/

type TaskQueryParams struct {
	WhichTasks WhichTasks
	From       time.Time
	To         time.Time
	Category   string
}

/*
type TaskQueryParamsFunc func(t *TaskQueryParams)

func All() TaskQueryParamsFunc {
	return func(t *TaskQueryParams) {
		t.WhichTasks = AllTasks
	}
}

func Inc() TaskQueryParamsFunc {
	return func(t *TaskQueryParams) {
		t.WhichTasks = IncTasks
	}
}

func Comp() TaskQueryParamsFunc {
	return func(t *TaskQueryParams) {
		t.WhichTasks = CompTasks
	}
}

func From(from time.Time) TaskQueryParamsFunc {
	return func(t *TaskQueryParams) {
		t.From = from
	}
}

func To(to time.Time) TaskQueryParamsFunc {
	return func(t *TaskQueryParams) {
		t.From = to
	}
}

func Category(category string) TaskQueryParamsFunc {
	return func(t *TaskQueryParams) {
		t.Category = category
	}
}
*/

// ================================================== WhichTasks ==================================================

type WhichTasks int

const (
	AllTasks WhichTasks = iota
	IncTasks
	CompTasks
)
