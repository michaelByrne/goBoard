package threadsvc

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/service/mocks"
	"testing"
	"time"
)

func TestNewThreadService(t *testing.T) {
	l := zap.NewNop()
	sugar := l.Sugar()

	mayFirst := time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC)
	maySecond := time.Date(2020, 5, 2, 0, 0, 0, 0, time.UTC)
	mayThird := time.Date(2020, 5, 3, 0, 0, 0, 0, time.UTC)
	mayFourth := time.Date(2020, 5, 4, 0, 0, 0, 0, time.UTC)
	mayFifth := time.Date(2020, 5, 5, 0, 0, 0, 0, time.UTC)
	maySixth := time.Date(2020, 5, 6, 0, 0, 0, 0, time.UTC)
	maySeventh := time.Date(2020, 5, 7, 0, 0, 0, 0, time.UTC)

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

	t.Run("successfully gets a list of threads by page forward", func(t *testing.T) {
		mockThreadRepo := &mocks.ThreadRepoMock{
			ListThreadsByCursorForwardFunc: func(limit int, cursor *time.Time, memberID int) ([]domain.Thread, error) {
				return []domain.Thread{
					{
						ID:             1,
						DateLastPosted: &maySecond,
						Subject:        "Hello, BCO",
					},
					{
						ID:             2,
						DateLastPosted: &mayFirst,
						Subject:        "ThreadPost a picture of yourself thread",
					},
					{
						ID:             3,
						DateLastPosted: &mayFifth,
						Subject:        "Who peed in the pool",
					},
				}, nil
			},
			PeekPreviousFunc: func(timestamp *time.Time, memberID int) (bool, error) {
				return false, nil
			},
		}

		mockMemberRepo := &mocks.MemberRepoMock{}

		svc := NewThreadService(mockThreadRepo, mockMemberRepo, sugar, 2)

		site, err := svc.GetThreadsWithCursorForward(2, false, &mayThird, 1)
		require.NoError(t, err)

		assert.Len(t, site.ThreadPage.Threads, 2)
		assert.Equal(t, 1, site.ThreadPage.Threads[0].ID)
		assert.Equal(t, 2, site.ThreadPage.Threads[1].ID)
		assert.Equal(t, "Hello, BCO", site.ThreadPage.Threads[0].Subject)
		assert.Equal(t, "ThreadPost a picture of yourself thread", site.ThreadPage.Threads[1].Subject)
		assert.Equal(t, &maySecond, site.ThreadPage.Threads[0].DateLastPosted)
		assert.Equal(t, &mayFirst, site.ThreadPage.Threads[1].DateLastPosted)
		assert.Equal(t, false, site.ThreadPage.HasPrevPage)
		assert.Equal(t, &mayFirst, site.PageCursor)
		assert.Equal(t, &maySecond, site.PrevPageCursor)
	})

	t.Run("successfully gets a list of threads by page reverse", func(t *testing.T) {
		mockThreadRepo := &mocks.ThreadRepoMock{
			ListThreadsInReverseFunc: func(limit int, cursor *time.Time, memberID int, ignored, favorited, participated bool) ([]domain.Thread, error) {
				return []domain.Thread{
					{
						ID:             1,
						DateLastPosted: &maySeventh,
						Subject:        "Hello, BCO",
					},
					{
						ID:             2,
						DateLastPosted: &maySixth,
						Subject:        "ThreadPost a picture of yourself thread",
					},
					{
						ID:             3,
						DateLastPosted: &mayFifth,
						Subject:        "Who peed in the pool",
					},
					{
						ID:             4,
						DateLastPosted: &mayFourth,
						Subject:        "soup's on!",
					},
					{
						ID:             5,
						DateLastPosted: &mayThird,
						Subject:        "I'm in outeep space",
					},
					{
						ID:             6,
						DateLastPosted: &maySecond,
						Subject:        "new experiences thread",
					},
					{
						ID:             7,
						DateLastPosted: &mayFirst,
						Subject:        "thread for eataly",
					},
				}, nil
			},
			PeekPreviousFunc: func(timestamp *time.Time, memberID int) (bool, error) {
				return true, nil
			},
		}

		mockMemberRepo := &mocks.MemberRepoMock{}

		svc := NewThreadService(mockThreadRepo, mockMemberRepo, sugar, 2)

		site, err := svc.GetThreadsWithCursorReverse(3, &mayFourth, 1, false, false, false)
		require.NoError(t, err)

		assert.Len(t, site.ThreadPage.Threads, 2)
		assert.Equal(t, &maySixth, site.PageCursor)
		assert.True(t, site.ThreadPage.HasPrevPage)
		assert.Equal(t, &maySeventh, site.PrevPageCursor)
	})
}
