CREATE TABLE new_user (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE
    CHECK (username != ''),
  hashed_password TEXT NOT NULL
    CHECK (hashed_password != ''),
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

INSERT INTO new_user (id, username, hashed_password, created_at, updated_at)
SELECT id, username, hashed_password, created_at, updated_at FROM user;

DROP TABLE user;

ALTER TABLE new_user RENAME TO user;
