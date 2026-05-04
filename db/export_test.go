package db

import (
	"fmt"
	"testing"
)

func SeedDB(r Repo) error {
	_, err := r.db.Exec(seedData)
	return err
}

func ResetDB(t *testing.T, r Repo) {
	t.Helper()
	if _, err := r.db.Exec(`
		DELETE FROM task;
		DELETE FROM session;
		DELETE FROM recurring;
		DELETE FROM user;
		DELETE FROM sqlite_sequence WHERE name IN ('task', 'session', 'recurring', 'user');
		`); err != nil {
		t.Fatal(fmt.Errorf("deleting data: %w", err))
	}
	if err := SeedDB(r); err != nil {
		t.Fatal(fmt.Errorf("seeding data: %w", err))
	}
}

var NextID = 10

// UserID [1, 4], TaskID [1, 9], RecurringID [1, 2].
// Categories: User1("go tests", "general tests", ""), User2("go tests"), User3("go tests"), User4().
// RecurringID: Only id: 1.
// Overdue Tasks: User3.
// Unscheduled Tasks: User3.
var seedData = `
	INSERT INTO user (username, hashed_password, created_at, updated_at) VALUES ("tester1", "password1", strftime('%s', '1997-10-19'), strftime('%s', '1997-10-19'));
	INSERT INTO user (username, hashed_password, created_at, updated_at) VALUES ("tester2", "password2", strftime('%s', 'now'), strftime('%s', 'now'));
	INSERT INTO user (username, hashed_password, created_at, updated_at) VALUES ("tester3", "password3", strftime('%s', 'now'), strftime('%s', 'now'));
	INSERT INTO user (username, hashed_password, created_at, updated_at) VALUES ("tester4", "password4", strftime('%s', 'now'), strftime('%s', 'now'));

	INSERT INTO recurring (period) VALUES ("2 days");
	INSERT INTO recurring (period) VALUES ("4 weeks");

	INSERT INTO task (title, category, description, due_date, user_id, recurring_id) VALUES (
	-- TaskID: 1
	"First Task",
	"go tests",
	"First task's description",
	strftime('%s', 'now') + 3600 * 24 * 7,
	1,
	NULL
	);
	-- TaskID: 2
	INSERT INTO task (title, category, description, due_date, user_id, recurring_id) VALUES (
	"Second Task",
	"go tests",
	"Another task description",
	strftime('%s',
	'now') + 3600 * 24 * 3,
	2,
	NULL
	);
	-- TaskID: 3
	INSERT INTO task (title, category, description, due_date, user_id, recurring_id) VALUES (
	"Empty Description", 
	"general tests", 
	NULL, 
	strftime('%s', 'now') + 3600 * 24 * 2, 
	1, 
	NULL
	);
	-- TaskID: 4
	INSERT INTO task (title, category, description, due_date, user_id, recurring_id) VALUES (
	"With RecurringPeriod", 
	"go tests", 
	"Seeing if the recurring_id works", 
	strftime('%s', 'now') + 3600 * 24 * 2, 
	3, 
	1
	);
	-- TaskID: 5
	INSERT INTO task (title, category, description, due_date, user_id, recurring_id) VALUES (
	"Empty Category", 
	NULL,
	"Testing empty category", 
	strftime('%s', 'now') + 3600 * 24 * 10, 
	1, 
	NULL
	);
	-- TaskID: 6
	INSERT INTO task (title, category, description, due_date, user_id, recurring_id) VALUES (
	"Unscheduled Task", 
	"go tests", 
	"For checking unscheduled", 
	NULL,
	3, 
	NULL
	);
	-- TaskID: 7
	INSERT INTO task (title, category, description, due_date, user_id, recurring_id) VALUES (
	"Overdue Task", 
	"go tests", 
	"For checking overdue", 
	strftime('%s', 'now') - 3600 * 24 * 2, 
	3, 
	NULL
	);
	-- TaskID: 8
	INSERT INTO task (title, category, description, due_date, user_id, recurring_id) VALUES (
	"Another in the go tests group for 1", 
	"go tests", 
	"Seeing how the grouping works by category", 
	strftime('%s', 'now') + 3600 * 24 * 21, 
	1, 
	NULL
	);
	-- TaskID: 9
	INSERT INTO task (title, category, description, due_date, user_id, recurring_id, completion_date, done) VALUES (
	"A completed task", 
	"go tests", 
	"For testing with completed tasks", 
	strftime('%s', 'now') + 3600 * 24 * 21, 
	3, 
	NULL,
	strftime('%s', 'now'),
	1
	);

	INSERT INTO session (id, user_id, created_at, expires_at) VALUES ('66ba76c7856f0aea25fb415bfba50721195efbb19e13ccabd7f384f2dfaa62', 1, strftime('%s', 'now') - 3600 * 24 * 2, strftime('%s', 'now') + 3600 * 24 * 28);
	INSERT INTO session (id, user_id, created_at, expires_at) VALUES ('56ba76c7856f0aea25fb415bfba50721195efbb19e13ccabd7f384f2dfaa62', 2, strftime('%s', 'now') - 3600 * 24 * 3, strftime('%s', 'now') + 3600 * 24 * 27);
	INSERT INTO session (id, user_id, created_at, expires_at) VALUES ('36ba76c7856f0aea25fb415bfba50721195efbb19e13ccabd7f384f2dfaa62', 3, strftime('%s', 'now') - 3600 * 24 * 40, strftime('%s', 'now') - 3600 * 24 * 10);
`
