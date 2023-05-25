package threadrepo

import (
	"context"
	_ "embed"
	"goBoard/internal/core/domain"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed queries/list_threads.sql
var listThreadsQuery string

//go:embed queries/count_threads.sql
var countThreadsQuery string

//go:embed queries/list_posts.sql
var listPostsQuery string

type ThreadRepo struct {
	connPool *pgxpool.Pool
}

func NewThreadRepo(pool *pgxpool.Pool) ThreadRepo {
	return ThreadRepo{pool}
}

func (r ThreadRepo) SavePost(post domain.Post) (int, error) {
	var id int
	err := r.connPool.QueryRow(context.Background(), "INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4) RETURNING id", post.ThreadID, post.MemberID, post.MemberIP, post.Text).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r ThreadRepo) GetPostByID(id int) (*domain.Post, error) {
	var post domain.Post
	var cidr pgtype.CIDR
	var err = r.connPool.QueryRow(context.Background(), "SELECT thread_post.id, thread_Id, member_id, m.name, member_ip, body, date_posted FROM thread_post LEFT JOIN member m on m.id = thread_post.member_id WHERE thread_post.id = $1", id).Scan(
		&post.ID,
		&post.ThreadID,
		&post.MemberID,
		&post.MemberName,
		&cidr,
		&post.Text,
		&post.Timestamp,
	)
	if err != nil {
		return nil, err
	}

	post.MemberIP = cidr.IPNet.String()

	return &post, nil
}

func (r ThreadRepo) GetThreadByID(id int) (*domain.Thread, error) {
	var thread domain.Thread
	err := r.connPool.QueryRow(context.Background(), "SELECT id, subject, date_posted, member_id, views FROM thread WHERE id = $1", id).Scan(&thread.ID, &thread.Subject, &thread.Timestamp, &thread.MemberID, &thread.Views)
	if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (r ThreadRepo) ListThreads(limit, offset int) (*domain.ThreadPage, error) {
	var threads []domain.Thread
	threadPage := &domain.ThreadPage{}
	rows, err := r.connPool.Query(context.Background(), listThreadsQuery, limit, offset, nil)
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
	threadPage.Threads = threads

	var totalThreads int
	threadCountRows, err := r.connPool.Query(context.Background(), countThreadsQuery)
	for threadCountRows.Next() {
		err := threadCountRows.Scan(&totalThreads)
		if err != nil {
			return nil, err
		}
	}
	threadPage.TotalPages = totalThreads / limit

	return threadPage, nil
}

func (r ThreadRepo) ListThreadsByMemberID(memberID int, limit, offset int) ([]domain.Thread, error) {
	var threads []domain.Thread
	rows, err := r.connPool.Query(context.Background(), listThreadsQuery, limit, offset, memberID)
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

func (r ThreadRepo) ListPostsForThread(limit, offset, id int) ([]domain.Post, error) {
	var posts []domain.Post
	var cidr pgtype.CIDR
	rows, err := r.connPool.Query(context.Background(), listPostsQuery, limit, offset, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.ID, &post.Timestamp, &post.MemberID, &post.MemberName, &post.Text, &cidr, &post.ThreadSubject, &post.ThreadID, &post.IsAdmin)
		if err != nil {
			return nil, err
		}

		post.MemberIP = cidr.IPNet.String()

		posts = append(posts, post)
	}

	return posts, nil
}
