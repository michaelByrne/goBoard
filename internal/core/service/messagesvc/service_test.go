package messagesvc

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/service/mocks"
	"testing"
	"time"
)

func TestMessageService_GetMessagesWithCursor(t *testing.T) {
	mayFirst := time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC)
	maySecond := time.Date(2020, 5, 2, 0, 0, 0, 0, time.UTC)
	//mayThird := time.Date(2020, 5, 3, 0, 0, 0, 0, time.UTC)
	mayFourth := time.Date(2020, 5, 4, 0, 0, 0, 0, time.UTC)
	//mayFifth := time.Date(2020, 5, 5, 0, 0, 0, 0, time.UTC)
	maySixth := time.Date(2020, 5, 6, 0, 0, 0, 0, time.UTC)
	maySeventh := time.Date(2020, 5, 7, 0, 0, 0, 0, time.UTC)

	l := zap.NewNop()
	sugar := l.Sugar()

	t.Run("should return messages with forward cursor", func(t *testing.T) {
		var peekTimestamp *time.Time
		var cursorCalledWith *time.Time
		mockMessageRepo := &mocks.MessageRepoMock{
			GetMessagesWithCursorForwardFunc: func(memberID int, limit int, cursor *time.Time) ([]domain.Message, error) {
				cursorCalledWith = cursor
				return []domain.Message{
					{
						ID:             1,
						Subject:        "subject",
						Body:           "body",
						MemberID:       1,
						DateLastPosted: &mayFirst,
					},
					{
						ID:             2,
						Subject:        "Luna the dog",
						Body:           "Luna is a good dog",
						MemberID:       2,
						DateLastPosted: &maySecond,
					},
					{
						ID:             3,
						Subject:        "TDD",
						Body:           "TDD is the best",
						DateLastPosted: &maySixth,
					},
					{
						ID:             4,
						Subject:        "Authentication",
						Body:           "Authentication is important",
						DateLastPosted: &maySeventh,
					},
				}, nil
			},
			PeekPreviousFunc: func(timestamp *time.Time) (bool, error) {
				peekTimestamp = timestamp
				return false, nil
			},
		}

		mockMemberRepo := &mocks.MemberRepoMock{}

		svc := NewMessageService(mockMessageRepo, mockMemberRepo, sugar, 3)

		messages, err := svc.GetMessagesWithCursor(1, false, &mayFourth)
		require.NoError(t, err)

		require.Equal(t, 3, len(messages))
		assert.Equal(t, &mayFourth, peekTimestamp)
		assert.Equal(t, &mayFourth, cursorCalledWith)
		assert.True(t, messages[0].HasNextPage)
		assert.False(t, messages[0].HasPrevPage)
		assert.Equal(t, &maySixth, messages[0].PageCursor)
	})

	t.Run("should return messages with backward cursor", func(t *testing.T) {
		var peekTimestamp *time.Time
		var cursorCalledWith *time.Time
		mockMessageRepo := &mocks.MessageRepoMock{
			GetMessagesWithCursorReverseFunc: func(memberID int, limit int, cursor *time.Time) ([]domain.Message, error) {
				cursorCalledWith = cursor
				return []domain.Message{
					{
						ID:             1,
						Subject:        "subject",
						Body:           "body",
						MemberID:       1,
						DateLastPosted: &mayFirst,
					},
					{
						ID:             2,
						Subject:        "Luna the dog",
						Body:           "Luna is a good dog",
						MemberID:       2,
						DateLastPosted: &maySecond,
					},
					{
						ID:             3,
						Subject:        "TDD",
						Body:           "TDD is the best",
						DateLastPosted: &maySixth,
					},
					{
						ID:             4,
						Subject:        "Authentication",
						Body:           "Authentication is important",
						DateLastPosted: &maySeventh,
					},
				}, nil
			},
			PeekPreviousFunc: func(timestamp *time.Time) (bool, error) {
				peekTimestamp = timestamp
				return true, nil
			},
		}

		mockMemberRepo := &mocks.MemberRepoMock{}

		svc := NewMessageService(mockMessageRepo, mockMemberRepo, sugar, 3)

		messages, err := svc.GetMessagesWithCursor(1, true, &mayFourth)
		require.NoError(t, err)

		require.Equal(t, 3, len(messages))
		assert.Equal(t, &mayFirst, peekTimestamp)
		assert.Equal(t, &mayFourth, cursorCalledWith)
		assert.True(t, messages[0].HasPrevPage)
		assert.True(t, messages[0].HasNextPage)
		assert.Equal(t, &maySixth, messages[0].PageCursor)
		assert.Equal(t, &mayFirst, messages[0].PrevPageCursor)
	})
}

func TestMessageService_GetMessageByID(t *testing.T) {
	mockMessageRepo := &mocks.MessageRepoMock{
		GetMessagePostsByIDFunc: func(memberID int, messageID int, limit int) ([]domain.MessagePost, error) {
			return []domain.MessagePost{
				{
					ID:            1,
					ParentID:      2,
					ParentSubject: "subject",
				},
				{
					ID:            2,
					ParentID:      2,
					ParentSubject: "testing!",
				},
				{
					ID:            3,
					ParentID:      2,
					ParentSubject: "luna the dog",
				},
			}, nil
		},
	}

	mockMemberRepo := &mocks.MemberRepoMock{}

	svc := NewMessageService(mockMessageRepo, mockMemberRepo, zap.NewNop().Sugar(), 3)

	message, err := svc.GetMessageByID(1, 2)
	require.NoError(t, err)

	assert.Equal(t, 2, message.ID)
	assert.Equal(t, "subject", message.Subject)
	assert.Equal(t, 3, len(message.Posts))
	assert.Equal(t, 1, message.Posts[0].Position)
	assert.Equal(t, 2, message.Posts[1].Position)
	assert.Equal(t, 3, message.Posts[2].Position)
}

func TestMessageService_NewPost(t *testing.T) {
	t.Run("should create new post", func(t *testing.T) {
		postToSave := domain.MessagePost{}
		usernameToGet := ""
		mockMessageRepo := &mocks.MessageRepoMock{
			SavePostFunc: func(post domain.MessagePost) (int, error) {
				postToSave = post
				return 1, nil
			},
		}

		mockMemberRepo := &mocks.MemberRepoMock{
			GetMemberIDByUsernameFunc: func(username string) (int, error) {
				usernameToGet = username
				return 666, nil
			},
		}

		svc := NewMessageService(mockMessageRepo, mockMemberRepo, zap.NewNop().Sugar(), 3)

		postID, err := svc.NewPost("great body", "127.0.0.1", "username", 1)
		require.NoError(t, err)

		assert.Equal(t, 1, postID)

		expectedPostToSave := domain.MessagePost{
			Body:       "great body",
			MemberIP:   "127.0.0.1",
			MemberName: "username",
			ParentID:   1,
			MemberID:   666,
		}

		assert.Equal(t, expectedPostToSave, postToSave)
		assert.Equal(t, "username", usernameToGet)
	})
}
