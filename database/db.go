package database

import (
	"fmt"
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

func ConnectFileDB(dbFileName string) (*DB, error) {
	db, err := sqlx.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, err
	}

	return &DB{RawDB: db}, err
}

func SetupDB(dbFileName, initFile string) error {
	db, err := sqlx.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	b, err := ioutil.ReadFile(initFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(b))
	if err != nil {
		return err
	}

	dbWrap := DB{RawDB: db}

	var email, name, password string
	fmt.Println("Enter an email for the admin account:")
	fmt.Scanln(&email)
	fmt.Println("Enter a name for the admin account:")
	fmt.Scanln(&name)
	fmt.Println("Enter the password for the admin account:")
	fmt.Scanln(&password)

	_, err = dbWrap.CreateUser(email, name, password)
	if err != nil {
		return err
	}

	return nil
}
