CREATE TABLE session (
  id TEXT PRIMARY KEY,
  user_id INTEGER NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
  expires_at INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE INDEX idx_session_user_id ON session(user_id);
