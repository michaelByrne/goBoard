package domain

import (
	"github.com/gorilla/sessions"
	"time"
)

type SiteContext struct {
	ThreadPage     ThreadPage
	Messages       []Message
	Member         Member
	Thread         Thread
	PageName       string
	PageCursor     *time.Time
	PrevPageCursor *time.Time
	Session        *sessions.Session
	Prefs          []Pref
	Username       string
}
