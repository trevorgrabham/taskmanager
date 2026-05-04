CREATE TABLE new_recurring (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  period INTEGER NOT NULL
);

INSERT INTO new_recurring (id, period)
SELECT id, CAST(TRIM(SUBSTR(period, 1, INSTR(period, ' ') - 1)) AS INTEGER) * 
  CASE TRIM(SUBSTR(period, INSTR(period, ' ') + 1))
    WHEN 'days' THEN 3600 * 24
    WHEN 'weeks' THEN 3600 * 24 * 7
    WHEN 'months' THEN 3600 * 24 * 30
  END 
FROM recurring;

CREATE TABLE new_task (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL
    CHECK (title != ''),
  category TEXT, 
  description TEXT,
  due_date INTEGER,
  completion_date INTEGER,
  done INTEGER NOT NULL DEFAULT 0
    CHECK (done IN (0, 1)),
  recurring_period INTEGER,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES user(id)
);

INSERT INTO new_task (id, title, category, description, due_date, completion_date, done, recurring_period, user_id)
SELECT id, title, category, description, due_date, completion_date, done, (
  SELECT period FROM new_recurring WHERE new_recurring.id = task.recurring_id 
  ), user_id 
FROM task;

DROP TABLE new_recurring;
DROP TABLE task;

ALTER TABLE new_task RENAME TO task;
