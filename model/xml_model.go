package model

import "time"

type Feed struct {
	Title   string    `xml:"title"`
	Updated time.Time `xml:"updated"`
	Category Category  `xml:"category"`
	Entries []Entry   `xml:"entry"`
}

type Entry struct {
	ID        string    `xml:"id"`
	Title     string    `xml:"title"`
	Content   string    `xml:"content"`
	Updated   time.Time `xml:"updated"`
	Published time.Time `xml:"published"`
	Author    Author    `xml:"author"`
	Link      Link      `xml:"link"`
}
type Category struct {
	Label string `xml:"label,attr"`
}
type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

type Link struct {
	Href string `xml:"href,attr"`
}