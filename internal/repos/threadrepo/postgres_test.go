package threadrepo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goBoard/db"
	"goBoard/internal/core/domain"
	"goBoard/internal/repos/seed"
	"testing"
	"time"
)

func TestNewThreadRepo(t *testing.T) {
	dbContainer, connPool, err := db.SetupTestDatabase()
	require.NoError(t, err)

	defer dbContainer.Terminate(context.Background())

	require.NoError(t, seed.SeedData(t, connPool))

	repo := NewThreadRepo(connPool, 2)

	t.Run("successfully lists all threads by member id", func(t *testing.T) {
		threads, err := repo.ListThreadsByMemberID(1, 10, 0)
		require.NoError(t, err)

		expectedThreads := []domain.Thread{
			{
				ID:             1,
				Timestamp:      nil,
				MemberID:       1,
				MemberName:     "admin",
				Subject:        "Hello, BCO",
				LastPostText:   "Attn. Roxy",
				LastPosterID:   1,
				LastPosterName: "admin",
				Views:          0,
				NumPosts:       2,
			},
		}

		threads[0].Timestamp = nil
		threads[0].DateLastPosted = nil

		assert.Equal(t, expectedThreads, threads)
	})

	t.Run("successfully lists posts by thread id", func(t *testing.T) {
		posts, err := repo.ListPostsForThread(10, 0, 1, 1)
		require.NoError(t, err)

		require.Len(t, posts, 2)
		assert.Equal(t, 2, posts[1].ID)
		assert.Equal(t, "Attn. Roxy", posts[0].Body)
		assert.Equal(t, 1, posts[0].ID)
		assert.Equal(t, "WCFRP", posts[1].Body)
	})

	t.Run("successfully gets threads by cursor forward", func(t *testing.T) {
		cursor := time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)
		janFirst := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		janSecond := time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)
		threads, err := repo.ListThreadsByCursorForward(2, &cursor, 1)
		require.NoError(t, err)

		require.Len(t, threads, 2)
		assert.Equal(t, &janSecond, threads[0].DateLastPosted)
		assert.Equal(t, &janFirst, threads[1].DateLastPosted)
	})

	t.Run("successfully gets threads by cursor in reverse", func(t *testing.T) {
		cursor := time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)
		janFourth := time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)
		janFifth := time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC)
		janSixth := time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC)
		threads, err := repo.ListThreadsInReverse(2, &cursor, 1, false, false, false)
		require.NoError(t, err)

		require.Len(t, threads, 3)
		assert.Equal(t, &janSixth, threads[0].DateLastPosted)
		assert.Equal(t, &janFifth, threads[1].DateLastPosted)
		assert.Equal(t, &janFourth, threads[2].DateLastPosted)
	})

	t.Run("successfully gets threads participated in by cursor in reverse", func(t *testing.T) {
		cursor := time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)
		janFifth := time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC)
		janFourth := time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)
		threads, err := repo.ListThreadsInReverse(2, &cursor, 1, false, false, true)
		require.NoError(t, err)

		require.Len(t, threads, 2)
		assert.Equal(t, &janFifth, threads[0].DateLastPosted)
		assert.Equal(t, &janFourth, threads[1].DateLastPosted)
	})

	t.Run("successfully gets threads favorited by cursor in reverse", func(t *testing.T) {
		cursor := time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)
		janFifth := time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC)
		janSixth := time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC)
		threads, err := repo.ListThreadsInReverse(2, &cursor, 1, false, true, false)
		require.NoError(t, err)

		require.Len(t, threads, 2)
		assert.Equal(t, &janSixth, threads[0].DateLastPosted)
		assert.Equal(t, &janFifth, threads[1].DateLastPosted)
	})

	t.Run("successfully gets threads ignored by cursor in reverse", func(t *testing.T) {
		cursor := time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)
		janFifth := time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC)
		janFourth := time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)
		threads, err := repo.ListThreadsInReverse(2, &cursor, 1, true, false, false)
		require.NoError(t, err)

		require.Len(t, threads, 2)
		assert.Equal(t, &janFifth, threads[0].DateLastPosted)
		assert.Equal(t, &janFourth, threads[1].DateLastPosted)
	})

	t.Run("successfully saves a thread", func(t *testing.T) {
		id, err := repo.SaveThread(domain.Thread{
			MemberID:      1,
			Subject:       "Hello, BCO",
			LastPosterID:  1,
			FirstPostText: "It's me Roxy",
			MemberIP:      "127.0.0.1",
		})
		require.NoError(t, err)

		var subject string
		err = connPool.QueryRow(context.Background(), "SELECT subject FROM thread WHERE id = $1", id).Scan(&subject)
		require.NoError(t, err)

		var body string
		err = connPool.QueryRow(context.Background(), "SELECT body FROM thread_post WHERE thread_id = $1", id).Scan(&body)
		require.NoError(t, err)

		assert.Equal(t, "Hello, BCO", subject)
		assert.Equal(t, "It's me Roxy", body)
	})
}

func pointerToType[T any](t T) *T {
	return &t
}
