package xml

import "time"

type Feed_xml struct {
	Title    string       `xml:"title"`
	Updated  time.Time    `xml:"updated"`
	Category Category_xml `xml:"category"`
	Entries  []Entry_xml      `xml:"entry"`
}

type Entry_xml struct {
	ID        string     `xml:"id"`
	Title     string     `xml:"title"`
	Content   string     `xml:"content"`
	Updated   time.Time  `xml:"updated"`
	Published time.Time  `xml:"published"`
	Author    Author_xml `xml:"author"`
	Link      Link_xml   `xml:"link"`
}
type Category_xml struct {
	Subreddit string `xml:"label,attr"`
}
type Author_xml struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}
type Link_xml struct {
	Href string `xml:"href,attr"`
}
