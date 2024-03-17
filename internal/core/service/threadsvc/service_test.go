package threadsvc

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/service/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewThreadService(t *testing.T) {
	l := zap.NewNop()
	sugar := l.Sugar()

	t.Run("successfully gets a thread by id", func(t *testing.T) {
		mockThreadRepo := &mocks.ThreadRepoMock{
			GetThreadByIDFunc: func(id int, memberID int) (*domain.Thread, error) {
				return &domain.Thread{
					ID:      1,
					Subject: "Hello, BCO",
				}, nil
			},
			ListPostsForThreadFunc: func(limit, offset, threadID, memberID int) ([]domain.ThreadPost, error) {
				return []domain.ThreadPost{
					{
						ID:   2,
						Body: "It's been a while",
					},
				}, nil
			},
		}

		mockMemberRepo := &mocks.MemberRepoMock{}

		svc := NewThreadService(mockThreadRepo, mockMemberRepo, sugar, 5)

		expectedThread := domain.Thread{
			ID:      1,
			Subject: "Hello, BCO",
			Posts: []domain.ThreadPost{
				{
					ID:       2,
					Body:     "It's been a while",
					Position: 1,
				},
			},
		}

		thread, err := svc.GetThreadByID(1, 1, 1, 1)
		require.NoError(t, err)

		assert.Len(t, mockThreadRepo.ListPostsForThreadCalls(), 1)
		assert.Len(t, mockThreadRepo.GetThreadByIDCalls(), 1)
		assert.Equal(t, 1, mockThreadRepo.ListPostsForThreadCalls()[0].ID)
		assert.Equal(t, 1, mockThreadRepo.GetThreadByIDCalls()[0].ID)
		assert.Equal(t, &expectedThread, thread)
	})

	t.Run("should bail if thread posts call returns an error", func(t *testing.T) {
		mockThreadRepo := &mocks.ThreadRepoMock{
			ListPostsForThreadFunc: func(limit, offset, threadID, memberID int) ([]domain.ThreadPost, error) {
				return nil, assert.AnError
			},
		}

		mockMemberRepo := &mocks.MemberRepoMock{}

		svc := NewThreadService(mockThreadRepo, mockMemberRepo, sugar, 5)

		thread, err := svc.GetThreadByID(1, 1, 1, 1)
		require.Error(t, err)

		assert.Nil(t, thread)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should bail if thread call returns an error", func(t *testing.T) {
		mockThreadRepo := &mocks.ThreadRepoMock{
			GetThreadByIDFunc: func(id int, memberID int) (*domain.Thread, error) {
				return nil, assert.AnError
			},
			ListPostsForThreadFunc: func(limit, offset, threadID, memberID int) ([]domain.ThreadPost, error) {
				return nil, nil
			},
		}

		mockMemberRepo := &mocks.MemberRepoMock{}

		svc := NewThreadService(mockThreadRepo, mockMemberRepo, sugar, 5)

		thread, err := svc.GetThreadByID(1, 1, 1, 1)
		require.Error(t, err)

		assert.Nil(t, thread)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("successfully saves a new post", func(t *testing.T) {
		expectedPostArg := domain.ThreadPost{
			Body:     "Hello, BCO",
			ParentID: 1,
			MemberID: 1,
			MemberIP: "127.0.0.1",
		}

		var actualPostArg domain.ThreadPost

		mockMemberRepo := &mocks.MemberRepoMock{
			GetMemberIDByUsernameFunc: func(username string) (int, error) {
				return 1, nil
			},
		}

		mockThreadRepo := &mocks.ThreadRepoMock{
			SavePostFunc: func(post domain.ThreadPost) (int, error) {
				actualPostArg = post
				return 1, nil
			},
		}

		svc := NewThreadService(mockThreadRepo, mockMemberRepo, sugar, 5)

		_, err := svc.NewPost("Hello, BCO", "127.0.0.1", "roxy", 1)
		require.NoError(t, err)

		assert.Equal(t, expectedPostArg, actualPostArg)
	})

	t.Run("successfully saves a new thread", func(t *testing.T) {
		expectedThreadArg := domain.Thread{
			MemberIP:      "127.0.0.1",
			MemberID:      1,
			Subject:       "Hello, BCO",
			FirstPostText: "Attn Roxy",
			LastPosterID:  1,
		}

		actualThreadArg := domain.Thread{}

		mockMemberRepo := &mocks.MemberRepoMock{
			GetMemberIDByUsernameFunc: func(username string) (int, error) {
				return 1, nil
			},
		}

		mockThreadRepo := &mocks.ThreadRepoMock{
			SaveThreadFunc: func(thread domain.Thread) (int, error) {
				actualThreadArg = thread
				return 1, nil
			},
		}

		svc := NewThreadService(mockThreadRepo, mockMemberRepo, sugar, 5)

		_, err := svc.NewThread("gofreescout", "127.0.0.1", "Attn Roxy", "Hello, BCO")
		require.NoError(t, err)

		assert.Equal(t, expectedThreadArg, actualThreadArg)
	})
}
