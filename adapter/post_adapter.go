package adapter

import (
	"encoding/xml"
	"rradar/model"
	"time"
)

type Feed struct {
	Title    string      `xml:"title"`
	Updated  time.Time   `xml:"updated"`
	Category CategoryXML `xml:"category"`
	Entries  []EntryXML  `xml:"entry"`
}

type EntryXML struct {
	ID        string    `xml:"id"`
	Title     string    `xml:"title"`
	Content   string    `xml:"content"`
	Updated   string    `xml:"updated"`
	Published time.Time `xml:"published"`
	Author    AuthorXML `xml:"author"`
	Link      LinkXML   `xml:"link"`
}

type CategoryXML struct {
	Subreddit string `xml:"label,attr"`
}

type AuthorXML struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

type LinkXML struct {
	Href string `xml:"href,attr"`
}

// Parse converts xml feed into domain Posts.
func AdaptXMLToPost(data []byte) ([]model.Post, error) {
	var feed Feed

	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	posts := make([]model.Post, 0, len(feed.Entries))

	for _, entry := range feed.Entries {
		post := model.Post{
			Title:     entry.Title,
			Content:   entry.Content,
			Link:      entry.Link.Href,
			Author:    entry.Author.Name,
			Published: entry.Published,
			Subreddit: feed.Category.Subreddit,
		}

		posts = append(posts, post)
	}

	return posts, nil
}
