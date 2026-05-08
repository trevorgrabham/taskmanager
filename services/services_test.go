package services_test

import (
	"errors"
	sqlite "local/taskmanager/db"
	"local/taskmanager/services"
	"testing"
	"time"
)

type stubRepo struct {
	callCounts map[string]int
}

func (r stubRepo) GetCount(method string) int { return r.callCounts[method] }

func (r stubRepo) AddTask(sqlite.Task) (sqlite.Task, error) {
	r.callCounts["AddTask"]++
	return sqlite.Task{}, nil
}
func (r stubRepo) CompleteTaskIfOwned(taskID, userID int) error {
	r.callCounts["CompleteTaskIfOwned"]++
	return nil
}
func (r stubRepo) DeleteTaskIfOwned(taskID, userID int) error {
	r.callCounts["DeleteTaskIfOwned"]++
	return nil
}
func (r stubRepo) DeleteSession(sessionID string) error { r.callCounts["DeleteSession"]++; return nil }
func (r stubRepo) GetOverdueTasks(userID int) ([]sqlite.Task, error) {
	r.callCounts["GetOverdueTasks"]++
	return nil, nil
}
func (r stubRepo) GetTaskByID(taskID, userID int) (sqlite.Task, error) {
	r.callCounts["GetTaskByID"]++
	return sqlite.Task{}, nil
}
func (r stubRepo) GetTasksForDay(userID int, day time.Time) ([]sqlite.Task, error) {
	r.callCounts["GetTasksForDay"]++
	return nil, nil
}
func (r stubRepo) GetUnscheduledTasks(userID int) ([]sqlite.Task, error) {
	r.callCounts["GetUnscheduledTasks"]++
	return nil, nil
}
func (r stubRepo) GetUserBySessionID(sessionID string) (sqlite.User, error) {
	r.callCounts["GetUserBySessionID"]++
	return sqlite.User{}, nil
}
func (r stubRepo) GetUserByUsername(username string) (sqlite.User, error) {
	r.callCounts["GetUserByUsername"]++
	return sqlite.User{}, nil
}
func (r stubRepo) GetUserCategorySuggestions(userID int) ([]string, error) {
	r.callCounts["GetUserCategorySuggestions"]++
	return nil, nil
}
func (r stubRepo) GetWeekOfTasks(userID int, startDay time.Time) ([]sqlite.Task, error) {
	r.callCounts["GetWeekOfTasks"]++
	return nil, nil
}
func (r stubRepo) SignupUserIfNotTaken(username, password string) (sqlite.User, error) {
	r.callCounts["SignupUserIfNotTaken"]++
	return sqlite.User{}, nil
}
func (r stubRepo) StartSession(sessionID string, userID int) error {
	r.callCounts["StartSession"]++
	return nil
}
func (r stubRepo) TaskToggleCompleteIfOwned(taskID, userID int) error {
	r.callCounts["TaskToggleCompleteIfOwned"]++
	return nil
}
func (r stubRepo) TaskUpdateDueDateIfOwned(taskID, userID int, day time.Time) (sqlite.Task, error) {
	r.callCounts["TaskUpdateDueDateIfOwned"]++
	return sqlite.Task{}, nil
}
func (r stubRepo) UpdateTaskIfOwned(t sqlite.Task) (sqlite.Task, error) {
	r.callCounts["UpdateTaskIfOwned"]++
	return sqlite.Task{}, nil
}

type mockUpdateTask struct {
	stubRepo
	returnTask sqlite.Task
	returnErr  error
}

func (m mockUpdateTask) UpdateTaskIfOwned(_ sqlite.Task) (sqlite.Task, error) {
	return m.returnTask, m.returnErr
}

type mockTaskUpdateDueDate struct {
	stubRepo
	returnTask sqlite.Task
	returnErr  error
}

func (m mockTaskUpdateDueDate) TaskUpdateDueDateIfOwned(_, _ int, _ time.Time) (sqlite.Task, error) {
	return m.returnTask, m.returnErr
}

type mockTaskToggleComplete struct {
	stubRepo
	returnTask sqlite.Task
	returnErr  error
}

func (m mockTaskToggleComplete) TaskToggleCompleteIfOwned(_, _ int) error {
	return m.returnErr
}

type mockLoginUser struct {
	stubRepo
	returnUser       sqlite.User
	returnUserErr    error
	returnSessionErr error
}

func (m mockLoginUser) GetUserByUsername(_ string) (sqlite.User, error) {
	return m.returnUser, m.returnUserErr
}
func (m mockLoginUser) StartSession(_ string, _ int) error { return m.returnSessionErr }

type mockSignupUser struct {
	stubRepo
	returnUser       sqlite.User
	returnUserErr    error
	returnSessionErr error
}

func (m mockSignupUser) SignupUserIfNotTaken(_, _ string) (sqlite.User, error) {
	return m.returnUser, m.returnUserErr
}
func (m mockSignupUser) StartSession(_ string, _ int) error { return m.returnSessionErr }

type mockWeekOfTasks struct {
	stubRepo
	returnTasks []sqlite.Task
	returnErr   error
}

func (m mockWeekOfTasks) GetWeekOfTasks(_ int, _ time.Time) ([]sqlite.Task, error) {
	return m.returnTasks, m.returnErr
}

type mockUserCategorySuggestions struct {
	stubRepo
	returnSuggestions []string
	returnErr         error
}

func (m mockUserCategorySuggestions) GetUserCategorySuggestions(_ int) ([]string, error) {
	return m.returnSuggestions, m.returnErr
}

type mockUserByUsername struct {
	stubRepo
	returnUser sqlite.User
	returnErr  error
}

func (m mockUserByUsername) GetUserByUsername(_ string) (sqlite.User, error) {
	return m.returnUser, m.returnErr
}

type mockUserBySessionID struct {
	stubRepo
	returnUser sqlite.User
	returnErr  error
}

func (m mockUserBySessionID) GetUserBySessionID(_ string) (sqlite.User, error) {
	return m.returnUser, m.returnErr
}

type mockUnscheduledTasks struct {
	stubRepo
	returnTasks []sqlite.Task
	returnErr   error
}

func (m mockUnscheduledTasks) GetUnscheduledTasks(_ int) ([]sqlite.Task, error) {
	return m.returnTasks, m.returnErr
}

type mockGetTasksForDay struct {
	stubRepo
	returnTasks []sqlite.Task
	returnErr   error
}

func (m mockGetTasksForDay) GetTasksForDay(_ int, _ time.Time) ([]sqlite.Task, error) {
	return m.returnTasks, m.returnErr
}

type mockGetTaskByID struct {
	stubRepo
	returnTask sqlite.Task
	returnErr  error
}

func (m mockGetTaskByID) GetTaskByID(_, _ int) (sqlite.Task, error) { return m.returnTask, m.returnErr }

type mockGetOverdueTasks struct {
	stubRepo
	returnTasks []sqlite.Task
	returnErr   error
}

func (m mockGetOverdueTasks) GetOverdueTasks(_ int) ([]sqlite.Task, error) {
	return m.returnTasks, m.returnErr
}

type mockGetDashboardData struct {
	stubRepo
	returnWeek           []sqlite.Task
	returnWeekErr        error
	returnOverdue        []sqlite.Task
	returnOverdueErr     error
	returnUnscheduled    []sqlite.Task
	returnUnscheduledErr error
}

func (m mockGetDashboardData) GetWeekOfTasks(_ int, _ time.Time) ([]sqlite.Task, error) {
	return m.returnWeek, m.returnWeekErr
}
func (m mockGetDashboardData) GetOverdueTasks(_ int) ([]sqlite.Task, error) {
	return m.returnOverdue, m.returnOverdueErr
}
func (m mockGetDashboardData) GetUnscheduledTasks(_ int) ([]sqlite.Task, error) {
	return m.returnUnscheduled, m.returnUnscheduledErr
}

type mockDeleteSession struct {
	stubRepo
	returnErr error
}

func (m mockDeleteSession) DeleteSession(_ string) error { return m.returnErr }

type mockDeleteTask struct {
	stubRepo
	returnErr error
}

func (m mockDeleteTask) DeleteTaskIfOwned(_, _ int) error { return m.returnErr }

type mockCompleteTask struct {
	stubRepo
	returnErr error
}

func (m mockCompleteTask) CompleteTaskIfOwned(_, _ int) error { return m.returnErr }

type mockAddTask struct {
	stubRepo
	returnTask sqlite.Task
	returnErr  error
}

func (m mockAddTask) AddTask(_ sqlite.Task) (sqlite.Task, error) { return m.returnTask, m.returnErr }

var randomErr = errors.New("unknown error")

func wantNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("wanted: nil, got: %v", err)
	}
}
func wantErrIs(target error) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		t.Helper()
		if !errors.Is(err, target) {
			t.Errorf("wanted: %v, got: %v", target, err)
		}
	}
}
func wantErrTaskValidation(wantedField services.TaskField, wantedMessage string) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		t.Helper()
		target := services.AsErrTaskValidation(err)
		if target == nil {
			t.Errorf("wanted ErrTaskValidation, got: %v (type %T)", err, err)
			return
		}
		if target.Field != wantedField {
			t.Errorf("wanted Field: %s, got: %s", wantedField, target.Field)
		}
		if target.Message != wantedMessage {
			t.Errorf("wanted Message: %s, got: %s", wantedMessage, target.Message)
		}
	}
}

func wantErrUserValidation(wantedField services.UserField, wantedMessage string) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		t.Helper()
		target := services.AsErrUserValidation(err)
		if target == nil {
			t.Errorf("wanted ErrUserValidation, got: %v (type %T)", err, err)
			return
		}
		if target.Field != wantedField {
			t.Errorf("wanted Field: %s, got: %s", wantedField, target.Field)
		}
		if target.Message != wantedMessage {
			t.Errorf("wanted Message: %s, got: %s", wantedMessage, target.Message)
		}
	}
}

func mustParseDateTime(t string) time.Time {
	date, err := time.ParseInLocation("2006-01-02 15:04", t, time.Local)
	if err != nil {
		panic(err)
	}

	return date
}

func mustParseDate(t string) time.Time {
	date, err := time.ParseInLocation("2006-01-02", t, time.Local)
	if err != nil {
		panic(err)
	}

	return date
}

func checkTask(t *testing.T, gotTask, wantedTask services.Task) {
	if gotTask != wantedTask {
		if gotTask.ID != wantedTask.ID {
			t.Errorf("Task ID wanted: %d, got: %d", wantedTask.ID, gotTask.ID)
		}
		if gotTask.UserID != wantedTask.UserID {
			t.Errorf("Task UserID wanted: %d, got: %d", wantedTask.UserID, gotTask.UserID)
		}
		if gotTask.Title != wantedTask.Title {
			t.Errorf("Task Title wanted: %s, got: %s", wantedTask.Title, gotTask.Title)
		}
		if gotTask.Category != wantedTask.Category {
			t.Errorf("Task Category wanted: %s, got: %s", wantedTask.Category, gotTask.Category)
		}
		if gotTask.Description != wantedTask.Description {
			t.Errorf("Task Description wanted: %s, got: %s", wantedTask.Description, gotTask.Description)
		}
		if gotTask.DueDate.Format("Mon Jan 2 2006") != wantedTask.DueDate.Format("Mon Jan 2 2006") {
			t.Errorf("Task DueDate wanted: %s, got: %s", wantedTask.DueDate.Format("Mon Jan 2 2006"), gotTask.DueDate.Format("Mon Jan 2 2006"))
		}
		if gotTask.CompletionDate != wantedTask.CompletionDate {
			t.Errorf("Task CompletionDate wanted: %s, got: %s", wantedTask.CompletionDate.Format("Mon Jan 2 2006"), gotTask.CompletionDate.Format("Mon Jan 2 2006"))
		}
		if gotTask.Done != wantedTask.Done {
			t.Errorf("Task Done wanted: %t, got: %t", wantedTask.Done, gotTask.Done)
		}
		if gotTask.RecurringPeriod != wantedTask.RecurringPeriod {
			t.Errorf("Task RecurringPeriod wanted: %d, got: %d", wantedTask.RecurringPeriod, gotTask.RecurringPeriod)
		}
	}
}

func checkUser(t *testing.T, gotUser, wantedUser services.User) {
	if gotUser != wantedUser {
		if gotUser.Username != wantedUser.Username {
			t.Errorf("User Username wanted: %s, got %s", wantedUser.Username, gotUser.Username)
		}
		if gotUser.Password != wantedUser.Password {
			t.Errorf("User Password wanted: %s, got %s", wantedUser.Password, gotUser.Password)
		}
		if gotUser.ID != wantedUser.ID {
			t.Errorf("User ID wanted: %d, got %d", wantedUser.ID, gotUser.ID)
		}
		if gotUser.CreatedAt.Compare(wantedUser.CreatedAt) != 0 {
			t.Errorf("User CreatedAt wanted: %s, got %s", wantedUser.CreatedAt.Format("2006-01-02"), gotUser.CreatedAt.Format("2006-01-02"))
		}
		if gotUser.UpdatedAt.Compare(wantedUser.UpdatedAt) != 0 {
			t.Errorf("User UpdatedAt wanted: %s, got %s", wantedUser.UpdatedAt.Format("2006-01-02"), gotUser.UpdatedAt.Format("2006-01-02"))
		}
	}
}
