package threadrepo

import (
	"context"
	"goBoard/db"
	"goBoard/internal/core/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const insertMember = `INSERT INTO member (name, pass, ip, email_signup, postalcode, secret) VALUES ($1, $2, $3, $4, $5, $6)`
const insertThread = `"INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ($1, $2, $3, $4)`

func TestNewThreadRepo(t *testing.T) {
	dbContainer, connPool, err := db.SetupTestDatabase()
	require.NoError(t, err)

	defer dbContainer.Terminate(context.Background())

	//require.NoError(t, seed.SeedData(t, connPool))

	_, err = connPool.Exec(context.Background(), insertMember, "admin", "admin", "127.0.0.1", "mpbyrne@gmail.com", "97217", "topsecret")
	require.NoError(t, err)

	_, err = connPool.Exec(context.Background(), insertMember, "gofreescout", "test", "127.0.0.2", "gofreescout@yahoo.com", "97217", "topsecret")
	require.NoError(t, err)

	_, err = connPool.Exec(context.Background(), `
	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('Hello, BCO', 1, 1, '2021-01-01T00:00:00Z');
	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('It stinks! A new moratorium thread', 2, 1, '2021-01-02T00:00:00Z');
	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('THYROID', 2, 1, '2021-01-03T00:00:00Z');
	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2017 politics thread', 2, 1, '2021-01-04T00:00:00Z');
	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2018 politics thread', 2, 1, '2021-01-05T00:00:00Z');
	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2019 politics thread', 2, 1, '2021-01-06T00:00:00Z');
	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2020 politics thread', 2, 1, '2021-01-07T00:00:00Z');
	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2021 politics thread', 2, 1, '2021-01-08T00:00:00Z');
	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2022 politics thread', 2, 1, '2021-01-09T00:00:00Z');
`)
	require.NoError(t, err)

	repo := NewThreadRepo(connPool, 2)

	t.Run("successfully starts at the beginning with no cursor", func(t *testing.T) {
		threads, cursorsOut, err := repo.ListThreads(context.Background(), domain.Cursors{}, 3)
		require.NoError(t, err)

		require.Len(t, threads, 3)

		assert.Equal(t, "2021-01-09T00:00:00Z", threads[0].DateLastPosted.Format(time.RFC3339Nano))
		assert.Equal(t, "2021-01-07T00:00:00Z", threads[2].DateLastPosted.Format(time.RFC3339Nano))

		assert.Empty(t, cursorsOut.Prev)
		assert.Equal(t, "2021-01-07T00:00:00Z", cursorsOut.Next)
	})

	t.Run("successfully handles a next curor", func(t *testing.T) {
		cursorsIn := domain.Cursors{
			Next: "2021-01-07T00:00:00Z",
		}

		threads, cursorsOut, err := repo.ListThreads(context.Background(), cursorsIn, 3)
		require.NoError(t, err)

		require.Len(t, threads, 3)

		assert.Equal(t, "2021-01-04T00:00:00Z", cursorsOut.Next)
		assert.Equal(t, "2021-01-06T00:00:00Z", cursorsOut.Prev)
	})

	t.Run("successfully handles forward direction on last page", func(t *testing.T) {
		cursorsIn := domain.Cursors{
			Next: "2021-01-04T00:00:00Z",
		}

		threads, cursorsOut, err := repo.ListThreads(context.Background(), cursorsIn, 3)
		require.NoError(t, err)

		require.Len(t, threads, 3)

		assert.Empty(t, cursorsOut.Next)
		assert.Equal(t, "2021-01-03T00:00:00Z", cursorsOut.Prev)
	})

	t.Run("successfully goes back in the middle", func(t *testing.T) {
		cursorsIn := domain.Cursors{
			Prev: "2021-01-02T00:00:00Z",
		}

		threads, cursorsOut, err := repo.ListThreads(context.Background(), cursorsIn, 3)
		require.NoError(t, err)

		require.Len(t, threads, 3)

		assert.Equal(t, "2021-01-03T00:00:00Z", cursorsOut.Next)
		assert.Equal(t, "2021-01-05T00:00:00Z", cursorsOut.Prev)
	})

	t.Run("successfully goes back to the beginning", func(t *testing.T) {
		cursorsIn := domain.Cursors{
			Prev: "2021-01-06T00:00:00Z",
		}

		threads, cursorsOut, err := repo.ListThreads(context.Background(), cursorsIn, 3)
		require.NoError(t, err)

		require.Len(t, threads, 3)

		assert.Empty(t, cursorsOut.Prev)
		assert.Equal(t, "2021-01-07T00:00:00Z", cursorsOut.Next)
	})

}
