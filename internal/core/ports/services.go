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
	UndotThread(ctx context.Context, memberID, threadID int) error
	DotThread(ctx context.Context, memberID, threadID int) error
	ToggleDot(ctx context.Context, memberID, threadID int) (bool, error)
	ToggleIgnore(ctx context.Context, memberID, threadID int, ignore bool) error
	ListThreads(ctx context.Context, cursors domain.Cursors, limit, memberID int) ([]domain.Thread, domain.Cursors, error)
	GetCollapsibleThreadByID(ctx context.Context, viewable, threadID, memberID int) (*domain.Thread, error)
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
}

type MessageService interface {
	SendMessage(subject, body, memberIP string, memberID int, recipientIDs []int) (int, error)
	GetMessagesWithCursor(memberID int, reverse bool, cursor *time.Time) ([]domain.Message, error)
	GetMessageByID(messageID, memberID int) (*domain.Message, error)
	NewPost(body, ip, memberName string, messageID int) (int, error)
	GetMessagePostByID(id int) (*domain.MessagePost, error)
}

type AuthenticationService interface {
	Authenticate(ctx context.Context, username, password string) (*domain.Member, error)
}
