package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type Repository struct {
	db *sql.DB
}

func New(path string) (*Repository, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// SQLite generally works best with a single writer connection.
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS post (
			label TEXT PRIMARY KEY,
			user_name TEXT,
			published DATETIME
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}
func (r *Repository) Drop() error {
	_, err := r.db.Exec(`
		DROP TABLE IF EXISTS post
	`)
	return err
}
