package threadrepo

import (
	"context"
	"goBoard/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const insertMember = `INSERT INTO member (name, pass, ip, email_signup, postalcode, secret) VALUES ($1, $2, $3, $4, $5, $6)`

// func TestNewThreadRepo_ListThreads(t *testing.T) {
// 	dbContainer, connPool, err := db.SetupTestDatabase()
// 	require.NoError(t, err)

// 	defer dbContainer.Terminate(context.Background())

// 	//require.NoError(t, seed.SeedData(t, connPool))

// 	_, err = connPool.Exec(context.Background(), insertMember, "admin", "admin", "127.0.0.1", "mpbyrne@gmail.com", "97217", "topsecret")
// 	require.NoError(t, err)

// 	_, err = connPool.Exec(context.Background(), insertMember, "gofreescout", "test", "127.0.0.2", "gofreescout@yahoo.com", "97217", "topsecret")
// 	require.NoError(t, err)

// 	_, err = connPool.Exec(context.Background(), `
// 	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('Hello, BCO', 1, 1, '2021-01-01T00:00:00Z');
// 	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('It stinks! A new moratorium thread', 2, 1, '2021-01-02T00:00:00Z');
// 	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('THYROID', 2, 1, '2021-01-03T00:00:00Z');
// 	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2017 politics thread', 2, 1, '2021-01-04T00:00:00Z');
// 	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2018 politics thread', 2, 1, '2021-01-05T00:00:00Z');
// 	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2019 politics thread', 2, 1, '2021-01-06T00:00:00Z');
// 	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2020 politics thread', 2, 1, '2021-01-07T00:00:00Z');
// 	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2021 politics thread', 2, 1, '2021-01-08T00:00:00Z');
// 	INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('2022 politics thread', 2, 1, '2021-01-09T00:00:00Z');
// `)
// 	require.NoError(t, err)

// 	_, err = connPool.Exec(context.Background(), `
// 	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 1, 'Hello, BCO!', '2021-01-01T00:00:00Z', '172.0.0.1');
// 	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (2, 2, 'It stinks!', '2021-01-02T00:00:00Z', '172.0.0.1');
// 	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (3, 2, 'THYROID!', '2021-01-03T00:00:00Z', '172.0.0.1');
// 	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (4, 2, '2017 politics!', '2021-01-04T00:00:00Z', '127.0.0.1');
// 	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (5, 2, '2018 politics!', '2021-01-05T00:00:00Z', '127.0.0.1');
// 	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (6, 2, '2019 politics!', '2021-01-06T00:00:00Z', '127.0.0.1');
// 	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (7, 2, '2020 politics!', '2021-01-07T00:00:00Z', '127.0.0.1');
// 	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (8, 2, '2021 politics!', '2021-01-08T00:00:00Z', '127.0.0.1');
// 	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (9, 2, '2022 politics!', '2021-01-09T00:00:00Z', '127.0.0.1');
// 	`)
// 	require.NoError(t, err)

// 	repo := NewThreadRepo(connPool, 2)

// 	t.Run("successfully starts at the beginning with no cursor", func(t *testing.T) {
// 		threads, cursorsOut, err := repo.ListThreads(context.Background(), domain.Cursors{}, 3, 1)
// 		require.NoError(t, err)

// 		require.Len(t, threads, 3)

// 		assert.Equal(t, "2021-01-09T00:00:00Z", threads[0].DateLastPosted.Format(time.RFC3339Nano))
// 		assert.Equal(t, "2021-01-07T00:00:00Z", threads[2].DateLastPosted.Format(time.RFC3339Nano))

// 		assert.Empty(t, cursorsOut.Prev)
// 		assert.Equal(t, "2021-01-07T00:00:00Z", cursorsOut.Next)
// 	})

// 	t.Run("successfully handles a next cursor", func(t *testing.T) {
// 		cursorsIn := domain.Cursors{
// 			Next: "2021-01-07T00:00:00Z",
// 		}

// 		threads, cursorsOut, err := repo.ListThreads(context.Background(), cursorsIn, 3, 1)
// 		require.NoError(t, err)

// 		require.Len(t, threads, 3)

// 		assert.Equal(t, "2021-01-06T00:00:00Z", threads[0].DateLastPosted.Format(time.RFC3339Nano))
// 		assert.Equal(t, "2021-01-04T00:00:00Z", threads[2].DateLastPosted.Format(time.RFC3339Nano))

// 		assert.Equal(t, "2021-01-04T00:00:00Z", cursorsOut.Next)
// 		assert.Equal(t, "2021-01-06T00:00:00Z", cursorsOut.Prev)
// 	})

// 	t.Run("successfully handles forward direction on last page", func(t *testing.T) {
// 		cursorsIn := domain.Cursors{
// 			Next: "2021-01-04T00:00:00Z",
// 		}

// 		threads, cursorsOut, err := repo.ListThreads(context.Background(), cursorsIn, 3, 1)
// 		require.NoError(t, err)

// 		require.Len(t, threads, 3)

// 		assert.Equal(t, "2021-01-03T00:00:00Z", threads[0].DateLastPosted.Format(time.RFC3339Nano))
// 		assert.Equal(t, "2021-01-01T00:00:00Z", threads[2].DateLastPosted.Format(time.RFC3339Nano))

// 		assert.Empty(t, cursorsOut.Next)
// 		assert.Equal(t, "2021-01-03T00:00:00Z", cursorsOut.Prev)
// 	})

// 	t.Run("successfully goes back in the middle", func(t *testing.T) {
// 		cursorsIn := domain.Cursors{
// 			Prev: "2021-01-02T00:00:00Z",
// 		}

// 		threads, cursorsOut, err := repo.ListThreads(context.Background(), cursorsIn, 3, 1)
// 		require.NoError(t, err)

// 		require.Len(t, threads, 3)

// 		assert.Equal(t, "2021-01-05T00:00:00Z", threads[0].DateLastPosted.Format(time.RFC3339Nano))
// 		assert.Equal(t, "2021-01-03T00:00:00Z", threads[2].DateLastPosted.Format(time.RFC3339Nano))

// 		assert.Equal(t, "2021-01-03T00:00:00Z", cursorsOut.Next)
// 		assert.Equal(t, "2021-01-05T00:00:00Z", cursorsOut.Prev)
// 	})

// 	t.Run("successfully goes back to the beginning", func(t *testing.T) {
// 		cursorsIn := domain.Cursors{
// 			Prev: "2021-01-06T00:00:00Z",
// 		}

// 		threads, cursorsOut, err := repo.ListThreads(context.Background(), cursorsIn, 3, 1)
// 		require.NoError(t, err)

// 		require.Len(t, threads, 3)

// 		assert.Equal(t, "2021-01-09T00:00:00Z", threads[0].DateLastPosted.Format(time.RFC3339Nano))
// 		assert.Equal(t, "2021-01-07T00:00:00Z", threads[2].DateLastPosted.Format(time.RFC3339Nano))

// 		assert.Empty(t, cursorsOut.Prev)
// 		assert.Equal(t, "2021-01-07T00:00:00Z", cursorsOut.Next)
// 	})

// 	t.Run("successfully handles a last page with four items in forward direction", func(t *testing.T) {
// 		cursorsIn := domain.Cursors{
// 			Next: "2021-01-02T00:00:00Z",
// 		}

// 		threads, cursorsOut, err := repo.ListThreads(context.Background(), cursorsIn, 4, 1)
// 		require.NoError(t, err)

// 		require.Len(t, threads, 1)

// 		assert.Equal(t, "2021-01-01T00:00:00Z", threads[0].DateLastPosted.Format(time.RFC3339Nano))

// 		assert.Empty(t, cursorsOut.Next)
// 		assert.Equal(t, "2021-01-01T00:00:00Z", cursorsOut.Prev)
// 	})
// }

func TestNewThreadRepo_ListPostsCollapsible(t *testing.T) {
	dbContainer, connPool, err := db.SetupTestDatabase()
	require.NoError(t, err)

	defer dbContainer.Terminate(context.Background())

	_, err = connPool.Exec(context.Background(), insertMember, "admin", "admin", "127.0.0.1", "mpbyrne@gmail.com", "97217", "topsecret")
	require.NoError(t, err)

	_, err = connPool.Exec(context.Background(), insertMember, "gofreescout", "test", "127.0.0.2", "gofreescout@yahoo.com", "97217", "topsecret")
	require.NoError(t, err)

	_, err = connPool.Exec(context.Background(), `INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ('Hello, BCO', 1, 1, '2021-01-01T00:00:00Z')`)
	require.NoError(t, err)

	_, err = connPool.Exec(context.Background(), `
	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 1, 'Hello, BCO!', '2021-01-01T00:00:00Z', '127.0.0.1');
	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 2, 'Welcome back!', '2021-01-02T00:00:00Z', '127.0.0.1');
	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 1, 'I am back!', '2021-01-03T00:00:00Z', '127.0.0.1');
	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 2, 'I am back too!', '2021-01-04T00:00:00Z', '127.0.0.1');
	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 1, 'Wow!', '2021-01-05T00:00:00Z', '127.0.0.1');
	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 2, 'Where have you been??', '2021-01-06T00:00:00Z', '127.0.0.1');
	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 1, 'I was on vacation!', '2021-01-07T00:00:00Z', '127.0.0.1');
	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 2, 'I was too!', '2021-01-08T00:00:00Z', '127.0.0.1');
	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 1, 'I was in Hawaii!', '2021-01-09T00:00:00Z', '127.0.0.1');
	INSERT INTO thread_post (thread_id, member_id, body, date_posted, member_ip) VALUES (1, 2, 'I was in Amish Country!', '2021-01-10T00:00:00Z', '127.0.0.1');
	`)
	require.NoError(t, err)

	repo := NewThreadRepo(connPool, 2)

	t.Run("successfully shows the first 3 posts", func(t *testing.T) {
		posts, count, err := repo.ListPostsCollapsible(context.Background(), 3, 1, 2)
		require.NoError(t, err)

		assert.Len(t, posts, 3)
		assert.Equal(t, 7, count)

		assert.Equal(t, "I was in Hawaii!", posts[1].Body)
		assert.Equal(t, 1, posts[1].MemberID)

		assert.Equal(t, "I was in Amish Country!", posts[2].Body)
		assert.Equal(t, 2, posts[2].MemberID)

		assert.Equal(t, "I was too!", posts[0].Body)
		assert.Equal(t, 2, posts[0].MemberID)
	})

	t.Run("successfully shows the next 3 posts", func(t *testing.T) {
		posts, count, err := repo.ListPostsCollapsible(context.Background(), 6, 1, 2)
		require.NoError(t, err)

		assert.Len(t, posts, 6)
		assert.Equal(t, 4, count)

		assert.Equal(t, "Wow!", posts[0].Body)
		assert.Equal(t, 1, posts[0].MemberID)

		assert.Equal(t, "I was on vacation!", posts[2].Body)
		assert.Equal(t, 1, posts[2].MemberID)

		assert.Equal(t, "I was in Amish Country!", posts[5].Body)
		assert.Equal(t, 2, posts[5].MemberID)
	})

	t.Run("successfully shows all posts", func(t *testing.T) {
		posts, count, err := repo.ListPostsCollapsible(context.Background(), 10, 1, 2)
		require.NoError(t, err)

		assert.Len(t, posts, 10)
		assert.Equal(t, 0, count)

		assert.Equal(t, "Hello, BCO!", posts[0].Body)
		assert.Equal(t, 1, posts[0].MemberID)

		assert.Equal(t, "I was in Amish Country!", posts[9].Body)
		assert.Equal(t, 2, posts[9].MemberID)
	})

	t.Run("returns an error when the number of posts to show is 0", func(t *testing.T) {
		_, _, err := repo.ListPostsCollapsible(context.Background(), 0, 1, 2)
		require.Error(t, err)

		assert.Equal(t, "toShow cannot be 0", err.Error())
	})
}
