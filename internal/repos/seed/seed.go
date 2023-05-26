package seed

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	insertThreadPost = `INSERT INTO thread_post (thread_id, member_id, member_ip, body, date_posted) VALUES ($1, $2, $3, $4, $5)`
	insertThread     = `INSERT INTO thread (member_id, first_post_id, last_member_id, subject) VALUES ($1, $2, $3, $4)`
	insertMember     = `INSERT INTO member (name, pass, ip, email_signup, postalcode, secret) VALUES ($1, $2, $3, $4, $5, $6)`
)

func SeedData(t *testing.T, pool *pgxpool.Pool) error {
	var err error
	_, err = pool.Exec(context.Background(), insertMember, "admin", "admin", "127.0.0.1", "admin@test.com", "12345", "topsecret")
	_, err = pool.Exec(context.Background(), insertThread, 1, 1, 1, "Hello, BCO")
	_, err = pool.Exec(context.Background(), insertThreadPost, 1, 1, "127.0.0.1", "Attn. Roxy", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), insertThreadPost, 1, 1, "127.0.0.2", "WCFRP", time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), insertMember, "gofreescout", "test", "127.0.0.2", "gofreescout@gmail.com", "48225", "topsecret")
	_, err = pool.Exec(context.Background(), insertThread, 2, 3, 1, "It stinks! A new moratorium thread")
	_, err = pool.Exec(context.Background(), insertThreadPost, 2, 2, "127.0.0.1", "I listened to a podcast earlier that had five minutes of ads at the beginning", time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), insertThreadPost, 2, 1, "127.0.0.1", "moratorium on anything to do with AI", time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), insertThreadPost, 2, 2, "127.0.0.1", "small d democratic", time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), "UPDATE thread SET date_last_posted = $1 WHERE id = 2", time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), "UPDATE thread SET date_last_posted = $1 WHERE id = 1", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))

	require.NoError(t, err)

	return nil
}
