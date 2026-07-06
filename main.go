package main

import (
	"fmt"
	"rradar/adapter"
	"rradar/config"
	"rradar/db"
	"rradar/domain"
	"rradar/model"
	"rradar/orchestrator"
	"rradar/service"
	classifierStrategy "rradar/strategy/classifier"
	notifierStrategy "rradar/strategy/notifier"
	"time"
)

func main() {
	// create a single db and use everywhere
	repo, err := db.New("posts.db")
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	// use env
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	// init classifier strategies
	classifiers := []domain.Classifier{
		classifierStrategy.NewGemini25Flash(config.GeminiAPIKey),
		classifierStrategy.NewGemini25FlashLite(config.GeminiAPIKey),
	}
	// init notifier strategy
	notifier := domain.Notifier(notifierStrategy.NewDiscord(config.DiscordWebhook))

	// init classification manager
	classificationOrchestrator := orchestrator.NewClassificationOrchestrator(classifiers)

	// init notification manager
	notificationOrchestrator := orchestrator.NewNotificationOrchestrator(notifier)

	// init ticker
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	// 1st minute subreddit 1
	// after that subreddit 1 is done wait for 1 minute run for 2nd subreddit
	for i := 0; i < len(model.Subreddits); i++ {
		subreddit := model.Subreddits[i]
		// scrape that subreddit
		feed := service.Scrape(subreddit)
		// get posts from this feed
		posts, err := adapter.AdaptXMLToPost(feed)

		//filter the feed only return posts that are published after the last in db
		// this might be empty
		filteredPosts, err := service.Filter(repo, posts)
		if len(filteredPosts) == 0 {
			fmt.Println("========================================")
			fmt.Println("No content for ", subreddit)
			fmt.Println("========================================")
			time.Sleep(time.Minute)
			if i == len(model.Subreddits)-1 {
				i = -1 // reset to start
			}
			continue
		}
		// update the db for this subreddit
		latest := filteredPosts[0]
		err = repo.UpdatePost(latest.Subreddit, latest.Author, latest.Published)
		if err != nil {
			panic(err)
		}
		// use the classification strategy manager to classify each of the feed's entries
		fmt.Println("========================================")
		fmt.Println("Classifying ", subreddit)
		fmt.Println("========================================")
		// fix the classifier part
		classifiedPosts, err := classificationOrchestrator.Classify(filteredPosts)
		if err != nil {
			panic(err)
		}
		for _, classifiedPost := range classifiedPosts {
			fmt.Println("========================================")
			fmt.Println("Title       :", classifiedPost.Title)
			fmt.Println("Interesting :", classifiedPost.Interesting)
			fmt.Println("Confidence  :", classifiedPost.Confidence)
			fmt.Println("Reason      :", classifiedPost.Reason)
			fmt.Println("Author      :", classifiedPost.Author)
			fmt.Println("Published   :", classifiedPost.Published)
			fmt.Println("Link        :", classifiedPost.Link)
			fmt.Println()
		}

		err = notificationOrchestrator.Notify(classifiedPosts)
		if err != nil {
			panic(err)
		}
		// notify the feed
		time.Sleep(time.Minute)
		if i == len(model.Subreddits)-1 {
			i = -1 // reset to start
		}
	}
}
