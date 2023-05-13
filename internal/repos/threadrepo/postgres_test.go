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
	})

}
