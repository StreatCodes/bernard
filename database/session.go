package database

import (
	"crypto/rand"
)

type Session struct {
	Token  []byte `db:"token"`
	UserID int64  `db:"user_id"`
}

const TokenLength = 64

func (db *DB) CreateSession(userId int64) ([]byte, error) {
	session := Session{
		Token:  make([]byte, TokenLength),
		UserID: userId,
	}
	_, err := rand.Read(session.Token)
	if err != nil {
		return nil, err
	}

	_, err = db.RawDB.NamedExec(`INSERT INTO sessions (token, user_id)
        VALUES (:token, :user_id)`, session)

	return session.Token, err
}
