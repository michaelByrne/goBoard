package domain

import "time"

type Thread struct {
	ID             int
	Subject        string
	MemberID       int
	MemberName     string
	FirstPostText  string
	MemberIP       string
	Views          int
	LastPosterName string
	LastPosterID   int
	LastPostText   string
	DateLastPosted *time.Time
	DatePosted     *time.Time
	Sticky         bool
	Locked         bool
	Legendary      bool
	NumPosts       int
	Posts          []ThreadPost
	Timestamp      *time.Time
	PageCursor     *time.Time
}
