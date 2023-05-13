package domain

import "time"

type Post struct {
	ID        string
	Timestamp *time.Time
	MemberID  int
	MemberIP  string
	ThreadID  int
	Text      string
}
