package postsvc

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
)

type PostService struct {
	postRepo ports.PostRepo
}

func NewPostService(postRepo ports.PostRepo) PostService {
	return PostService{postRepo}
}

func (s PostService) Save(post domain.Post) error {
	return s.postRepo.Save(post)
}

func (s PostService) GetByID(id string) (domain.Post, error) {
	return s.postRepo.GetByID(id)
}

func (s PostService) GetByThreadID(threadID string) ([]domain.Post, error) {
	return s.postRepo.GetByThreadID(threadID)
}
