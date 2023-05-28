package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"math/rand"
	"time"
)

var DATA = [30]string{"Lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit", "Mauris", "faucibus", "lectus", "eget", "cursus", "tempus", "ligula", "orci", "mattis", "massa", "nec", "eleifend", "lorem", "ipsum", "congue", "erat", "Pellentesque", "suscipit", "semper", "sapien", "sed", "luctus"}

func main() {
	dbURI := "postgres://boardking:test@localhost:5432/board?sslmode=disable"
	var pool, err = pgxpool.Connect(context.Background(), dbURI)
	if err != nil {
		fmt.Println(err)
	}

	begin, err := pool.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	b := &pgx.Batch{}

	b.Queue("INSERT INTO member (id, name, pass, postalcode, email_signup, secret, ip) VALUES ($1, $2, $3, $4, $5, $6, $7)", 1, "gofreescout", "test", "97217", "mpbyrne@hotmail.com", "fishsticks", "172.0.0.1")
	b.Queue("INSERT INTO member (id, name, pass, postalcode, email_signup, secret, ip) VALUES ($1, $2, $3, $4, $5, $6, $7)", 2, "elliott", "test2", "97217", "admin@admin.net", "admin", "172.0.0.1")

	for i := 1; i < 100; i++ {
		b.Queue("INSERT INTO thread (id, subject, member_id, last_member_id) VALUES ($1, $2, $3, $4)", i, DATA[i%30], 1, 1)

		b.Queue("INSERT INTO thread (id, subject, member_id, last_member_id) VALUES ($1, $2, $3, $4)", i*100, DATA[i%30], 2, 2)
	}

	for i := 1; i < 100; i++ {
		b.Queue("INSERT INTO board.public.thread_post (id, member_id, thread_id, body, member_ip) VALUES ($1, $2, $3, $4, $5)", i, 1, i%100, DATA[i%30], "172.0.0.1")
		b.Queue("INSERT INTO board.public.thread_post (id, member_id, thread_id, body, member_ip) VALUES ($1, $2, $3, $4, $5)", i*100, 2, i*100, DATA[i%30], "172.0.0.1")
		randDay := rand.Intn(28) + 1
		randMonth := rand.Intn(12) + 1
		randYear := rand.Intn(20) + 2000
		b.Queue("UPDATE thread SET date_last_posted = $1 WHERE id = $2", time.Date(randYear, time.Month(randMonth), randDay, 0, 0, 0, 0, time.UTC), i)
		randDay = rand.Intn(28) + 1
		randMonth = rand.Intn(12) + 1
		randYear = rand.Intn(20) + 2000
		b.Queue("UPDATE thread SET date_last_posted = $1 WHERE id = $2", time.Date(randYear, time.Month(randMonth), randDay, 0, 0, 0, 0, time.UTC), i*100)
	}

	results := begin.SendBatch(context.Background(), b)

	var qerr error
	var rows pgx.Rows
	for qerr == nil {
		rows, qerr = results.Query()
		rows.Close()
	}

	err = begin.Commit(context.Background())
	if err != nil {
		log.Fatal(err)
	}

}
