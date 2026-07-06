package model

import "time"

type Post struct {
	Title     string
	Content   string
	Link      string
	Author    string
	Published time.Time
	Subreddit string
}