package threadrepo

import (
	"context"
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
