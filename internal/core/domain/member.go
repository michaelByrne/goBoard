package domain

import "time"

type Member struct {
	ID               int
	Name             string
	Pass             string
	Secret           string
	DateJoined       *time.Time
	FirstPosted      *time.Time
	LastPosted       *time.Time
	LastView         *time.Time
	TotalThreads     int
	TotalThreadPosts int
	Email            string
	PostalCode       string
	Banned           bool
	IP               string
	Prefs            MemberPrefs
	IsAdmin          bool
}

type MemberPref struct {
	ID    int
	Value string
	Type  string
}

type MemberPrefs map[string]MemberPref
