package ports

import (
	"context"
	"goBoard/internal/core/domain"
	"time"
)

//go:generate moq -pkg mocks -out ../service/mocks/thread_repo_moq.go . ThreadRepo

type ThreadRepo interface {
	SavePost(post domain.ThreadPost) (int, error)
	GetPostByID(id int) (*domain.ThreadPost, error)
	GetThreadByID(id, memberID int) (*domain.Thread, error)
	//ListThreads(limit, offset int) (*domain.SiteContext, error)
	ListThreadsByMemberID(memberID int, limit, offset int) ([]domain.Thread, error)
	SaveThread(thread domain.Thread) (int, error)
	ListPostsForThread(limit, offset, id, memberID int) ([]domain.ThreadPost, error)
	PeekPrevious(timestamp *time.Time, memberID int) (bool, error)
	UndotThread(ctx context.Context, memberID, threadID int) error
	DotThread(ctx context.Context, memberID, threadID int) error
	ToggleDot(ctx context.Context, memberID, threadID int) (bool, error)
	ToggleIgnore(ctx context.Context, memberID, threadID int, ignore bool) error
	ListThreads(ctx context.Context, cursors domain.Cursors, limit, memberID int) ([]domain.Thread, domain.Cursors, error)
	ListPostsCollapsible(ctx context.Context, toShow, threadID, memberID int) (posts []domain.ThreadPost, collapsed int, err error)
}

//go:generate moq -pkg mocks -out ../service/mocks/member_repo_moq.go . MemberRepo

type MemberRepo interface {
	SaveMember(member domain.Member) (int, error)
	GetMemberByID(id int) (*domain.Member, error)
	GetMemberIDByUsername(username string) (int, error)
	GetMemberByUsername(username string) (*domain.Member, error)
	GetMemberPrefs(memberID int) (*domain.MemberPrefs, error)
	GetAllPrefs(ctx context.Context) ([]domain.Pref, error)
	UpdatePrefs(ctx context.Context, memberID int, updatedPrefs domain.MemberPrefs) error
	UpdateMember(ctx context.Context, member domain.Member) error
}

//go:generate moq -pkg mocks -out ../service/mocks/message_repo_moq.go . MessageRepo

type MessageRepo interface {
	SaveMessage(message domain.Message) (int, error)
	GetMessagesWithCursorForward(memberID, limit int, cursor *time.Time) ([]domain.Message, error)
	GetMessagesWithCursorReverse(memberID, limit int, cursor *time.Time) ([]domain.Message, error)
	GetMessagePostsByID(memberID, messageID, limit int) ([]domain.MessagePost, error)
	SavePost(post domain.MessagePost) (int, error)
	GetMessagePostByID(id int) (*domain.MessagePost, error)
	PeekPrevious(timestamp *time.Time) (bool, error)
}

//go:generate moq -pkg mocks -out ../service/mocks/authentication_repo_moq.go . AuthenticationRepo

type AuthenticationRepo interface {
	Authenticate(username, password string) (int, error)
}
