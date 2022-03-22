package database

type Session struct {
	Token  []byte `db:"token"`
	UserID int64  `db:"user_id"`
}

func (db *DB) CreateSession(userId int64) ([]byte, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := Session{
		Token:  token,
		UserID: userId,
	}

	_, err = db.RawDB.NamedExec(`INSERT INTO sessions (token, user_id)
        VALUES (:token, :user_id)`, session)

	return session.Token, err
}
