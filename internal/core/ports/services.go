package ports

import (
	"goBoard/internal/core/domain"
	"html/template"
	"time"
)

type ThreadService interface {
	NewPost(body, ip, memberName string, threadID int) (int, error)
	GetPostByID(id int) (*domain.ThreadPost, error)
	GetThreadByID(limit, offset, id int) (*domain.Thread, error)
	ListThreads(limit, offset int) (*domain.SiteContext, error)
	NewThread(memberName, memberIP, body, subject string) (int, error)
	GetThreadsWithCursorForward(limit int, firstPage bool, cursor *time.Time) (*domain.SiteContext, error)
	GetThreadsWithCursorReverse(limit int, cursor *time.Time) (*domain.SiteContext, error)
	ConvertPostBodyBbcodeToHtml(body string) (*template.HTML, error)
}

type MemberService interface {
	Save(member domain.Member) (int, error)
	GetMemberByID(id int) (*domain.Member, error)
	GetMemberIDByUsername(username string) (int, error)
	GetMemberByUsername(username string) (*domain.SiteContext, error)
	ValidateMembers(names []string) ([]domain.Member, error)
}

type MessageService interface {
	SendMessage(subject, body, memberIP string, memberID int, recipientIDs []int) (int, error)
	GetMessagesWithCursor(memberID int, reverse bool, cursor *time.Time) ([]domain.Message, error)
	GetMessageByID(messageID, memberID int) (*domain.Message, error)
	NewPost(body, ip, memberName string, messageID int) (int, error)
	GetMessagePostByID(id int) (*domain.MessagePost, error)
}
