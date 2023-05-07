package ports

import "goBoard/internal/core/domain"

type PostRepo interface {
	Save(post domain.Post) error
	GetByID(id string) (domain.Post, error)
	GetByThreadID(threadID string) ([]domain.Post, error)
}
