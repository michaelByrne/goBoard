package ports

import "goBoard/internal/core/domain"

//go:generate moq -pkg mocks -out ../service/mocks/thread_repo_moq.go . ThreadRepo

type ThreadRepo interface {
	SavePost(post domain.Post) (int, error)
	GetPostByID(id int) (*domain.Post, error)
	GetPostsByThreadID(threadID int) ([]domain.Post, error)
	GetThreadByID(id int) (*domain.Thread, error)
	ListThreads(limit int) ([]domain.Thread, error)
	SaveThread(thread domain.Thread) (int, error)
}

type MemberRepo interface {
	SaveMember(member domain.Member) (int, error)
	GetMemberByID(id int) (*domain.Member, error)
}
