package threadsvc

import (
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
)

type ThreadService struct {
	threadRepo ports.ThreadRepo
	memberRepo ports.MemberRepo
	logger     *zap.SugaredLogger
}

func NewThreadService(postRepo ports.ThreadRepo, memberREpo ports.MemberRepo, logger *zap.SugaredLogger) ThreadService {
	return ThreadService{
		threadRepo: postRepo,
		logger:     logger,
		memberRepo: memberREpo,
	}
}

func (s ThreadService) NewPost(body, ip, memberName string, threadID int) (int, error) {
	memberID, err := s.memberRepo.GetMemberIDByUsername(memberName)
	if err != nil {
		s.logger.Errorf("error getting member id by username: %v", err)
		return 0, err
	}

	id, err := s.threadRepo.SavePost(domain.Post{
		Text:     body,
		MemberIP: ip,
		ThreadID: threadID,
		MemberID: memberID,
	})
	if err != nil {
		s.logger.Errorf("error saving post: %v", err)
		return 0, err
	}

	return id, nil
}

func (s ThreadService) GetPostByID(id int) (*domain.Post, error) {
	return s.threadRepo.GetPostByID(id)
}

func (s ThreadService) GetThreadByID(limit, offset, id int) (*domain.Thread, error) {
	posts, err := s.threadRepo.ListPostsForThread(limit, offset, id)
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
