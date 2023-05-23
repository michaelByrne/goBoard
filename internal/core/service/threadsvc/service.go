package threadsvc

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"

	"go.uber.org/zap"
)

type ThreadService struct {
	threadRepo ports.ThreadRepo
	memberRepo ports.MemberRepo
	logger     *zap.SugaredLogger
}

func NewThreadService(postRepo ports.ThreadRepo, memberRepo ports.MemberRepo, logger *zap.SugaredLogger) ThreadService {
	return ThreadService{
		threadRepo: postRepo,
		logger:     logger,
		memberRepo: memberRepo,
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

	for idx, post := range posts {
		postPtr := &post
		postPtr.ThreadPosition = idx + 1
		thread.Posts = append(thread.Posts, *postPtr)
	}

	return thread, nil
}

func (s ThreadService) ListThreads(limit, offset int) (domain.ThreadPage, error) {
	return s.threadRepo.ListThreads(limit, offset)
}

func (s ThreadService) NewThread(memberName, memberIP, body, subject string) (int, error) {
	id, err := s.memberRepo.GetMemberIDByUsername(memberName)
	if err != nil {
		s.logger.Errorf("error getting member id by username: %v", err)
		return 0, err
	}

	thread := domain.Thread{
		Subject:       subject,
		FirstPostText: body,
		MemberID:      id,
		LastPosterID:  id,
		MemberIP:      memberIP,
	}

	threadID, err := s.threadRepo.SaveThread(thread)
	if err != nil {
		s.logger.Errorf("error saving thread: %v", err)
		return 0, err
	}

	return threadID, nil
}
