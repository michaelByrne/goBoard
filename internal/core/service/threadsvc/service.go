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
	_, err := s.threadRepo.SavePost(post)
	if err != nil {
		return err
	}

	return nil
}

func (s ThreadService) GetByID(id int) (*domain.Post, error) {
	return s.threadRepo.GetPostByID(id)
}

func (s ThreadService) GetPostsByThreadID(threadID int) ([]domain.Post, error) {
	return s.threadRepo.GetPostsByThreadID(threadID)
}
