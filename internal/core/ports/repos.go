package ports

import "goBoard/internal/core/domain"

type ThreadRepo interface {
	SavePost(post domain.Post) (int, error)
	GetPostByID(id string) (domain.Post, error)
	GetPostsByThreadID(threadID string) ([]domain.Post, error)
}
