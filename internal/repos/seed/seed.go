package seed

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

const (
	insertThreadPost = `INSERT INTO thread_post (thread_id, member_id, member_ip, body, date_posted) VALUES ($1, $2, $3, $4, $5)`
	insertThread     = `INSERT INTO thread (member_id, first_post_id, last_member_id, subject) VALUES ($1, $2, $3, $4)`
	insertMember     = `INSERT INTO member (name, pass, ip, email_signup, postalcode, secret) VALUES ($1, $2, $3, $4, $5, $6)`
	insertFavorite   = `INSERT INTO favorite (member_id, thread_id) VALUES ($1, $2)`
	updateIgnored    = `UPDATE thread_member SET ignore = $1 WHERE thread_id = $2 AND member_id = $3`
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
	_, err = pool.Exec(context.Background(), insertThread, 2, 3, 1, "THYROID")
	_, err = pool.Exec(context.Background(), insertThreadPost, 3, 1, "127.0.0.1", "Doctor BCO: I'm not sure if I have a thyroid problem", time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), insertMember, "elliott", "pass", "127.0.0.1", "elliott@test.com", "12345", "topsecret")
	_, err = pool.Exec(context.Background(), insertThread, 3, 2, 1, "2017 politics thread")
	_, err = pool.Exec(context.Background(), insertThread, 3, 2, 1, "2018 politics thread")
	_, err = pool.Exec(context.Background(), "UPDATE thread SET date_last_posted = $1 WHERE id = 1", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), "UPDATE thread SET date_last_posted = $1 WHERE id = 2", time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), "UPDATE thread SET date_last_posted = $1 WHERE id = 3", time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), "UPDATE thread SET date_last_posted = $1 WHERE id = 4", time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), "UPDATE thread SET date_last_posted = $1 WHERE id = 5", time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC))
	_, err = pool.Exec(context.Background(), insertFavorite, 1, 1)
	_, err = pool.Exec(context.Background(), insertFavorite, 1, 4)
	_, err = pool.Exec(context.Background(), updateIgnored, true, 1, 1)
	_, err = pool.Exec(context.Background(), updateIgnored, true, 3, 1)

	require.NoError(t, err)

	return nil
}
