package threadrepo

import (
	"context"
	_ "embed"
	"goBoard/internal/core/domain"
	"time"

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

//go:embed queries/lists_posts_cursor_forward.sql
var listPostsCursorQuery string

//go:embed queries/list_threads_cursor_forward.sql
var listThreadsCursorForwardQuery string

//go:embed queries/list_threads_cursor_reverse.sql
var listThreadsCursorReverseQuery string

//go:embed queries/get_thread_by_id.sql
var getThreadByIDQuery string

type ThreadRepo struct {
	connPool              *pgxpool.Pool
	defaultMaxThreadLimit int
}

func NewThreadRepo(pool *pgxpool.Pool, defaultMaxThreadLimit int) ThreadRepo {
	return ThreadRepo{
		connPool:              pool,
		defaultMaxThreadLimit: defaultMaxThreadLimit,
	}
}

func (r ThreadRepo) SavePost(post domain.ThreadPost) (int, error) {
	var id int
	err := r.connPool.QueryRow(context.Background(), "INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4) RETURNING id", post.ParentID, post.MemberID, post.MemberIP, post.Body).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r ThreadRepo) GetPostByID(id int) (*domain.ThreadPost, error) {
	var post domain.ThreadPost
	var cidr pgtype.CIDR
	var err = r.connPool.QueryRow(context.Background(), "SELECT thread_post.id, thread_Id, member_id, m.name, member_ip, body, date_posted FROM thread_post LEFT JOIN member m on m.id = thread_post.member_id WHERE thread_post.id = $1", id).Scan(
		&post.ID,
		&post.ParentID,
		&post.MemberID,
		&post.MemberName,
		&cidr,
		&post.Body,
		&post.Timestamp,
	)
	if err != nil {
		return nil, err
	}

	post.MemberIP = cidr.IPNet.String()

	return &post, nil
}

func (r ThreadRepo) GetThreadByID(id, memberID int) (*domain.Thread, error) {
	var thread domain.Thread
	err := r.connPool.QueryRow(context.Background(), getThreadByIDQuery, id, memberID).Scan(&thread.ID, &thread.Subject, &thread.Timestamp, &thread.MemberID, &thread.Views, &thread.Dotted, &thread.Ignored)
	if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (r ThreadRepo) ListThreads(limit, offset int) (*domain.SiteContext, error) {
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

	siteContext := &domain.SiteContext{ThreadPage: *threadPage}

	return siteContext, nil
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

func (r ThreadRepo) ListPostsForThread(limit, offset, id, memberID int) ([]domain.ThreadPost, error) {
	var posts []domain.ThreadPost
	var cidr pgtype.CIDR
	rows, err := r.connPool.Query(context.Background(), listPostsQuery, limit, offset, id, memberID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post domain.ThreadPost
		err := rows.Scan(&post.ID, &post.Timestamp, &post.MemberID, &post.MemberName, &post.Body, &cidr, &post.ParentSubject, &post.ParentID, &post.IsAdmin)
		if err != nil {
			return nil, err
		}

		post.MemberIP = cidr.IPNet.String()

		posts = append(posts, post)
	}

	return posts, nil
}

func (r ThreadRepo) ListPostsForThreadByCursor(limit, id int, cursor *time.Time) ([]domain.ThreadPost, error) {
	var posts []domain.ThreadPost
	var cidr pgtype.CIDR
	rows, err := r.connPool.Query(context.Background(), listPostsQuery, limit, id, cursor)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post domain.ThreadPost
		err := rows.Scan(&post.ID, &post.Timestamp, &post.MemberID, &post.MemberName, &post.Body, &cidr, &post.ParentSubject, &post.ParentID, &post.IsAdmin)
		if err != nil {
			return nil, err
		}

		post.MemberIP = cidr.IPNet.String()

		posts = append(posts, post)
	}

	return posts, nil
}

func (r ThreadRepo) ListThreadsByCursorForward(limit int, cursor *time.Time, memberID int) ([]domain.Thread, error) {
	var threads []domain.Thread
	rows, err := r.connPool.Query(context.Background(), listThreadsCursorForwardQuery, limit, cursor, memberID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var thread domain.Thread
		err := rows.Scan(
			&thread.ID,
			&thread.DateLastPosted,
			&thread.DatePosted,
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
			&thread.Dotted,
		)
		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}

func (r ThreadRepo) ListThreadsByCursorReverse(limit int, cursor *time.Time, memberID int) ([]domain.Thread, error) {
	var threads []domain.Thread
	rows, err := r.connPool.Query(context.Background(), listThreadsCursorReverseQuery, cursor, limit, memberID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var thread domain.Thread
		err := rows.Scan(
			&thread.ID,
			&thread.DateLastPosted,
			&thread.DatePosted,
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
			&thread.Dotted,
		)
		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}

func (r ThreadRepo) PeekPrevious(timestamp *time.Time) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM thread WHERE date_last_posted > $1)"
	var exists bool
	err := r.connPool.QueryRow(context.Background(), query, timestamp).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r ThreadRepo) UndotThread(ctx context.Context, memberID, threadID int) error {
	query := "UPDATE thread_member SET undot=true WHERE thread_id=$1 AND member_id=$2"
	_, err := r.connPool.Exec(ctx, query, threadID, memberID)
	if err != nil {
		return err
	}

	return nil
}

func (r ThreadRepo) ToggleIgnore(ctx context.Context, memberID, threadID int, ignore bool) error {
	query := "UPDATE thread_member SET ignore=$3 WHERE thread_id=$1 AND member_id=$2"
	_, err := r.connPool.Exec(ctx, query, threadID, memberID, ignore)
	if err != nil {
		return err
	}

	return nil
}
