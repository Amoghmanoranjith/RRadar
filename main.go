package main

import (
	"rradar/db"
	"rradar/model"
	"rradar/service"
	"time"
)

func main() {
	repo, err := db.New("posts.db")
	if err != nil {
		panic(err)
	}
	defer repo.Close()
	// init ticker
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	// 1st minute subreddit 1
	// after that subreddit 1 is done wait for 1 minute run for 2nd subreddit
	for i := 0; i < len(model.Subreddits); i++ {
		subreddit := model.Subreddits[i]
		feed := service.Scrape(subreddit)

		// repo.UpdatePost with the latest one among all
		latest := feed.Entries[0]
		err = repo.UpdatePost(feed.Subreddit, latest.Author, latest.Published)
		if err != nil {
			panic(err)
		}

		// process this dataArray
		time.Sleep(time.Minute)
		if i == len(model.Subreddits)-1 {
			i = -1 // reset to start
		}
	}
	// _ = repo.Drop()
}
