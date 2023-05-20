package ports

import "goBoard/internal/core/domain"

type ThreadService interface {
	Save(post domain.Post) (int, error)
	GetPostByID(id int) (*domain.Post, error)
	GetThreadByID(limit, offset, id int) (*domain.Thread, error)
	ListThreads(limit, offset int) ([]domain.Thread, error)
	NewThread(thread domain.Thread) (int, error)
}

type MemberService interface {
	Save(member domain.Member) (int, error)
	GetMemberByID(id int) (*domain.Member, error)
}
