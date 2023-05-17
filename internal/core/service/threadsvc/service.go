package threadsvc

import (
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
)

type ThreadService struct {
	threadRepo ports.ThreadRepo
	logger     *zap.SugaredLogger
}

func NewThreadService(postRepo ports.ThreadRepo, logger *zap.SugaredLogger) ThreadService {
	return ThreadService{
		threadRepo: postRepo,
		logger:     logger,
	}
}

func (s ThreadService) Save(post domain.Post) (int, error) {
	id, err := s.threadRepo.SavePost(post)
	if err != nil {
		s.logger.Errorf("error saving post: %v", err)
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
		s.logger.Errorf("error getting posts by thread id: %v", err)
		return nil, err
	}

	thread, err := s.threadRepo.GetThreadByID(id)
	if err != nil {
		s.logger.Errorf("error getting thread by id: %v", err)
		return nil, err
	}

	thread.Posts = posts

	return thread, nil
}

func (s ThreadService) ListThreads(limit, offset int) ([]domain.Thread, error) {
	return s.threadRepo.ListThreads(limit, offset)
}

func (s ThreadService) NewThread(thread domain.Thread) (int, error) {
	id, err := s.threadRepo.SaveThread(thread)
	if err != nil {
		s.logger.Errorf("error saving thread: %v", err)
		return 0, err
	}

	return id, nil
}
