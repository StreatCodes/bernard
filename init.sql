CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	password BLOB NOT NULL
);

CREATE INDEX email_index on users(email);

CREATE TABLE sessions (
	token BLOB PRIMARY KEY,
	user_id INTEGER
);

CREATE TABLE logs (
	time DATETIME PRIMARY KEY,
	host_id INTEGER,
	level INTEGER,
	service TEXT,
	content TEXT
);

CREATE INDEX filter_index on logs(time, host_id);