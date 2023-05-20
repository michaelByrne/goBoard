package threadsvc

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/service/mocks"
	"testing"
)

func TestNewThreadService(t *testing.T) {
	l := zap.NewNop()
	sugar := l.Sugar()

	t.Run("successfully gets a thread by id", func(t *testing.T) {
		mockRepo := &mocks.ThreadRepoMock{
			GetThreadByIDFunc: func(id int) (*domain.Thread, error) {
				return &domain.Thread{
					ID:      1,
					Subject: "Hello, BCO",
				}, nil
			},
			GetPostsByThreadIDFunc: func(threadID int) ([]domain.Post, error) {
				return []domain.Post{
					{
						ID:   2,
						Text: "It's been a while",
					},
				}, nil
			},
		}

		svc := NewThreadService(mockRepo, sugar)

		expectedThread := domain.Thread{
			ID:      1,
			Subject: "Hello, BCO",
			Posts: []domain.Post{
				{
					ID:   2,
					Text: "It's been a while",
				},
			},
		}

		thread, err := svc.GetThreadByID(1)
		require.NoError(t, err)

		assert.Len(t, mockRepo.GetPostsByThreadIDCalls(), 1)
		assert.Len(t, mockRepo.GetThreadByIDCalls(), 1)
		assert.Equal(t, 1, mockRepo.GetPostsByThreadIDCalls()[0].ThreadID)
		assert.Equal(t, 1, mockRepo.GetThreadByIDCalls()[0].ID)
		assert.Equal(t, &expectedThread, thread)
	})

	t.Run("should bail if thread posts call returns an error", func(t *testing.T) {
		mockRepo := &mocks.ThreadRepoMock{
			GetPostsByThreadIDFunc: func(threadID int) ([]domain.Post, error) {
				return nil, assert.AnError
			},
		}

		svc := NewThreadService(mockRepo, sugar)

		thread, err := svc.GetThreadByID(1)
		require.Error(t, err)

		assert.Nil(t, thread)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should bail if thread call returns an error", func(t *testing.T) {
		mockRepo := &mocks.ThreadRepoMock{
			GetThreadByIDFunc: func(id int) (*domain.Thread, error) {
				return nil, assert.AnError
			},
			GetPostsByThreadIDFunc: func(threadID int) ([]domain.Post, error) {
				return nil, nil
			},
		}

		svc := NewThreadService(mockRepo, sugar)

		thread, err := svc.GetThreadByID(1)
		require.Error(t, err)

		assert.Nil(t, thread)
		assert.Equal(t, assert.AnError, err)
	})
}
