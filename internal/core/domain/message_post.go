package domain

import "time"

type MessagePost struct {
	ID            int
	ParentID      int
	Timestamp     *time.Time
	MemberIP      string
	MemberID      int
	MemberName    string
	Body          string
	ParentSubject string
	Position      int
}
