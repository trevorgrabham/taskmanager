CREATE TABLE recurring (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  period TEXT UNIQUE NOT NULL
);

CREATE TABLE task (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  category TEXT, 
  description TEXT,
  due_date INTEGER,
  completion_date INTEGER,
  done INTEGER NOT NULL DEFAULT 0
    CHECK (done IN (0, 1)),
  recurring_id INTEGER,
  FOREIGN KEY (recurring_id) REFERENCES recurring(id)
);

