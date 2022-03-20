package database

import (
	"io/ioutil"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	RawDB *sqlx.DB
}

func CreateInMemoryDB(initFile string) (*DB, error) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(initFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(b))
	if err != nil {
		return nil, err
	}
	return &DB{RawDB: db}, err
}
