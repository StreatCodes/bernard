package database

type Host struct {
	ID          int64  `db:"id"`
	Domain      string `db:"domain"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Key         []byte `db:"key"`
}

func (db *DB) CreateHost(domain, name, description string) (int64, error) {
	token, err := generateToken()
	if err != nil {
		return -1, err
	}

	res, err := db.RawDB.Exec(`INSERT INTO hosts (domain, name, description, key) VALUES (?, ?, ?, ?)`, domain, name, description, token)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (db *DB) SelectHostById(id int64) (Host, error) {
	host := Host{}
	err := db.RawDB.Get(&host, `SELECT * FROM hosts WHERE id = ?`, id)

	return host, err
}

func (db *DB) SelectHostByKey(key []byte) (Host, error) {
	host := Host{}
	err := db.RawDB.Get(&host, `SELECT * FROM hosts WHERE key = ?`, key)

	return host, err
}

func (db *DB) SearchHostByDomain(query string) ([]Host, error) {
	hosts := []Host{}
	err := db.RawDB.Select(&hosts, `SELECT * FROM hosts WHERE domain LIKE ? COLLATE NOCASE`, "%"+query+"%")

	return hosts, err
}
