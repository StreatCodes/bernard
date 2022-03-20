package database

import (
	"time"
)

type Log struct {
	Time    time.Time `db:"time"`
	HostID  int64     `db:"host_id"`
	Level   int64     `db:"level"`
	Service string    `db:"service"`
	Content string    `db:"content"`
}

func (db *DB) InsertLog(log Log) error {
	_, err := db.RawDB.NamedExec(`INSERT INTO logs (time, host_id, level, service, content)
        VALUES (:time, :host_id, :level, :service, :content)`, log)
	return err
}

func (db *DB) SelectSince(since time.Time) ([]Log, error) {
	logs := []Log{}
	err := db.RawDB.Select(&logs, `SELECT * FROM logs WHERE time > ?`, since)
	return logs, err
}

func (db *DB) SelectBetween(start time.Time, end time.Time) ([]Log, error) {
	logs := []Log{}
	err := db.RawDB.Select(&logs, `SELECT * FROM logs WHERE time > ? AND time < ?`, start, end)
	return logs, err
}
