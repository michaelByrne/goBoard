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
	SaveThread(thread domain.Thread) (int, error)
	ListPostsForThread(limit, offset, id, memberID int) ([]domain.ThreadPost, error)
	ToggleDot(ctx context.Context, memberID, threadID int) (bool, error)
	ToggleIgnore(ctx context.Context, memberID, threadID int) (bool, error)
	ToggleFavorite(ctx context.Context, memberID, threadID int) (bool, error)
	ListThreads(ctx context.Context, cursors domain.Cursors, limit, memberID int, filter domain.ThreadFilter) ([]domain.Thread, domain.Cursors, error)
	ListPostsCollapsible(ctx context.Context, toShow, threadID, memberID int) (posts []domain.ThreadPost, collapsed int, err error)
	ViewThread(ctx context.Context, memberID, threadID int) (int, error)
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
	GetMessagePostsCollapsible(ctx context.Context, viewable, messageID, memberID int) ([]domain.MessagePost, int, error)
	GetMessageByID(ctx context.Context, messageID, memberID int) (*domain.Message, error)
	ListMessages(ctx context.Context, cursors domain.Cursors, limit, memberID int) ([]domain.Message, domain.Cursors, error)
	ViewMessage(ctx context.Context, memberID, messageID int) (int, error)
	DeleteMessage(ctx context.Context, memberID, messageID int) error
	GetMessageParticipants(ctx context.Context, messageID int) ([]string, error)
	GetNewMessageCounts(ctx context.Context, memberID int) (*domain.MessageCounts, error)
}

//go:generate moq -pkg mocks -out ../service/mocks/authentication_repo_moq.go . AuthenticationRepo

type AuthenticationRepo interface {
	Authenticate(ctx context.Context, username, password string) (*domain.Token, error)
}

//go:generate moq -pkg mocks -out ../service/mocks/image_repo_moq.go . ImageRepo

type ImageRepo interface {
	UploadImage(ctx context.Context, imageBytes []byte) (*domain.Image, error)
	ResizeImage(imageBytes []byte) ([]byte, error)
	PresignURL(ctx context.Context, key string) (string, error)
}

type ThemeRepo interface {
	GetTheme(ctx context.Context, name string) (string, error)
}
