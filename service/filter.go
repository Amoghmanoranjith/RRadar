package service

import (
	"database/sql"
	"errors"
	"fmt"
	"rradar/db"
	modelXML "rradar/model/xml"
	"slices"
	"time"
)

func Filter(repo *db.Repository, feed modelXML.Feed) (filteredFeed modelXML.Feed){

	// we filter
	// find the last published post stored in db
	published, _, err := repo.GetPost(feed.Subreddit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			published = time.Time{} // id want this to be equal to oldest so that i can analyse all the posts once
		} else {
			panic(err) 
		}
	}
	// first sort the entries in latest published at top
	slices.SortFunc(feed.Entries, func(a, b modelXML.Entry) int {
		if a.Published.After(b.Published) {
			return -1
		}
		if a.Published.Before(b.Published) {
			return 1
		}
		return 0
	})

	// stroe the entries in a separate array
	var filteredEntries []modelXML.Entry

	// then find the first entry that occured after published
	for _, entry := range feed.Entries {
		if entry.Published.After(published) {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	// we sys out these entries the others are not required
	if len(filteredEntries) == 0 {
		fmt.Println("No new entries")
		return modelXML.Feed{}
	}
	// store in an array of data
	for _, entry := range filteredEntries {
		fmt.Println("================================")
		fmt.Println("Title    :", entry.Title)
		fmt.Println("Author   :", entry.Author)
		fmt.Println("Published:", entry.Published)
		fmt.Println("Link     :", entry.Link)

	}
	feed.Entries = filteredEntries
	return feed
}

