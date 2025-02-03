package ports

import (
	"context"
	"goBoard/internal/core/domain"
	"html/template"
	"time"
)

type ThreadService interface {
	NewPost(body, ip, memberName string, threadID int) (int, error)
	GetPostByID(id int) (*domain.ThreadPost, error)
	GetThreadByID(limit, offset, id, memberID int) (*domain.Thread, error)
	NewThread(memberName, memberIP, body, subject string) (int, error)
	ConvertPostBodyBbcodeToHtml(body string) (*template.HTML, error)
	ToggleDot(ctx context.Context, memberID, threadID int) (bool, error)
	ToggleIgnore(ctx context.Context, memberID, threadID int) (bool, error)
	ToggleFavorite(ctx context.Context, memberID, threadID int) (bool, error)
	ListThreads(ctx context.Context, cursors domain.Cursors, limit, memberID int, filter domain.ThreadFilter) ([]domain.Thread, domain.Cursors, error)
	GetCollapsibleThreadByID(ctx context.Context, viewable, threadID, memberID int) (*domain.Thread, error)
	ViewThread(ctx context.Context, memberID, threadID int) (int, error)
}

type MemberService interface {
	Save(member domain.Member) (int, error)
	GetMemberByID(id int) (*domain.Member, error)
	GetMemberIDByUsername(username string) (int, error)
	GetMemberByUsername(username string) (*domain.Member, error)
	ValidateMembers(names []string) ([]domain.Member, error)
	GetAllPrefs(ctx context.Context) ([]domain.Pref, error)
	GetMergedPrefs(ctx context.Context, memberID int) ([]domain.Pref, error)
	UpdatePrefs(ctx context.Context, memberID int, updatedPrefs domain.MemberPrefs) error
	UpdateMember(ctx context.Context, member domain.Member) error
	UpdatePostalCode(ctx context.Context, memberID int, postalCode string) error
}

type MessageService interface {
	SendMessage(subject, body, memberIP string, memberID int, recipientIDs []int) (int, error)
	GetMessagesWithCursor(memberID int, reverse bool, cursor *time.Time) ([]domain.Message, error)
	GetMessageByID(messageID, memberID int) (*domain.Message, error)
	NewPost(body, ip, memberName string, messageID int) (int, error)
	GetMessagePostByID(id int) (*domain.MessagePost, error)
	GetCollapsibleMessageByID(ctx context.Context, viewable, messageID, memberID int) (*domain.Message, error)
	ListMessages(ctx context.Context, cursors domain.Cursors, limit, memberID int) ([]domain.Message, domain.Cursors, error)
	ViewMessage(ctx context.Context, memberID, messageID int) (int, error)
	DeleteMessage(ctx context.Context, memberID, messageID int) error
	GetNewMessageCounts(ctx context.Context, memberID int) (*domain.MessageCounts, error)
}

type AuthenticationService interface {
	Authenticate(ctx context.Context, username, password string) (*domain.Member, *domain.Token, error)
}

type ImageService interface {
	UploadImage(ctx context.Context, imageBytes []byte) (string, string, error)
	RefreshPresign(ctx context.Context, key string) (string, error)
	PresignPostImages(ctx context.Context, body string) (string, error)
}

type ThemeService interface {
	GetTheme(ctx context.Context, name string) (string, error)
}
