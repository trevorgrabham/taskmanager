package db_test

import (
	"errors"
	sqlite "local/taskmanager/db"
	"log"
	"os"
	"testing"
	"time"
)

var r sqlite.Repo
var testingDBFile = ":memory:"

func TestMain(m *testing.M) {
	var err error
	if r, err = sqlite.NewRepo(testingDBFile); err != nil {
		log.Fatal(err)
	}

	if err = sqlite.SeedDB(r); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	r.Close()
	os.Exit(code)
}

func checkTask(t *testing.T, wantedTask, gotTask sqlite.Task) {
	if wantedTask != gotTask {
		if wantedTask.ID != gotTask.ID {
			t.Errorf("TaskID wanted %d, got %d", wantedTask.ID, gotTask.ID)
		}
		if wantedTask.UserID != gotTask.UserID {
			t.Errorf("TaskUserID wanted %d, got %d", wantedTask.UserID, gotTask.UserID)
		}
		if wantedTask.Category != gotTask.Category {
			t.Errorf("TaskCategory wanted %v, got %v", wantedTask.Category, gotTask.Category)
		}
		if wantedTask.Description != gotTask.Description {
			t.Errorf("TaskDescription wanted %v, got %v", wantedTask.Description, gotTask.Description)
		}
		if wantedTask.DueDate.Valid != gotTask.DueDate.Valid {
			t.Errorf("TaskDueDate wanted %v, got %v", wantedTask.DueDate, gotTask.DueDate)
		}
		if wantedTask.CompletionDate.Valid != gotTask.CompletionDate.Valid {
			t.Errorf("TaskCompletionDate wanted %v, got %v", wantedTask.CompletionDate, gotTask.CompletionDate)
		}
		if wantedTask.Done != gotTask.Done {
			t.Errorf("TaskDone wanted %t, got %t", wantedTask.Done, gotTask.Done)
		}
		if wantedTask.RecurringPeriod != gotTask.RecurringPeriod {
			t.Errorf("TaskRecurringPeriod wanted %v, got %v", wantedTask.RecurringPeriod, gotTask.RecurringPeriod)
		}
	}
}

func checkUser(t *testing.T, gotUser, wantedUser sqlite.User) {
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
		if time.Unix(gotUser.CreatedAt, 0).Format("2006-01-02") != time.Unix(wantedUser.CreatedAt, 0).Format("2006-01-02") {
			t.Errorf("User CreatedAt wanted: %s, got %s", time.Unix(wantedUser.CreatedAt, 0).Format("2006-01-02"), time.Unix(gotUser.CreatedAt, 0).Format("2006-01-02"))
		}
		if time.Unix(gotUser.UpdatedAt, 0).Format("2006-01-02") != time.Unix(wantedUser.UpdatedAt, 0).Format("2006-01-02") {
			t.Errorf("User UpdatedAt wanted: %s, got %s", time.Unix(wantedUser.UpdatedAt, 0).Format("2006-01-02"), time.Unix(gotUser.UpdatedAt, 0).Format("2006-01-02"))
		}
	}
}

func wantNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("wanted error: nil, got: %v", err)
	}
}
func wantErrIs(target error) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		t.Helper()
		if !errors.Is(err, target) {
			t.Errorf("wanted error: %v, got: %v", target, err)
		}
	}
}

func mustParseDateTime(t string) time.Time {
	date, err := time.ParseInLocation("2006-01-02 15:04", t, time.UTC)
	if err != nil {
		panic(err)
	}

	return date
}

func mustParseDate(t string) time.Time {
	date, err := time.ParseInLocation("2006-01-02", t, time.UTC)
	if err != nil {
		panic(err)
	}

	return date
}
