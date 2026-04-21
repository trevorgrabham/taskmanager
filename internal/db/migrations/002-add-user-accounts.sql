CREATE TABLE user (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE INDEX idx_user_username ON user(username);

INSERT INTO user (username, hashed_password) VALUES ("default_user", "nopass");

CREATE TABLE new_task (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  category TEXT, 
  description TEXT,
  due_date INTEGER,
  completion_date INTEGER,
  done INTEGER NOT NULL DEFAULT 0
    CHECK (done IN (0, 1)),
  recurring_id INTEGER,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (recurring_id) REFERENCES recurring(id),
  FOREIGN KEY (user_id) REFERENCES user(id)
);

INSERT INTO new_task (id, title, category, description, due_date, completion_date, done, recurring_id, user_id)
SELECT id, title, category, description, due_date, completion_date, done, recurring_id, 1 FROM task;

DROP TABLE task;

ALTER TABLE new_task RENAME TO task;
