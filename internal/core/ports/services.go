package ports

import "goBoard/internal/core/domain"

type ThreadService interface {
	NewPost(body, ip, memberName string, threadID int) (int, error)
	GetPostByID(id int) (*domain.Post, error)
	GetThreadByID(limit, offset, id int) (*domain.Thread, error)
	ListThreads(limit, offset int) ([]domain.Thread, error)
	NewThread(memberName, memberIP, body, subject string) (int, error)
}

type MemberService interface {
	Save(member domain.Member) (int, error)
	GetMemberByID(id int) (*domain.Member, error)
	GetMemberIDByUsername(username string) (int, error)
	GetMemberByUsername(username string) (*domain.Member, error)
}
