package model

import "time"

type ClassifiedPost struct {
	Title       string
	Content     string
	Link        string
	Author      string
	Published   time.Time
	Subreddit   string
	Interesting bool
	Confidence  float32
	Reason      string
}
