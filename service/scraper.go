package service

import (
	"fmt"
	"io"
	"log/slog"
	"rradar/http"
)

func Scrape(
	subreddit string,
) []byte {
	url := "https://www.reddit.com/r/" + subreddit + "/new.rss"

	fmt.Println("Fetching:", url)
	resp, err := http.Client.Get(url)
	if err != nil {
		slog.Error(
			"failed to scrape",
			"subreddit", subreddit,
			"operation", "Scrape",
			"cause", "http.Client.Get",
			"error", err,
		)
		return []byte{}
	}

	fmt.Println("Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()	
	if err != nil {
		slog.Error(
			"failed to read data",
			"subreddit", subreddit,
			"operation", "Scrape",
			"cause", "io.ReadAll",
			"error", err,
		)		
		return []byte{}
	}

	return body
}
