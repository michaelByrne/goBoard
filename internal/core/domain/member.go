package domain

import "time"

type Member struct {
	ID          int
	Name        string
	Pass        string
	Secret      string
	FirstPosted *time.Time
	LastPosted  *time.Time
	Email       string
	PostalCode  string
	Banned      bool
	IP          string
}
