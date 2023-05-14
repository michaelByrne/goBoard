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

func (s ThreadService) Save(post domain.Post) (int, error) {
	id, err := s.threadRepo.SavePost(post)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s ThreadService) GetPostByID(id int) (*domain.Post, error) {
	return s.threadRepo.GetPostByID(id)
}

func (s ThreadService) GetPostsByThreadID(threadID int) ([]domain.Post, error) {
	return s.threadRepo.GetPostsByThreadID(threadID)
}

func (s ThreadService) GetThreadByID(id int) (*domain.Thread, error) {
	posts, err := s.threadRepo.GetPostsByThreadID(id)
	if err != nil {
		return nil, err
	}

	thread, err := s.threadRepo.GetThreadByID(id)
	if err != nil {
		return nil, err
	}

	thread.Posts = posts

	return thread, nil
}

func (s ThreadService) ListThreads(limit int) ([]domain.Thread, error) {
	return s.threadRepo.ListThreads(limit)
}

func (s ThreadService) NewThread(thread domain.Thread) (int, error) {
	id, err := s.threadRepo.SaveThread(thread)
	if err != nil {
		return 0, err
	}

	return id, nil
}
