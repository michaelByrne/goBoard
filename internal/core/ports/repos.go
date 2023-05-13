package ports

import "goBoard/internal/core/domain"

//go:generate moq -out thread_repo_moq.go . ThreadRepo

type ThreadRepo interface {
	SavePost(post domain.Post) (int, error)
	GetPostByID(id int) (*domain.Post, error)
	GetPostsByThreadID(threadID int) ([]domain.Post, error)
	GetThreadByID(id int) (*domain.Thread, error)
}
