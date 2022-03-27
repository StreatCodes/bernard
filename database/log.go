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

func (db *DB) SelectLogsBetween(start time.Time, end time.Time) ([]Log, error) {
	logs := []Log{}
	err := db.RawDB.Select(&logs, `SELECT * FROM logs WHERE time > ? AND time < ?`, start, end)
	return logs, err
}

func (db *DB) SelectHostLogsBetween(hostID int64, start time.Time, end time.Time) ([]Log, error) {
	logs := []Log{}
	err := db.RawDB.Select(&logs, `SELECT * FROM logs
        WHERE host_id = ? AND time > ? AND time < ?`, hostID, start, end)
	return logs, err
}
