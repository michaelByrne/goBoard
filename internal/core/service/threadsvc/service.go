package threadsvc

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
)

type ThreadService struct {
	threadRepo ports.ThreadRepo
}

func NewThreadService(postRepo ports.ThreadRepo) ThreadService {
	return ThreadService{postRepo}
}

func (s ThreadService) Save(post domain.Post) error {
	return s.threadRepo.SavePost(post)
}

func (s ThreadService) GetByID(id string) (domain.Post, error) {
	return s.threadRepo.GetPostByID(id)
}

func (s ThreadService) GetPostsByThreadID(threadID string) ([]domain.Post, error) {
	return s.threadRepo.GetPostsByThreadID(threadID)
}
