package messagesvc

import (
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"time"
)

type MessageService struct {
	logger              *zap.SugaredLogger
	messageRepo         ports.MessageRepo
	memberRepo          ports.MemberRepo
	maxMessageViewLimit int
}

func NewMessageService(messageRepo ports.MessageRepo, memberRepo ports.MemberRepo, logger *zap.SugaredLogger, maxMessageViewLimit int) MessageService {
	return MessageService{
		logger:              logger,
		messageRepo:         messageRepo,
		memberRepo:          memberRepo,
		maxMessageViewLimit: maxMessageViewLimit,
	}
}

func (s MessageService) SendMessage(subject, body, memberIP string, memberID int, recipientIDs []int) (int, error) {
	s.logger.Infof("sending message with subject: %s, body: %s, to members: %v from member: %d", subject, body, recipientIDs, memberID)

	message := domain.Message{
		Subject:      subject,
		Body:         body,
		MemberID:     memberID,
		MemberIP:     memberIP,
		RecipientIDs: recipientIDs,
	}

	messageID, err := s.messageRepo.SaveMessage(message)
	if err != nil {
		s.logger.Errorf("error sending message: %s", err)
		return 0, err
	}

	return messageID, nil
}

func (s MessageService) GetMessagesWithCursor(memberID int, reverse bool, cursor *time.Time) ([]domain.Message, error) {
	s.logger.Infof("getting messages for member: %d", memberID)

	if !reverse {
		messages, err := s.messageRepo.GetMessagesWithCursorForward(memberID, s.maxMessageViewLimit, cursor)
		if err != nil {
			s.logger.Errorf("error getting messages for member: %d", memberID)
			return nil, err
		}

		if len(messages) > s.maxMessageViewLimit {
			messages = messages[:s.maxMessageViewLimit]
			messages[0].HasNextPage = true
			messages[0].PageCursor = messages[s.maxMessageViewLimit-1].DateLastPosted
		}

		hasPrevious, err := s.messageRepo.PeekPrevious(cursor)
		if err != nil {
			s.logger.Errorf("error peeking previous: %s", err)
			return nil, err
		}

		if hasPrevious {
			messages[0].HasPrevPage = true
			messages[0].PrevPageCursor = messages[0].DateLastPosted
		}

		return messages, nil
	}

	messages, err := s.messageRepo.GetMessagesWithCursorReverse(memberID, s.maxMessageViewLimit, cursor)
	if err != nil {
		s.logger.Errorf("error getting messages for member: %d", memberID)
		return nil, err
	}

	if len(messages) > s.maxMessageViewLimit {
		messages = messages[:s.maxMessageViewLimit]
		messages[0].HasNextPage = true
		messages[0].PageCursor = messages[s.maxMessageViewLimit-1].DateLastPosted
	}

	hasPrevious, err := s.messageRepo.PeekPrevious(messages[0].DateLastPosted)
	if err != nil {
		s.logger.Errorf("error peeking previous: %s", err)
		return nil, err
	}

	if hasPrevious {
		messages[0].HasPrevPage = true
		messages[0].PrevPageCursor = messages[0].DateLastPosted
	}

	return messages, nil
}

func (s MessageService) GetMessageByID(messageID, memberID int) (*domain.Message, error) {
	s.logger.Infof("getting message with id: %d", messageID)

	posts, err := s.messageRepo.GetMessagePostsByID(memberID, messageID, 30)
	if err != nil {
		s.logger.Errorf("error getting posts with message id: %d", messageID)
		return nil, err
	}

	s.logger.Infof("got %d posts for message id %d", len(posts), messageID)

	message := &domain.Message{
		ID:      posts[0].ParentID,
		Subject: posts[0].ParentSubject,
	}

	for idx, post := range posts {
		postPtr := &post
		postPtr.Position = idx + 1
		message.Posts = append(message.Posts, *postPtr)
	}

	return message, nil
}

func (s MessageService) NewPost(body, memberIP, memberName string, messageID int) (int, error) {
	s.logger.Infof("creating new post with message id %d", messageID)

	memberID, err := s.memberRepo.GetMemberIDByUsername(memberName)
	if err != nil {
		s.logger.Errorf("error getting member id for member name: %s", memberName)
		return 0, err
	}

	post := domain.MessagePost{
		Body:       body,
		MemberIP:   memberIP,
		MemberName: memberName,
		ParentID:   messageID,
		MemberID:   memberID,
	}

	postID, err := s.messageRepo.SavePost(post)
	if err != nil {
		s.logger.Errorf("error creating new post: %s\n", err)
		return 0, err
	}

	s.logger.Infof("created new post with id: %d", postID)

	return postID, nil
}

func (s MessageService) GetMessagePostByID(postID int) (*domain.MessagePost, error) {
	s.logger.Infof("getting post with id: %d", postID)

	post, err := s.messageRepo.GetMessagePostByID(postID)
	if err != nil {
		s.logger.Errorf("error getting post with id: %d", postID)
		return nil, err
	}

	return post, nil
}
