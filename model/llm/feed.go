package llm

type Feed struct {
	Subreddit string
	Entries   []Entry
}

// func (f F) ToFeed() Feed {
// 	feed := Feed{
// 		Subreddit: f.Category.Subreddit,
// 		Entries:   make([]Entry, 0, len(f.Entries)),
// 	}

// 	for _, e := range f.Entries {
// 		feed.Entries = append(feed.Entries, Entry{
// 			Title:     e.Title,
// 			Content:   e.Content,
// 			Link:      e.Link.Href,
// 			Author:    e.Author.Name,
// 			Published: e.Published,
// 		})
// 	}

// 	return feed
// }
