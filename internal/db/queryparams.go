package db

import (
	"time"
)

type TaskQueryParams struct {
	WhichTasks WhichTasks
	From       time.Time
	To         time.Time
	Category   string
}

// ================================================== WhichTasks ==================================================

type WhichTasks int

const (
	AllTasks WhichTasks = iota
	IncTasks
	CompTasks
)
