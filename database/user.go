package database

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Name     string `db:"name"`
	Password []byte `db:"password"`
}

func (db *DB) CreateUser(email string, name string, password string) (int64, error) {
	bPassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bPassword, bcrypt.DefaultCost)

	if err != nil {
		return -1, err
	}

	res, err := db.RawDB.Exec(`INSERT INTO users (email, name, password) VALUES (?, ?, ?)`, email, name, hash)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (db *DB) SelectUserByCreds(email string, password string) (User, error) {
	user := User{}
	err := db.RawDB.Get(&user, `SELECT * FROM users WHERE email = ?`, email)
	if err != nil {
		return user, err
	}

	bPassword := []byte(password)

	err = bcrypt.CompareHashAndPassword(user.Password, bPassword)

	return user, err
}

func (db *DB) SelectUserByToken(token []byte) (User, error) {
	user := User{}
	err := db.RawDB.Get(&user, `SELECT users.* FROM users
		INNER JOIN sessions ON users.id=sessions.user_id
		WHERE sessions.token = ?`, token)

	return user, err
}
