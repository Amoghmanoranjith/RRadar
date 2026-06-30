package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"rradar/http"
	modelXML "rradar/model/xml"
)

func Scrape(
	subreddit string,
) modelXML.Feed {
	url := "https://www.reddit.com/r/" + subreddit + "/new.rss"

	fmt.Println("Fetching:", url)
	resp, err := http.Client.Get(url)
	if err != nil {
		fmt.Println("Request failed:", err)
		return modelXML.Feed{}
	}

	fmt.Println("Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Read failed:", err)
		return modelXML.Feed{}
	}

	var feedXML modelXML.Feed_xml
	if err := xml.Unmarshal(body, &feedXML); err != nil {
		fmt.Println("XML parse failed:", err)
		return modelXML.Feed{}
	}

	// Convert XML model -> domain model
	feed := feedXML.ToFeed()

	return feed
}
