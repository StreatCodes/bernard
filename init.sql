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

CREATE TABLE hosts (
	id INTEGER PRIMARY KEY,
	domain TEXT NOT NULL,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	key BLOB
);

CREATE INDEX domain_index on hosts(domain);

CREATE TABLE logs (
	time DATETIME PRIMARY KEY,
	host_id INTEGER,
	level INTEGER,
	service TEXT,
	content TEXT
);

CREATE INDEX filter_index on logs(time, host_id);