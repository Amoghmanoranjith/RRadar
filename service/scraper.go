package service

import (
	"fmt"
	"io"
	"rradar/http"
)

func Scrape(
	subreddit string,
) []byte {
	url := "https://www.reddit.com/r/" + subreddit + "/new.rss"

	fmt.Println("Fetching:", url)
	resp, err := http.Client.Get(url)
	if err != nil {
		fmt.Println("Scrape failed:", err)
		return []byte{}
	}

	fmt.Println("Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()	
	if err != nil {
		fmt.Println("Read failed:", err)
		return []byte{}
	}

	return body
}
