package ports

import "goBoard/internal/core/domain"

type ThreadRepo interface {
	SavePost(post domain.Post) (int, error)
	GetPostByID(id int) (*domain.Post, error)
	GetPostsByThreadID(threadID int) ([]domain.Post, error)
}
