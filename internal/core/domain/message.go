package domain

import "time"

type Message struct {
	ID             int
	Subject        string
	FirstPostID    int
	DatePosted     *time.Time
	NumPosts       int
	Views          int
	MemberID       int
	DateLastPosted *time.Time
	MemberIP       string
	Body           string
	RecipientIDs   []int
	MemberName     string
	LastPosterID   int
	LastPosterName string
	Posts          []MessagePost
	PageCursor     *time.Time
	PrevPageCursor *time.Time
	HasPrevPage    bool
	HasNextPage    bool
	RowNumber      int
	NumCollapsed   int
	Participants   []string
}

type MessageCounts struct {
	Unread   int
	NewPosts int
}
