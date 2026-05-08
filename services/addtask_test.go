package services_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// AddTask inserts task if it is a valid task
//
// Validation:
//   - ErrTaskValidation {Field: TaskFieldTaskID}
//     taskID set ( != 0 )
//   - ErrTaskValidation {Field: TaskFieldTitle}
//     title empty
//   - ErrTaskValidation {Field: TaskFieldDone}
//     done is true (UpdateTask not meant to alter completed tasks)
//   - ErrTaskValidation {Field: TaskFieldCompletionDate}
//     completionDate set (UpdateTask not meant to alter completed tasks)
//   - ErrUserValidation {Field: UserFieldID}
//     userID negative or zero valued (< 1)
//     userID doesn't exist (repo: ErrConstraintFailure)
//
// Errors:
//
// Happy Path:
//   - returns inserted task after adding to the repo
func TestAddTask(t *testing.T) {
	cases := []struct {
		name string

		paramTask services.Task

		returnTask sqlite.Task
		returnErr  error

		wantTask services.Task
		checkErr func(*testing.T, error)
	}{

		// case: Has TaskID
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageInvalidTaskID}
		{
			"Has TaskID",

			services.Task{
				ID:              12,
				UserID:          4,
				Title:           "Has TaskID",
				Category:        "go tests",
				Description:     "This task already has an ID",
				DueDate:         mustParseDate("2000-10-10"),
				RecurringPeriod: 3600 * 24 * 5,
			},

			sqlite.Task{
				ID:              12,
				UserID:          4,
				Title:           "Has TaskID",
				Category:        sql.NullString{String: "go tests", Valid: true},
				Description:     sql.NullString{String: "This task already has an ID", Valid: true},
				DueDate:         sql.NullInt64{Int64: mustParseDate("2000-10-10").Unix(), Valid: true},
				RecurringPeriod: sql.NullInt64{Int64: 3600 * 24 * 5, Valid: true},
			},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

		// case: Has negative TaskID
		// expected: ErrTaskValidation {Field: TaskFieldTaskID, Message: MessageInvalidTaskID}
		{
			"Negative TaskID",

			services.Task{
				ID:              -21,
				UserID:          3,
				Title:           "Negative TaskID",
				Category:        "go tests",
				Description:     "This task has a negative ID",
				DueDate:         mustParseDate("2000-01-15"),
				RecurringPeriod: 3600 * 24 * 14,
			},

			sqlite.Task{
				ID:              -21,
				UserID:          3,
				Title:           "Negative TaskID",
				Category:        sql.NullString{String: "go tests", Valid: true},
				Description:     sql.NullString{String: "This task has a negative ID", Valid: true},
				DueDate:         sql.NullInt64{Int64: mustParseDate("2000-01-15").Unix(), Valid: true},
				RecurringPeriod: sql.NullInt64{Int64: 3600 * 24 * 14, Valid: true},
			},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTaskID, services.MessageInvalidTaskID),
		},

		// case: Empty UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Empty UserID",

			services.Task{
				UserID:      0,
				Title:       "Empty UserID",
				Category:    "go tests",
				Description: "This task has no owner",
				DueDate:     mustParseDate("2007-07-23"),
			},

			sqlite.Task{
				UserID:      0,
				Title:       "Empty UserID",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "This task has no owner", Valid: true},
				DueDate:     sql.NullInt64{Int64: mustParseDate("2007-07-23").Unix(), Valid: true},
			},
			nil,

			services.Task{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Has negative UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Negative UserID",

			services.Task{
				UserID:      -54,
				Title:       "Negative UserID",
				Category:    "go tests",
				Description: "This task has an imaginary owner",
				DueDate:     mustParseDate("2005-02-14"),
			},

			sqlite.Task{
				UserID:      -54,
				Title:       "Negative UserID",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "This task has an imaginary owner", Valid: true},
				DueDate:     sql.NullInt64{Int64: mustParseDate("2005-02-14").Unix(), Valid: true},
			},
			nil,

			services.Task{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Empty title
		// expected: ErrTaskValidation {Field: TaskFieldTitle, Message: MessageEmptyTitle}
		{
			"Empty Title",

			services.Task{
				UserID:          4,
				Title:           "",
				Category:        "go tests",
				Description:     "This task has no title",
				RecurringPeriod: 3600 * 24,
			},

			sqlite.Task{
				UserID:          4,
				Title:           "",
				Category:        sql.NullString{String: "go tests", Valid: true},
				Description:     sql.NullString{String: "This task has no title", Valid: true},
				RecurringPeriod: sql.NullInt64{Int64: 3600 * 24, Valid: true},
			},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldTitle, services.MessageEmptyTitle),
		},

		// case: Empty category
		// expected: added task, nil error
		{
			"Empty Category",

			services.Task{
				UserID:          2,
				Title:           "No category for this one",
				Description:     "Still gets a description though",
				DueDate:         mustParseDate("2026-01-01"),
				RecurringPeriod: 3600 * 24 * 30,
			},

			sqlite.Task{
				UserID:          2,
				Title:           "No category for this one",
				Description:     sql.NullString{String: "Still gets a description though", Valid: true},
				DueDate:         sql.NullInt64{Int64: mustParseDate("2026-01-01").Unix(), Valid: true},
				RecurringPeriod: sql.NullInt64{Int64: 3600 * 24 * 30, Valid: true},
			},
			nil,

			services.Task{
				UserID:          2,
				Title:           "No category for this one",
				Description:     "Still gets a description though",
				DueDate:         mustParseDate("2026-01-01"),
				RecurringPeriod: 3600 * 24 * 30,
			},
			wantNoErr,
		},

		// case: With due date and time
		// expected: added task, nil error
		{
			"With DueDate Time",

			services.Task{
				UserID:      3,
				Title:       "With a duedate time",
				Category:    "go tests",
				Description: "Making sure that times are stored properly",
				DueDate:     mustParseDateTime("2026-01-01 01:23"),
			},

			sqlite.Task{
				UserID:      3,
				Title:       "With a duedate time",
				Category:    sql.NullString{String: "go tests", Valid: true},
				Description: sql.NullString{String: "Making sure that times are stored properly", Valid: true},
				DueDate:     sql.NullInt64{Int64: mustParseDateTime("2026-01-01 01:23").Unix(), Valid: true},
			},
			nil,

			services.Task{
				UserID:      3,
				Title:       "With a duedate time",
				Category:    "go tests",
				Description: "Making sure that times are stored properly",
				DueDate:     mustParseDateTime("2026-01-01 01:23"),
			},
			wantNoErr,
		},

		// case: CompletionDate set
		// expected: ErrTaskValidation {Field: TaskFieldCompletionDate, Message: MessageTaskAlreadyCompleted}
		{
			"CompletionDate Set",

			services.Task{
				UserID:         1,
				Title:          "With a completion date",
				Category:       "go tests",
				Description:    "This should generate an error",
				DueDate:        mustParseDate("2006-06-06"),
				CompletionDate: mustParseDate("1997-01-01"),
			},

			sqlite.Task{
				UserID:         1,
				Title:          "With a completion date",
				Category:       sql.NullString{String: "go tests", Valid: true},
				Description:    sql.NullString{String: "Just to make sure it gets ignored", Valid: true},
				DueDate:        sql.NullInt64{Int64: mustParseDate("2006-06-06").Unix(), Valid: true},
				CompletionDate: sql.NullInt64{Int64: mustParseDate("1997-01-01").Unix(), Valid: true},
			},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldCompletionDate, services.MessageTaskAlreadyCompleted),
		},

		// case: Done set
		// expected: ErrTaskValidation {Field: TaskFieldDone, Message: MessageTaskAlreadyCompleted}
		{
			"Done Set",

			services.Task{
				UserID:          5,
				Title:           "With Done set",
				Category:        "go tests",
				Description:     "This shouldn't go well",
				DueDate:         mustParseDate("3000-01-01"),
				Done:            true,
				RecurringPeriod: 3600 * 24 * 3,
			},

			sqlite.Task{
				UserID:          5,
				Title:           "With Done set",
				Category:        sql.NullString{String: "go tests", Valid: true},
				Description:     sql.NullString{String: "This shouldn't go well", Valid: true},
				DueDate:         sql.NullInt64{Int64: mustParseDate("3000-01-01").Unix(), Valid: true},
				Done:            true,
				RecurringPeriod: sql.NullInt64{Int64: 3600 * 24 * 3, Valid: true},
			},
			nil,

			services.Task{},
			wantErrTaskValidation(services.TaskFieldDone, services.MessageTaskAlreadyCompleted),
		},

		// case: Repo returned an ErrConstraintFailure (UserID failed FK constraint)
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Repo ErrConstraintFailure",

			services.Task{
				UserID:          7,
				Title:           "This shouldn't matter",
				Category:        "go tests",
				Description:     "If I see any of this data, something went wrong",
				DueDate:         mustParseDate("1994-04-02"),
				RecurringPeriod: 3600 * 24 * 6,
			},

			sqlite.Task{},
			sqlite.ErrConstraintFailure,

			services.Task{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Repo returned an ErrInternalRepo
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo",

			services.Task{
				UserID: 1,
				Title: "Internal Repo Error",
				Category: "go tests",
				Description: "Something went wrong in the repo",
				DueDate: mustParseDate("2012-02-02"),
			},

			sqlite.Task{
			},
			sqlite.ErrInternalRepo,

			services.Task{},
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: added task, nil error
		{
			"Happy Path",

			services.Task{
				UserID: 3,
				Title: "Happy Path",
				Category: "go tests",
				Description: "Just glad to finally get here",
				RecurringPeriod: 3600 * 24,
			},

			sqlite.Task{
				UserID: 3,
				Title: "Happy Path",
				Category: sql.NullString{ String: "go tests", Valid: true}, 
				Description: sql.NullString{ String: "Just glad to finally get here", Valid: true}, 
				RecurringPeriod: sql.NullInt64{ Int64: 3600 * 24, Valid: true},
			},
			nil,

			services.Task{
				UserID: 3,
				Title: "Happy Path",
				Category: "go tests",
				Description: "Just glad to finally get here",
				RecurringPeriod: 3600 * 24,
			},
			wantNoErr,
		}	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockAddTask{
				stubRepo:   stubRepo{callCounts: make(map[string]int)},
				returnTask: tc.returnTask,
				returnErr:  tc.returnErr,
			}
			s := services.NewService(mock)

			gotTask, gotErr := s.AddTask(tc.paramTask)

			tc.checkErr(t, gotErr)

			checkTask(t, gotTask, tc.wantTask)
		})
	}
}
