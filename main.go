package main

import(
	"rradar/model"
	"rradar/db"
	"rradar/cmd"
	"time"
)

func main(){
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
		cmd.Process(repo, subreddit)
		time.Sleep(time.Minute)
		if i == len(model.Subreddits)-1 {
			i = -1 // reset to start
		}
	}
	// _ = repo.Drop()
}