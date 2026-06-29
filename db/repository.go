package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"time"
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

// **********************************
// repo ops

func (r *Repository) GetPost(label string) (published time.Time, userName string, err error) {
	err = r.db.QueryRow(`
		SELECT published, user_name
		FROM post
		WHERE label = ?
	`, label).Scan(&published, &userName)
	return
}

func (r *Repository) UpdatePost(
	label string,
	userName string,
	published time.Time,
) error {
	_, err := r.db.Exec(`
		INSERT OR REPLACE INTO post (
			label,
			user_name,
			published
		)
		VALUES (?, ?, ?)
	`, label, userName, published)

	return err
}
