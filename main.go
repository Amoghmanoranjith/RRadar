package main

import (
	"log/slog"
	"rradar/adapter"
	"rradar/config"
	"rradar/db"
	"rradar/domain"
	"rradar/logger"
	"rradar/model"
	"rradar/orchestrator"
	"rradar/service"
	classifierStrategy "rradar/strategy/classifier"
	notifierStrategy "rradar/strategy/notifier"
	"time"
)

func main() {
	logger, err := logger.GetLogger()
	if err != nil {
		panic(err)
	}

	slog.SetDefault(logger)

	slog.Info("Starting Reddit Radar")

	// create a single db and use everywhere
	repo, err := db.New("posts.db")
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		panic(err)
	}
	defer func() {
		repo.Close()
		slog.Info("Database connection closed")
	}()
	// load configuration
	config, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		panic(err)
	}

	slog.Info("Configuration loaded")

	// init classifier strategies
	classifiers := []domain.Classifier{
		classifierStrategy.NewMistral(config.MistralAPIKey),
		classifierStrategy.NewGemini25Flash(config.GeminiAPIKey),
	}

	// init notifier strategy
	notifier := domain.Notifier(
		notifierStrategy.NewDiscord(config.DiscordWebhook),
	)

	// init orchestrators
	classificationOrchestrator := orchestrator.NewClassificationOrchestrator(classifiers)
	notificationOrchestrator := orchestrator.NewNotificationOrchestrator(notifier)

	slog.Info(
		"Application initialized",
		"classifiers", len(classifiers),
		"subreddits", len(model.Subreddits),
	)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for i := 0; i < len(model.Subreddits); i++ {
		subreddit := model.Subreddits[i]

		slog.Info(
			"Processing subreddit",
			"subreddit", subreddit,
		)

		// scrape subreddit
		feed := service.Scrape(subreddit)

		posts, err := adapter.AdaptXMLToPost(feed)
		if err != nil {
			slog.Error(
				"Failed to parse feed",
				"subreddit", subreddit,
				"error", err,
			)
			panic(err)
		}

		slog.Info(
			"Feed parsed",
			"subreddit", subreddit,
			"posts", len(posts),
		)

		filteredPosts, err := service.Filter(repo, posts)
		if err != nil {
			slog.Error(
				"Failed to filter posts",
				"subreddit", subreddit,
				"error", err,
			)
			panic(err)
		}

		if len(filteredPosts) == 0 {
			slog.Info(
				"No new posts",
				"subreddit", subreddit,
			)

			time.Sleep(time.Minute)

			if i == len(model.Subreddits)-1 {
				i = -1
			}

			continue
		}

		slog.Info(
			"New posts found",
			"subreddit", subreddit,
			"count", len(filteredPosts),
		)

		for _, post := range filteredPosts {
			slog.Info(
				"Filter result",
				"title", post.Title,
				"author", post.Author,
				"published", post.Published,
				"link", post.Link,
			)
		}

		// update latest processed post
		latest := filteredPosts[0]

		err = repo.UpdatePost(
			latest.Subreddit,
			latest.Author,
			latest.Published,
		)
		if err != nil {
			slog.Error(
				"Failed to update database",
				"subreddit", subreddit,
				"error", err,
			)
			panic(err)
		}

		slog.Info(
			"Starting classification",
			"subreddit", subreddit,
			"posts", len(filteredPosts),
		)

		classifiedPosts, err := classificationOrchestrator.Classify(filteredPosts)
		if err != nil {
			slog.Error(
				"Classification failed",
				"subreddit", subreddit,
				"error", err,
			)
			panic(err)
		}

		for _, post := range classifiedPosts {
			slog.Info(
				"Classification result",
				"title", post.Title,
				"interesting", post.Interesting,
				"confidence", post.Confidence,
				"reason", post.Reason,
				"author", post.Author,
				"published", post.Published,
				"link", post.Link,
			)
		}

		slog.Info(
			"Sending notifications",
			"subreddit", subreddit,
		)

		err = notificationOrchestrator.Notify(classifiedPosts)
		if err != nil {
			slog.Error(
				"Notification failed",
				"subreddit", subreddit,
				"error", err,
			)
			panic(err)
		}

		slog.Info(
			"Finished processing subreddit",
			"subreddit", subreddit,
		)

		time.Sleep(time.Minute)

		if i == len(model.Subreddits)-1 {
			i = -1
		}
	}
}
