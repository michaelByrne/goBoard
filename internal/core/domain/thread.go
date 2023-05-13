package domain

import "time"

type Thread struct {
	ID             int
	Subject        string
	MemberID       int
	MemberName     string
	Views          int
	LastPosterName string
	LastPosterID   int
	LastPostText   string
	DateLastPosted string
	Sticky         bool
	Locked         bool
	Legendary      bool
	NumPosts       int
	Posts          []Post
	Timestamp      *time.Time
}
