package cmd

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"rradar/db"
	"rradar/model"
	"slices"
	"time"
)

func Process(repo *db.Repository,
	subreddit string,
) {
	url := "https://www.reddit.com/r/" + subreddit + "/new.rss"
	// init db
	fmt.Println("Fetching:", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	fmt.Println("Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Read failed:", err)
		return
	}

	var feed model.Feed
	if err := xml.Unmarshal(body, &feed); err != nil {
		fmt.Println("XML parse failed:", err)
		return
	}

	// validation completed moving to logic
	published, _, err := repo.GetPost(feed.Category.Label)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			published = time.Time{} // id want this to be equal to oldest so that i can analyse all the posts once
		} else {
			panic(err)
		}
	}
	// first sort the entries in latest published at top
	slices.SortFunc(feed.Entries, func(a, b model.Entry) int {
		if a.Published.After(b.Published) {
			return -1
		}
		if a.Published.Before(b.Published) {
			return 1
		}
		return 0
	})
	// stroe the entries in a separate array
	var filteredEntries []model.Entry

	// then find the first entry that occured after published
	for _, entry := range feed.Entries {
		if entry.Published.After(published) {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	// we sys out these entries the others are not required
	if len(filteredEntries) == 0 {
		fmt.Println("No new entries")
		return
	}
	for _, entry := range filteredEntries {
		fmt.Println("================================")
		fmt.Println("Title    :", entry.Title)
		fmt.Println("Author   :", entry.Author.Name)
		fmt.Println("Published:", entry.Published)
		fmt.Println("Link     :", entry.Link.Href)
	}
	// repo.UpdatePost with the latest one among all
	latest := filteredEntries[0]
	err = repo.UpdatePost(feed.Category.Label, latest.Author.Name, latest.Published)

	if err != nil {
		panic(err)
	}
}
