package services_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
	"time"
)

// UpdateTask updates Task if UserID is its owner.
//
// Validation:
//   - ErrTaskValidation {Field: TaskFieldTaskID}
//     taskID negative or zero valued (< 1)
//     taskID doesn't belong to userID (repo: ErrNotOwner)
//     taskID doesn't exist
//   - ErrTaskValidation {Field: TaskFieldTitle}
//     title empty
//   - ErrTaskValidation {Field: TaskFieldDone}
//     done is true (UpdateTask not meant to alter completed tasks)
//   - ErrTaskValidation {Field: TaskFieldCompletionDate}
//     completionDate set (UpdateTask not meant to alter completed tasks)
//   - ErrUserValidation {Field: UserFieldID}
//     userID negative or zero valued (< 1)
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - returns the updated task with all fields populated
func TestUpdateTask(t *testing.T) {
	cases := []struct {
		name string

		paramTask services.Task

		returnTask sqlite.Task
		returnErr  error

		wantTask services.Task
		checkErr func(*testing.T, error)
	}{

		// case: Empty TaskID
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageInvalidTaskID}
		{
			"Empty TaskID",

			services.Task{ID: 0, UserID: 3, Title: "Empty TaskID", Category: "go tests"},

			sqlite.Task{ID: 0, UserID: 3, Title: "Empty TaskID", Category: sql.NullString{String: "go tests", Valid: true}},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

		// case: Negative TaskID
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageInvalidTaskID}
		{
			"Negative TaskID",

			services.Task{ID: -35, UserID: 3, Title: "Negative TaskID", Category: "go tests"},

			sqlite.Task{ID: -35, UserID: 3, Title: "Negative TaskID", Category: sql.NullString{String: "go tests", Valid: true}},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

		// case: Empty UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Empty UserID",

			services.Task{ID: 1, UserID: 0, Title: "Empty UserID", Category: "go tests"},

			sqlite.Task{ID: 1, UserID: 0, Title: "Empty UserID", Category: sql.NullString{String: "go tests", Valid: true}},
			nil,

			services.Task{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Negative UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Negative UserID",

			services.Task{ID: 1, UserID: -87, Title: "Negative UserID", Category: "go tests"},

			sqlite.Task{ID: 1, UserID: -87, Title: "Negative UserID", Category: sql.NullString{String: "go tests", Valid: true}},
			nil,

			services.Task{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Empty title
		// expected: ErrTaskValidation { Field: TaskFieldTitle, Message: MessageEmptyTitle}
		{
			"Empty Title",

			services.Task{ID: 3, UserID: 1, Category: "go tests"},

			sqlite.Task{ID: 3, UserID: 1, Category: sql.NullString{String: "go tests", Valid: true}},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTitle, services.MessageEmptyTitle),
		},

		// case: Non-empty completion date
		// expected: ErrTaskValidation { Field: TaskFieldCompletionDate, Message: MessageTaskAlreadyCompleted}
		{
			"Completion Date Set",

			services.Task{ID: 12, UserID: 3, Title: "Completion Date Set", Category: "go tests", CompletionDate: time.Now()},

			sqlite.Task{ID: 12, UserID: 3, Title: "Completion Date Set", Category: sql.NullString{String: "go tests", Valid: true}, CompletionDate: sql.NullInt64{Int64: time.Now().Unix(), Valid: true}},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldCompletionDate, services.MessageTaskAlreadyCompleted),
		},

		// case: Done is true
		// expected: ErrTaskValidation { Field: TaskFieldDone, Message: MessageTaskAlreadyCompleted}
		{
			"Done Set",

			services.Task{ID: 13, UserID: 2, Title: "Done Set", Category: "go tests", Done: true},

			sqlite.Task{ID: 13, UserID: 2, Title: "Done Set", Category: sql.NullString{String: "go tests", Valid: true}, Done: true},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldDone, services.MessageTaskAlreadyCompleted),
		},

		// case: Repo returned empty Task, nil error
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageTaskNotFound}
		{
			"Repo Empty Task",

			services.Task{ID: 123, UserID: 1, Title: "This Task doesn't exist", Category: "go tests"},

			sqlite.Task{},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageTaskNotFound),
		},

		// case: Repo returned ErrNotOwner
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageNotOwner}
		{
			"Repo ErrNotOwner",

			services.Task{ID: 1, UserID: 6, Title: "This user doesn't own me", Category: "go tests"},

			sqlite.Task{},
			sqlite.ErrNotOwner,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageNotOwner),
		},

		// case: Repo returned ErrInternalRepo
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			services.Task{ID: 12, UserID: 5, Title: "Internal Repo Error", Category: "go tests"},

			sqlite.Task{},
			sqlite.ErrInternalRepo,

			services.Task{},
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: Task updated, nil error
		{
			"Happy Path",

			services.Task{ID: 7, UserID: 2, Title: "Happy Path", Category: "happy", Description: "Such a happy path", DueDate: mustParseDate("2020-02-02"), RecurringPeriod: 3600 * 24 * 7},

			sqlite.Task{ID: 7, UserID: 2, Title: "Happy Path", Category: sql.NullString{String: "happy", Valid: true}, Description: sql.NullString{String: "Such a happy path", Valid: true}, DueDate: sql.NullInt64{Int64: mustParseDate("2020-02-02").Unix(), Valid: true}, RecurringPeriod: sql.NullInt64{Int64: 3600 * 24 * 7, Valid: true}},
			nil,

			services.Task{ID: 7, UserID: 2, Title: "Happy Path", Category: "happy", Description: "Such a happy path", DueDate: mustParseDate("2020-02-02"), RecurringPeriod: 3600 * 24 * 7},
			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockUpdateTask{
				stubRepo:   stubRepo{callCounts: make(map[string]int)},
				returnTask: tc.returnTask,
				returnErr:  tc.returnErr,
			}
			s := services.NewService(mock)

			gotTask, gotErr := s.UpdateTask(tc.paramTask)

			tc.checkErr(t, gotErr)

			checkTask(t, gotTask, tc.wantTask)
		})
	}
}
