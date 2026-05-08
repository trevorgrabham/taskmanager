package services_test

import (
	"database/sql"
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
)

// GetDashboardData aggergates tasks due this week, overdue tasks, and unscheduled tasks for the userID
//
// Validation:
//   - ErrUserValidation {Field: UserFieldID}
//     userID negative or zero valued (< 1)
//
// Errors:
//   - ErrInternalRepo
//     transient database error
//
// Happy Path:
//   - returns a Dashboard containing lists of tasks grouped by their dueDate
func TestGetDashboardData(t *testing.T) {
	cases := []struct {
		name string

		paramUserID int

		returnWeek           []sqlite.Task
		returnWeekErr        error
		returnOverdue        []sqlite.Task
		returnOverdueErr     error
		returnUnscheduled    []sqlite.Task
		returnUnscheduledErr error

		wantDash services.Dashboard
		checkErr func(*testing.T, error)
	}{

		// case: Empty UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Empty UserID",

			0,

			[]sqlite.Task{{
				ID:      12,
				UserID:  3,
				Title:   "Only thing this week",
				DueDate: sql.NullInt64{Int64: mustParseDate("2100-01-01").Unix(), Valid: true},
			}},
			nil,
			[]sqlite.Task{{
				ID:      10,
				UserID:  3,
				Title:   "Only one overdue",
				DueDate: sql.NullInt64{Int64: mustParseDate("2000-01-01").Unix(), Valid: true},
			}},
			nil,
			[]sqlite.Task{{
				ID:     11,
				UserID: 3,
				Title:  "Only unscheduled task",
			}},
			nil,

			services.Dashboard{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: Negative UserID
		// expected: ErrUserValidation {Field: UserFieldID, Message: MessageInvalidUserID}
		{
			"Negative UserID",

			-32,

			[]sqlite.Task{{
				ID:      12,
				UserID:  3,
				Title:   "Only thing this week",
				DueDate: sql.NullInt64{Int64: mustParseDate("2100-01-01").Unix(), Valid: true},
			}},
			nil,
			[]sqlite.Task{{
				ID:      10,
				UserID:  3,
				Title:   "Only one overdue",
				DueDate: sql.NullInt64{Int64: mustParseDate("2000-01-01").Unix(), Valid: true},
			}},
			nil,
			[]sqlite.Task{{
				ID:     11,
				UserID: 3,
				Title:  "Only unscheduled task",
			}},
			nil,

			services.Dashboard{},
			wantErrUserValidation(services.UserFieldID, services.MessageInvalidUserID),
		},

		// case: No tasks for the week
		// expected: empty list for the week, nil error
		{
			"Nothing For The Week",

			3,

			nil,
			nil,
			[]sqlite.Task{{
				ID:      10,
				UserID:  3,
				Title:   "Only one overdue",
				DueDate: sql.NullInt64{Int64: mustParseDate("2000-01-01").Unix(), Valid: true},
			}},
			nil,
			[]sqlite.Task{{
				ID:     11,
				UserID: 3,
				Title:  "Only unscheduled task",
			}},
			nil,

			services.Dashboard{
				nil,
				services.TaskList{{
					ID:      10,
					UserID:  3,
					Title:   "Only one overdue",
					DueDate: mustParseDate("2000-01-01"),
				}},
				services.TaskList{{
					ID:     11,
					UserID: 3,
					Title:  "Only unscheduled task",
				}},
			},
			wantNoErr,
		},

		// case: No overdue tasks
		// expected: empty list for overdue tasks, nil error
		{
			"Nothing Overdue",

			3,

			[]sqlite.Task{{
				ID:      12,
				UserID:  3,
				Title:   "Only thing this week",
				DueDate: sql.NullInt64{Int64: mustParseDate("2100-01-01").Unix(), Valid: true},
			}},
			nil,
			nil,
			nil,
			[]sqlite.Task{{
				ID:     11,
				UserID: 3,
				Title:  "Only unscheduled task",
			}},
			nil,

			services.Dashboard{
				services.TaskList{{
					ID:      12,
					UserID:  3,
					Title:   "Only thing this week",
					DueDate: mustParseDate("2100-01-01"),
				}},
				nil,
				services.TaskList{{
					ID:     11,
					UserID: 3,
					Title:  "Only unscheduled task",
				}},
			},
			wantNoErr,
		},

		// case: No unscheduled tasks
		// expected: empty list for unscheduled tasks, nil error
		{
			"Nothing Unscheduled",

			3,

			[]sqlite.Task{{
				ID:      12,
				UserID:  3,
				Title:   "Only thing this week",
				DueDate: sql.NullInt64{Int64: mustParseDate("2100-01-01").Unix(), Valid: true},
			}},
			nil,
			[]sqlite.Task{{
				ID:      10,
				UserID:  3,
				Title:   "Only one overdue",
				DueDate: sql.NullInt64{Int64: mustParseDate("2000-01-01").Unix(), Valid: true},
			}},
			nil,
			nil,
			nil,

			services.Dashboard{
				services.TaskList{{
					ID:      12,
					UserID:  3,
					Title:   "Only thing this week",
					DueDate: mustParseDate("2100-01-01"),
				}},
				services.TaskList{{
					ID:      10,
					UserID:  3,
					Title:   "Only one overdue",
					DueDate: mustParseDate("2000-01-01"),
				}},
				nil,
			},
			wantNoErr,
		},

		// case: Repo returned ErrInternalRepo from GetWeekOfTasks
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo Week",

			3,

			nil,
			sqlite.ErrInternalRepo,
			[]sqlite.Task{{
				ID:      10,
				UserID:  3,
				Title:   "Only one overdue",
				DueDate: sql.NullInt64{Int64: mustParseDate("2000-01-01").Unix(), Valid: true},
			}},
			nil,
			[]sqlite.Task{{
				ID:     11,
				UserID: 3,
				Title:  "Only unscheduled task",
			}},
			nil,

			services.Dashboard{},
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Repo returned ErrInternalRepo from GetOverdueTasks
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo Overdue",

			3,
			[]sqlite.Task{{
				ID:      12,
				UserID:  3,
				Title:   "Only thing this week",
				DueDate: sql.NullInt64{Int64: mustParseDate("2100-01-01").Unix(), Valid: true},
			}},
			nil,
			nil,
			sqlite.ErrInternalRepo,
			[]sqlite.Task{{
				ID:     11,
				UserID: 3,
				Title:  "Only unscheduled task",
			}},
			nil,

			services.Dashboard{},
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Repo returned ErrInternalRepo from GetUnscheduledTasks
		// expected: ErrInternalRepo (: err)
		{
			"Repo ErrInternalRepo Unscheduled",

			3,
			[]sqlite.Task{{
				ID:      12,
				UserID:  3,
				Title:   "Only thing this week",
				DueDate: sql.NullInt64{Int64: mustParseDate("2100-01-01").Unix(), Valid: true},
			}},
			nil,
			[]sqlite.Task{{
				ID:      10,
				UserID:  3,
				Title:   "Only one overdue",
				DueDate: sql.NullInt64{Int64: mustParseDate("2000-01-01").Unix(), Valid: true},
			}},
			nil,
			nil,
			sqlite.ErrInternalRepo,

			services.Dashboard{},
			wantErrIs(services.ErrInternalRepo),
		},

		// case: Happy path
		// expected: list of tasks, nil error
		{
			"Happy Path",

			3,

			[]sqlite.Task{{
				ID:      12,
				UserID:  3,
				Title:   "Only thing this week",
				DueDate: sql.NullInt64{Int64: mustParseDate("2100-01-01").Unix(), Valid: true},
			}},
			nil,
			[]sqlite.Task{{
				ID:      10,
				UserID:  3,
				Title:   "Only one overdue",
				DueDate: sql.NullInt64{Int64: mustParseDate("2000-01-01").Unix(), Valid: true},
			}},
			nil,
			[]sqlite.Task{{
				ID:     11,
				UserID: 3,
				Title:  "Only unscheduled task",
			}},
			nil,

			services.Dashboard{
				services.TaskList{{
					ID:      12,
					UserID:  3,
					Title:   "Only thing this week",
					DueDate: mustParseDate("2100-01-01"),
				}},
				services.TaskList{{
					ID:      10,
					UserID:  3,
					Title:   "Only one overdue",
					DueDate: mustParseDate("2000-01-01"),
				}},
				services.TaskList{{
					ID:     11,
					UserID: 3,
					Title:  "Only unscheduled task",
				}},
			},
			wantNoErr,
		}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := mockGetDashboardData{
				stubRepo:             stubRepo{callCounts: make(map[string]int)},
				returnWeek:           tc.returnWeek,
				returnWeekErr:        tc.returnWeekErr,
				returnOverdue:        tc.returnOverdue,
				returnOverdueErr:     tc.returnOverdueErr,
				returnUnscheduled:    tc.returnUnscheduled,
				returnUnscheduledErr: tc.returnUnscheduledErr,
			}
			s := services.NewService(mock)

			gotDash, gotErr := s.GetDashboardData(tc.paramUserID)

			tc.checkErr(t, gotErr)

			if len(gotDash.WeekOfTasks) != len(tc.wantDash.WeekOfTasks) {
				t.Errorf("WeekOfTasks wanted: %v\ngot: %v", tc.wantDash.WeekOfTasks, gotDash.WeekOfTasks)
			}
			for i := range gotDash.WeekOfTasks {
				checkTask(t, gotDash.WeekOfTasks[i], tc.wantDash.WeekOfTasks[i])
			}

			if len(gotDash.OverdueTasks) != len(tc.wantDash.OverdueTasks) {
				t.Errorf("OverdueTasks wanted: %v\ngot: %v", tc.wantDash.OverdueTasks, gotDash.OverdueTasks)
			}
			for i := range gotDash.OverdueTasks {
				checkTask(t, gotDash.OverdueTasks[i], tc.wantDash.OverdueTasks[i])
			}

			if len(gotDash.UnscheduledTasks) != len(tc.wantDash.UnscheduledTasks) {
				t.Errorf("UnscheduledTasks wanted: %v\ngot: %v", tc.wantDash.UnscheduledTasks, gotDash.UnscheduledTasks)
			}
			for i := range gotDash.UnscheduledTasks {
				checkTask(t, gotDash.UnscheduledTasks[i], tc.wantDash.UnscheduledTasks[i])
			}

		})
	}
}
