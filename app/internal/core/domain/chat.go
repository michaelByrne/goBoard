package domain

import "time"

type ChatGroup struct {
	ID    int    `json:"id"`
	Topic string `json:"topic"`
	Chats []Chat `json:"chats"`
}

type Chat struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	MemberID  int       `json:"member_id"`
	Chat      string    `json:"chat"`
	Timestamp time.Time `json:"timestamp"`
}
