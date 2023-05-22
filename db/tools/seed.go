package main

import (
	"strconv"
	"fmt"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	dbURI := "postgres://boardking:test@localhost:5432/board?sslmode=disable"
	var pool, err = pgxpool.Connect(context.Background(), dbURI)
	if err != nil { 
		fmt.Println(err) 
	}
	
	pool.Exec(context.Background(), `INSERT INTO member (name, pass, ip, email_signup, postalcode, secret) VALUES ($1, $2, $3, $4, $5, $6)`, "admin", "admin", "127.0.0.1", "admin@test.com", "12345", "topsecret")
	pool.Exec(context.Background(), `INSERT INTO member (name, pass, ip, email_signup, postalcode, secret) VALUES ($1, $2, $3, $4, $5, $6)`, "gofreescout", "test", "127.0.0.2", "gofreescout@gmail.com", "48225", "topsecret")

	n := 1
	for n < 10000 {
		t := strconv.Itoa(n)
		pool.Exec(context.Background(), `INSERT INTO thread (member_id, first_post_id, last_member_id, subject) VALUES ($1, $2, $3, $4)`, 1, 1, 1, "Hello, BCO"+t)
		pool.Exec(context.Background(), `INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4)`, n, 1, "127.0.0.1", "Attn. Roxy")
		pool.Exec(context.Background(), `INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4)`, n, 1, "127.0.0.2", "WCFRP")
		pool.Exec(context.Background(), `INSERT INTO thread (member_id, first_post_id, last_member_id, subject) VALUES ($1, $2, $3, $4)`, 2, 3, 1, "It stinks! A new moratorium thread"+t)
		pool.Exec(context.Background(), `INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4)`, n+1, 2, "127.0.0.1", "I listened to a podcast earlier that had five minutes of ads at the beginning")
		pool.Exec(context.Background(), `INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4)`, n+1, 1, "127.0.0.1", "moratorium on anything to do with AI")
		pool.Exec(context.Background(), `INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4)`, n+1, 2, "127.0.0.1", "small d democratic")
		n++
	}
}
