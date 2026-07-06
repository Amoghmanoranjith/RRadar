package service

import (
	"database/sql"
	"errors"
	"fmt"
	"rradar/db"
	"rradar/model"
	"slices"
	"time"
)

// we assume all the posts belong to the same subreddit
func Filter(repo *db.Repository, posts []model.Post) ([]model.Post, error) {

	// we filter
	// find the last published post stored in db
	published, _, err := repo.GetPost(posts[0].Subreddit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			published = time.Time{} // id want this to be equal to oldest so that i can analyse all the posts once
		} else {
			return nil, err
		}
	}
	// first sort the entries in latest published at top
	slices.SortFunc(posts, func(a, b model.Post) int {
		if a.Published.After(b.Published) {
			return -1
		}
		if a.Published.Before(b.Published) {
			return 1
		}
		return 0
	})

	// stroe the entries in a separate array
	var filteredPosts []model.Post

	// then find the first entry that occured after published
	for _, entry := range posts {
		if entry.Published.After(published) {
			filteredPosts = append(filteredPosts, entry)
		}
	}
	// we sys out these entries the others are not required
	fmt.Println("================================")
	fmt.Println("Filtering complete")
	fmt.Println("================================")

	return filteredPosts, nil
}
