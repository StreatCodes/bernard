package database

import (
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	db, err := CreateInMemoryDB("../init.sql")
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}

	id, err := db.CreateUser("john.doe@example.com", "John Doe", "password")
	if err != nil {
		t.Errorf("Error creating user: %s", err)
	}

	user, err := db.SelectUserByCreds("john.doe@example.com", "password")
	if err != nil {
		t.Errorf("Error looking up user: %s", err)
	}

	if user.ID != id {
		t.Errorf("Expected user id (%d) to match returned id (%d)", user.ID, id)
	}
	if user.Name != "John Doe" {
		t.Errorf("Expected user name to be 'John Doe' found %s", user.Name)
	}
}

func TestInvalidPassword(t *testing.T) {
	db, err := CreateInMemoryDB("../init.sql")
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}

	_, err = db.CreateUser("john.doe@example.com", "John Doe", "password")
	if err != nil {
		t.Errorf("Error creating user: %s", err)
	}

	_, err = db.SelectUserByCreds("john.doe@example.com", "invalid")
	if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		t.Errorf("Expected ErrMismatchedHashAndPassword found: %s", err)
	}
}

func TestCreateSession(t *testing.T) {
	db, err := CreateInMemoryDB("../init.sql")
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}

	id, err := db.CreateUser("john.doe@example.com", "John Doe", "password")
	if err != nil {
		t.Errorf("Error creating user: %s", err)
	}

	token, err := db.CreateSession(id)
	if len(token) != TokenLength {
		t.Errorf("Expected token length to be %d found %d", TokenLength, len(token))
	}

	user, err := db.SelectUserByToken(token)
	if user.ID != id {
		t.Errorf("Expected user id (%d) to match returned id (%d)", user.ID, id)
	}
}
