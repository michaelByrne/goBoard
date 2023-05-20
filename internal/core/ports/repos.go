package ports

import "goBoard/internal/core/domain"

//go:generate moq -pkg mocks -out ../service/mocks/thread_repo_moq.go . ThreadRepo

type ThreadRepo interface {
	SavePost(post domain.Post) (int, error)
	GetPostByID(id int) (*domain.Post, error)
	GetThreadByID(id int) (*domain.Thread, error)
	ListThreads(limit, offset int) ([]domain.Thread, error)
	ListThreadsByMemberID(memberID int, limit, offset int) ([]domain.Thread, error)
	SaveThread(thread domain.Thread) (int, error)
	ListPostsForThread(limit, offset, id int) ([]domain.Post, error)
}

type MemberRepo interface {
	SaveMember(member domain.Member) (int, error)
	GetMemberByID(id int) (*domain.Member, error)
}
