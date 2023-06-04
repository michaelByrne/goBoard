package domain

import (
	"github.com/gorilla/sessions"
	"time"
)

type SiteContext struct {
	ThreadPage     ThreadPage
	MemberPage     Member
	Member         Member
	PageName       string
	PageCursor     *time.Time
	PrevPageCursor *time.Time
	Session        *sessions.Session
	Prefs          []Pref
	Username       string
}
