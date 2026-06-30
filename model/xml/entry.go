package xml

import "time"

type Entry struct {
	Title     string
	Content   string
	Link      string
	Author    string
	Published time.Time
}
