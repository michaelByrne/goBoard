package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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

	//b.Queue("INSERT INTO member (id, name, pass, postalcode, email_signup, secret, ip) VALUES ($1, $2, $3, $4, $5, $6, $7)", 1, "gofreescout", "test", "97217", "mpbyrne@hotmail.com", "fishsticks", "172.0.0.1")
	//b.Queue("INSERT INTO member (id, name, pass, postalcode, email_signup, secret, ip) VALUES ($1, $2, $3, $4, $5, $6, $7)", 2, "elliott", "test2", "97217", "admin@admin.net", "admin", "172.0.0.1")

	// for i := 1; i < 150; i++ {
	// 	//dateLastPosted := time.Date(2021, 1, i, 0, 0, 0, 0, time.UTC)
	// 	b.Queue("INSERT INTO thread (subject, member_id, last_member_id) VALUES ($1, $2, $3)", DATA[i%30], 1, 1)
	// }

	// for i := 1; i < 50; i++ {
	// 	dateLastPosted := time.Date(2021, 1, i*2, 0, 0, 0, 0, time.UTC)
	// 	b.Queue("INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ($1, $2, $3, $4)", DATA[i%30], 2, 1, dateLastPosted.Format(time.RFC3339Nano))
	// }

	for i := 4; i < 149; i++ {
		postDate := randDate()
		b.Queue("INSERT INTO thread (subject, member_id, last_member_id, date_last_posted) VALUES ($1, $2, $3, $4)", DATA[i%30], 1, 1, postDate.Format(time.RFC3339Nano))
		b.Queue("INSERT INTO thread_post (member_id, thread_id, body, member_ip, date_posted) VALUES ($1, $2, $3, $4, $5)", 1, i, DATA[i%30], "172.0.0.1", postDate.Format(time.RFC3339Nano))
		b.Queue("UPDATE thread SET date_last_posted = $1 WHERE id = $2", postDate.Format(time.RFC3339Nano), i)
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

func randDate() time.Time {
	min := time.Date(2024, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
