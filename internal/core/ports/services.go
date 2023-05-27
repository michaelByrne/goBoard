package ports

import (
	"goBoard/internal/core/domain"
	"time"
)

type ThreadService interface {
	NewPost(body, ip, memberName string, threadID int) (int, error)
	GetPostByID(id int) (*domain.Post, error)
	GetThreadByID(limit, offset, id int) (*domain.Thread, error)
	ListThreads(limit, offset int) (*domain.SiteContext, error)
	NewThread(memberName, memberIP, body, subject string) (int, error)
	GetThreadsWithCursor(limit int, firstPage bool, cursor *time.Time) (*domain.SiteContext, error)
	GetThreadsWithCursorReverse(limit int, cursor *time.Time) (*domain.SiteContext, error)
}

type MemberService interface {
	Save(member domain.Member) (int, error)
	GetMemberByID(id int) (*domain.Member, error)
	GetMemberIDByUsername(username string) (int, error)
	GetMemberByUsername(username string) (*domain.SiteContext, error)
}
