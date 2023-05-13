package threadrepo

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"goBoard/internal/core/domain"
)

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
