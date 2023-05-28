package ports

import (
	"goBoard/internal/core/domain"
	"time"
)

//go:generate moq -pkg mocks -out ../service/mocks/thread_repo_moq.go . ThreadRepo

type ThreadRepo interface {
	SavePost(post domain.Post) (int, error)
	GetPostByID(id int) (*domain.Post, error)
	GetThreadByID(id int) (*domain.Thread, error)
	ListThreads(limit, offset int) (*domain.SiteContext, error)
	ListThreadsByMemberID(memberID int, limit, offset int) ([]domain.Thread, error)
	SaveThread(thread domain.Thread) (int, error)
	ListPostsForThread(limit, offset, id int) ([]domain.Post, error)
	ListPostsForThreadByCursor(limit, id int, cursor *time.Time) ([]domain.Post, error)
	ListThreadsByCursorForward(limit int, cursor *time.Time) ([]domain.Thread, error)
	ListThreadsByCursorReverse(limit int, cursor *time.Time) ([]domain.Thread, error)
	PeekPrevious(timestamp *time.Time) (bool, error)
}

//go:generate moq -pkg mocks -out ../service/mocks/member_repo_moq.go . MemberRepo

type MemberRepo interface {
	SaveMember(member domain.Member) (int, error)
	GetMemberByID(id int) (*domain.Member, error)
	GetMemberIDByUsername(username string) (int, error)
	GetMemberByUsername(username string) (*domain.SiteContext, error)
}
