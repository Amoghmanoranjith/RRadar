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
			subreddit TEXT PRIMARY KEY,
			author TEXT,
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

func (r *Repository) GetPost(subreddit string) (published time.Time, author string, err error) {
	err = r.db.QueryRow(`
		SELECT published, author
		FROM post
		WHERE subreddit = ?
	`, subreddit).Scan(&published, &author)
	return
}

func (r *Repository) UpdatePost(
	subreddit string,
	author string,
	published time.Time,
) error {
	_, err := r.db.Exec(`
		INSERT OR REPLACE INTO post (
			subreddit,
			author,
			published
		)
		VALUES (?, ?, ?)
	`, subreddit, author, published)

	return err
}
