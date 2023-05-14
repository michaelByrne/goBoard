package ports

import "goBoard/internal/core/domain"

type ThreadService interface {
	Save(post domain.Post) error
	GetPostByID(id int) (*domain.Post, error)
	GetPostsByThreadID(threadID int) ([]domain.Post, error)
	GetThreadByID(id int) (*domain.Thread, error)
	ListThreads(limit int) ([]domain.Thread, error)
	NewThread(thread domain.Thread) error
}

type MemberService interface {
	Save(member domain.Member) error
	GetMemberByID(id int) (*domain.Member, error)
}
