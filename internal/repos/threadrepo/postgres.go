package threadrepo

import (
	"context"
	_ "embed"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"goBoard/internal/core/domain"
	"net"
)

//go:embed queries/list_threads.sql
var listThreadsQuery string

type ThreadRepo struct {
	connPool *pgxpool.Pool
}

func NewThreadRepo(pool *pgxpool.Pool) ThreadRepo {
	return ThreadRepo{pool}
}

func (r ThreadRepo) SavePost(post domain.Post) (int, error) {
	ip, _, err := net.ParseCIDR(post.MemberIP)
	if err != nil {
		return 0, err
	}

	var id int
	err = r.connPool.QueryRow(context.Background(), "INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4) RETURNING id", post.ThreadID, post.MemberID, ip, post.Text).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r ThreadRepo) GetPostByID(id int) (*domain.Post, error) {
	var post domain.Post
	var cidr pgtype.CIDR
	err := r.connPool.QueryRow(context.Background(), "SELECT id, thread_Id, member_id, member_ip, body, date_posted FROM thread_post WHERE id = $1", id).Scan(&post.ID, &post.ThreadID, &post.MemberID, &cidr, &post.Text, &post.Timestamp)
	if err != nil {
		return nil, err
	}

	post.MemberIP = cidr.IPNet.String()

	return &post, nil
}

func (r ThreadRepo) GetPostsByThreadID(threadID int) ([]domain.Post, error) {
	var posts []domain.Post
	var cidr pgtype.CIDR
	rows, err := r.connPool.Query(context.Background(), "SELECT id, thread_Id, member_id, member_ip, body, date_posted FROM thread_post WHERE thread_id = $1", threadID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.ID, &post.ThreadID, &post.MemberID, &cidr, &post.Text, &post.Timestamp)
		if err != nil {
			return nil, err
		}

		post.MemberIP = cidr.IPNet.String()

		posts = append(posts, post)
	}

	return posts, nil
}

func (r ThreadRepo) GetThreadByID(id int) (*domain.Thread, error) {
	var thread domain.Thread
	err := r.connPool.QueryRow(context.Background(), "SELECT id, subject, date_posted, member_id, views FROM thread WHERE id = $1", id).Scan(&thread.ID, &thread.Subject, &thread.Timestamp, &thread.MemberID, &thread.Views)
	if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (r ThreadRepo) ListThreads(limit int) ([]domain.Thread, error) {
	var threads []domain.Thread
	rows, err := r.connPool.Query(context.Background(), listThreadsQuery, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var thread domain.Thread
		err := rows.Scan(
			&thread.ID,
			&thread.DateLastPosted,
			&thread.MemberID,
			&thread.MemberName,
			&thread.LastPosterID,
			&thread.LastPosterName,
			&thread.Subject,
			&thread.NumPosts,
			&thread.Views,
			&thread.LastPostText,
			&thread.Sticky,
			&thread.Locked,
			&thread.Legendary,
		)
		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}

func (r ThreadRepo) SaveThread(thread domain.Thread) (int, error) {
	tx, err := r.connPool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	var threadID int
	err = tx.QueryRow(context.Background(), "INSERT INTO thread (subject, member_id, last_member_id) VALUES ($1, $2, $3) RETURNING id", thread.Subject, thread.MemberID, thread.LastPosterID).Scan(&threadID)
	if err != nil {
		return 0, err
	}

	var postID int
	err = tx.QueryRow(context.Background(), "INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4) RETURNING id", threadID, thread.MemberID, thread.MemberIP, thread.FirstPostText).Scan(&postID)
	if err != nil {
		return 0, err
	}

	return threadID, nil
}

func (r ThreadRepo) DeleteThread(id int) error {
	_, err := r.connPool.Exec(context.Background(), "DELETE FROM thread WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
