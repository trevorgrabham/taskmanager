package services_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
	"time"
)

// UpdateDueDate changes the DueDate for the task identified by TaskID if UserID is the owner. Removes the DueDate if day.IsZero()
//
// Validation:
//   - ErrTaskValidation {Field: TaskFieldTaskID}
//     taskID negative or zero valued (< 1)
//     taskID doesn't belong to userID (repo: ErrNotOwner)
//     taskID doesn't exist
//   - ErrUserValidation {Field: UserFieldID}
//     userID negative or zero valued (< 1)
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - returns the updated task with all fields populated
func TestUpdateDueDate(t *testing.T) {
	cases := []struct {
		name string

		paramTaskID int
		paramUserID int
		paramDay    time.Time

		returnTask sqlite.Task
		returnErr  error

		wantTask services.Task
		checkErr func(*testing.T, error)
	}{
		// case: Empty TaskID
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MesageInvalidTaskID}
		{
			"Empty TaskID",

			0,
			3,
			time.Now(),

			sqlite.Task{
				ID:          31,
				UserID:      3,
				Title:       "Test Title",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Test description", Valid: true},
			},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

		// case: Negative TaskID
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MesageInvalidTaskID}
		{
			"Negatie TaskID",

			-32,
			3,
			time.Now(),

			sqlite.Task{
				ID:          32,
				UserID:      3,
				Title:       "Test Title",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Test description", Valid: true},
			},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

		// case: Empty UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MesageInvalidUserID}
		{
			"Empty UserID",

			32,
			0,
			time.Now(),

			sqlite.Task{
				ID:          32,
				UserID:      2,
				Title:       "Empty Title",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Empty description", Valid: true},
			},
			nil,

			services.Task{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Negative UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MesageInvalidUserID}
		{
			"Negative UserID",

			32,
			-3,
			time.Now(),

			sqlite.Task{
				ID:          32,
				UserID:      3,
				Title:       "Empty Title",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Empty description", Valid: true},
			},
			nil,

			services.Task{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Empty day
		// expected: Due date removed, nil error
		{
			"Empty Day",

			13,
			2,
			time.Time{},

			sqlite.Task{
				ID:          12,
				UserID:      2,
				Title:       "Title without a day",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Just to see what happens", Valid: true},
			},
			nil,

			services.Task{
				ID:          12,
				UserID:      2,
				Title:       "Title without a day",
				Category:    "go tests",
				Description: "Just to see what happens",
			},
			wantNoErr,
		},

		// case: Repo returned ErrNotOwner
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageNotOwner}
		{
			"UserID Not Owner",

			12,
			3,
			time.Now(),

			sqlite.Task{ },
			sqlite.ErrNotOwner,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageNotOwner),
		},

		// case: Repo returned empty
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageTaskNotFound}
		{
			"TaskID Not Exist",

			341,
			2,
			time.Now(),

			sqlite.Task{},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageTaskNotFound),
		},

		// case: Repo returned ErrInternalRepo
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			23,
			3,
			time.Now(),

			sqlite.Task{},
			sqlite.ErrInternalRepo,

			services.Task{},
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: Updated due date, nil error
		{
			"Happy Path",

			12,
			2,
			time.Time{},

			sqlite.Task{
				ID: 12,
				UserID: 2,
				Title: "Happy Title",
				Category: sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Happy description", Valid: true},
			},
			nil,

			services.Task{
				ID: 12,
				UserID: 2,
				Title: "Happy Title",
				Category:  "go tests",
				Description:  "Happy description",
			},
			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockTaskUpdateDueDate{
				stubRepo:   stubRepo{callCounts: make(map[string]int)},
				returnTask: tc.returnTask,
				returnErr:  tc.returnErr,
			}
			s := services.NewService(mock)

			gotTask, gotErr := s.UpdateDueDate(tc.paramTaskID, tc.paramUserID, tc.paramDay)

			tc.checkErr(t, gotErr)

			checkTask(t, gotTask, tc.wantTask)
		})
	}
}
