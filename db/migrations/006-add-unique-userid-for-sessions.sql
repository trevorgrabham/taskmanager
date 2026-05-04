CREATE TABLE new_session (
  id TEXT PRIMARY KEY
    CHECK (id != ""),
  user_id INTEGER NOT NULL UNIQUE,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
  expires_at INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES user(id)
);

INSERT INTO new_session (id, user_id, created_at, expires_at) 
SELECT id, user_id, created_at, expires_at FROM session;

DROP TABLE session;

ALTER TABLE new_session RENAME TO session;

CREATE INDEX idx_session_user_id ON session(user_id);
