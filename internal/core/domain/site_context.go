package domain

import (
	"time"
)

type SiteContext struct {
	ThreadPage     ThreadPage
	MemberPage     Member
	Member         Member
	PageName       string
	PageCursor     *time.Time
	PrevPageCursor *time.Time
	//Session        *sessions.Session
}
