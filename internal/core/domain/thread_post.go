package domain

import "time"

type ThreadPost struct {
	ID            int
	Timestamp     *time.Time
	MemberID      int
	MemberIP      string
	ParentID      int
	Body          string
	ParentSubject string
	IsAdmin       bool
	MemberName    string
	Position      int
}
