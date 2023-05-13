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

	require.NoError(t, err)

	return nil
}
