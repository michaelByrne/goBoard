package domain

import "time"

type Thread struct {
	ID        int
	Subject   string
	MemberID  int
	Views     int
	Posts     []Post
	Timestamp *time.Time
}
