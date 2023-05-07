package domain

import "time"

type Post struct {
	ID        string
	Timestamp *time.Time
	AuthorID  string
	ThreadID  string
	Text      string
}
