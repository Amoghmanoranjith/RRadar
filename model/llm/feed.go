package llm

import (
		modelXML "rradar/model/xml"
)

type Feed struct {
	Subreddit string
	Entries   []Entry
}

func NewFeedWithEmptyEntries(f modelXML.Feed) Feed {
	feed := Feed{
		Subreddit: f.Subreddit,
		Entries:   make([]Entry, 0, len(f.Entries)),
	}
	return feed
}

