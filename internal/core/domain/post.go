package domain

import "time"

type Post struct {
	ID             int
	Timestamp      *time.Time
	MemberID       int
	MemberIP       string
	ThreadID       int
	Text           string
	ThreadSubject  string
	IsAdmin        bool
	MemberName     string
	ThreadPosition int
}
