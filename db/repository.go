package db

import (
	"database/sql"
	"log/slog"
	"time"

	_ "modernc.org/sqlite"
)

type Repository struct {
	db *sql.DB
}

func New(path string) (*Repository, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		slog.Error(
			"failed to create a db connection",
			"operation", "New",
			"cause", "sql.Open",
			"error", err,
		)
		return nil, err
	}

	// SQLite generally works best with a single writer connection.
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		slog.Error(
			"failed to establish connection with db",
			"operation", "New",
			"cause", "db.Ping",
			"error", err,
		)
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
		slog.Error(
			"failed to execute query",
			"operation", "New",
			"cause", "db.Exec",
			"error", err,
		)
		return nil, err
	}
	r := &Repository{
		db: db,
	}
	// r.Drop() // ************************************8testing rn we need empty dbs
	return r, nil
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
	if err != nil {
		slog.Error(
			"failed to get post from db",
			"operation", "GetPost",
			"cause", "r.db.QueryRow",
			"error", err,
		)
		return time.Time{}, "", err
	}
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
	if err != nil {
		slog.Error(
			"failed to update post in db",
			"operation", "UpdatePost",
			"cause", "r.db.Exec",
			"error", err,
		)
	}
	return err
}
