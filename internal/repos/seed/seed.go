package seed

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	insertThreadPost = `INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4)`
	insertThread     = `INSERT INTO thread (member_id, first_post_id, last_member_id, subject) VALUES ($1, $2, $3, $4)`
	insertMember     = `INSERT INTO member (name, pass, ip, email_signup, postalcode, secret) VALUES ($1, $2, $3, $4, $5, $6)`
)

func SeedData(t *testing.T, pool *pgxpool.Pool) error {
	var err error
	_, err = pool.Exec(context.Background(), insertMember, "admin", "admin", "127.0.0.1", "admin@test.com", "12345", "topsecret")
	_, err = pool.Exec(context.Background(), insertThread, 1, 1, 1, "Hello, BCO")
	_, err = pool.Exec(context.Background(), insertThreadPost, 1, 1, "127.0.0.1", "Attn. Roxy")
	_, err = pool.Exec(context.Background(), insertThreadPost, 1, 1, "127.0.0.2", "WCFRP")
	_, err = pool.Exec(context.Background(), insertMember, "gofreescout", "test", "127.0.0.2", "gofreescout@gmail.com", "48225", "topsecret")
	_, err = pool.Exec(context.Background(), insertThread, 2, 3, 1, "It stinks! A new moratorium thread")
	_, err = pool.Exec(context.Background(), insertThreadPost, 2, 2, "127.0.0.1", "I listened to a podcast earlier that had five minutes of ads at the beginning")
	_, err = pool.Exec(context.Background(), insertThreadPost, 2, 1, "127.0.0.1", "moratorium on anything to do with AI")
	_, err = pool.Exec(context.Background(), insertThreadPost, 2, 2, "127.0.0.1", "small d democratic")

	require.NoError(t, err)

	return nil
}
