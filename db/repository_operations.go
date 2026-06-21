package db

import (
	"time"
)

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