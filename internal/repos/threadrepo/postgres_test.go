package threadrepo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goBoard/db"
	"goBoard/internal/core/domain"
	"goBoard/internal/repos/seed"
	"testing"
)

func TestNewThreadRepo(t *testing.T) {
	dbContainer, connPool, err := db.SetupTestDatabase()
	require.NoError(t, err)

	defer dbContainer.Terminate(context.Background())

	require.NoError(t, seed.SeedData(t, connPool))

	repo := NewThreadRepo(connPool)

	t.Run("successfully adds a post", func(t *testing.T) {
		id, err := repo.SavePost(domain.Post{
			ThreadID: 1,
			MemberID: 1,
			MemberIP: "67.189.58.94",
			Text:     "Hello, BCO",
		})
		require.NoError(t, err)

		var body string
		err = connPool.QueryRow(context.Background(), "SELECT body FROM thread_post WHERE id = $1", id).Scan(&body)
		require.NoError(t, err)

		assert.Equal(t, "Hello, BCO", body)

		_, err = connPool.Exec(context.Background(), "DELETE FROM thread_post WHERE id = $1", id)
		require.NoError(t, err)
	})

	t.Run("fails to add post if the thread doesn't exist", func(t *testing.T) {
		_, err := repo.SavePost(domain.Post{
			ThreadID: 666,
			MemberID: 1,
			MemberIP: "67.189.58.94",
			Text:     "Hello, BCO",
		})
		require.Error(t, err)
	})

	t.Run("successfully gets a post by id", func(t *testing.T) {
		post, err := repo.GetPostByID(1)
		require.NoError(t, err)

		expectedPost := domain.Post{
			ID:        1,
			Timestamp: nil,
			MemberID:  1,
			MemberIP:  "127.0.0.1/32",
			ThreadID:  1,
			Text:      "Attn. Roxy",
		}

		post.Timestamp = nil

		assert.Equal(t, &expectedPost, post)
	})

	t.Run("successfully gets posts by thread id", func(t *testing.T) {
		posts, err := repo.GetPostsByThreadID(1)
		require.NoError(t, err)
		require.Len(t, posts, 2)

		expectedPosts := []domain.Post{
			{
				ID:        1,
				Timestamp: nil,
				MemberID:  1,
				MemberIP:  "127.0.0.1/32",
				ThreadID:  1,
				Text:      "Attn. Roxy",
			},
			{
				ID:        2,
				Timestamp: nil,
				MemberID:  1,
				MemberIP:  "127.0.0.2/32",
				ThreadID:  1,
				Text:      "WCFRP",
			},
		}

		posts[0].Timestamp = nil
		posts[1].Timestamp = nil

		assert.Equal(t, expectedPosts, posts)
	})

	t.Run("successfully gets a thread by id", func(t *testing.T) {
		thread, err := repo.GetThreadByID(1)
		require.NoError(t, err)

		expectedThread := domain.Thread{
			ID:        1,
			Timestamp: nil,
			MemberID:  1,
			Subject:   "Hello, BCO",
			Views:     0,
		}

		thread.Timestamp = nil

		assert.Equal(t, &expectedThread, thread)
	})

	t.Run("successfully lists all threads", func(t *testing.T) {
		threads, err := repo.ListThreads(10)
		require.NoError(t, err)

		expectedThreads := []domain.Thread{
			{
				ID:             2,
				Timestamp:      nil,
				MemberID:       2,
				MemberName:     "gofreescout",
				Views:          0,
				LastPosterName: "gofreescout",
				LastPosterID:   2,
				LastPostText:   "I listened to a podcast earlier that had five minutes of ads at the beginning",
				NumPosts:       3,
				Subject:        "It stinks! A new moratorium thread",
			},
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
		threads[0].DateLastPosted = ""
		threads[1].Timestamp = nil
		threads[1].DateLastPosted = ""

		assert.Equal(t, expectedThreads, threads)
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
