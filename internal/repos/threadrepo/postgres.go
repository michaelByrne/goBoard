package threadrepo

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"goBoard/internal/core/domain"
	"time"

	sq "github.com/Masterminds/squirrel"
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

//go:embed queries/list_posts_collapsible.sql
var listPostsCollapsibleQuery string

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

func (r ThreadRepo) ListPostsCollapsible(ctx context.Context, toShow, threadID, memberID int) (posts []domain.ThreadPost, collapsed int, err error) {
	if toShow == 0 {
		err = errors.New("toShow cannot be 0")
		return
	}

	rows, err := r.connPool.Query(ctx, listPostsCollapsibleQuery, threadID, memberID, toShow)
	if err != nil {
		return
	}

	for rows.Next() {
		var post domain.ThreadPost
		var cidr pgtype.CIDR
		err = rows.Scan(&post.ID, &post.Timestamp, &post.MemberID, &post.MemberName, &post.Body, &cidr, &post.ParentSubject, &post.ParentID, &post.IsAdmin, &post.RowNumber, &collapsed)
		if err != nil {
			return
		}

		post.MemberIP = cidr.IPNet.String()

		posts = append(posts, post)
	}

	return
}

func (r ThreadRepo) SavePost(post domain.ThreadPost) (int, error) {
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

	var id int
	err = tx.QueryRow(context.Background(), "INSERT INTO thread_post (thread_id, member_id, member_ip, body) VALUES ($1, $2, $3, $4) RETURNING id", post.ParentID, post.MemberID, post.MemberIP, post.Body).Scan(&id)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(context.Background(), "UPDATE thread_viewer SET dotted = true WHERE thread_id = $1 AND member_id = $2 AND undot IS FALSE", post.ParentID, post.MemberID)
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
	tvQuery := "INSERT INTO thread_viewer (thread_id, member_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"

	_, err := r.connPool.Exec(context.Background(), tvQuery, id, memberID)
	if err != nil {
		return nil, err
	}

	var thread domain.Thread
	err = r.connPool.QueryRow(context.Background(), getThreadByIDQuery, id, memberID).Scan(&thread.ID, &thread.Subject, &thread.Timestamp, &thread.MemberID, &thread.Views, &thread.Dotted, &thread.Undot, &thread.Ignored, &thread.Favorite)
	if err != nil {
		return nil, err
	}

	return &thread, nil
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

	_, err = tx.Exec(context.Background(), "INSERT INTO thread_viewer (thread_id, member_id, dotted) VALUES ($1, $2, true)", threadID, thread.MemberID)
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

func (r ThreadRepo) ToggleDot(ctx context.Context, memberID, threadID int) (bool, error) {
	query := `
	UPDATE thread_viewer 
	SET 
 		dotted = NOT dotted,
  		undot = CASE
        	WHEN dotted IS TRUE THEN true
        	ELSE undot
    	END
	WHERE thread_id = $1 
 	   AND member_id = $2 
	RETURNING dotted;
	`
	var dot bool
	err := r.connPool.QueryRow(ctx, query, threadID, memberID).Scan(&dot)
	if err != nil {
		return false, err
	}

	return dot, nil
}

func (r ThreadRepo) ToggleFavorite(ctx context.Context, memberID, threadID int) (bool, error) {
	query := "SELECT toggle_favorite($1, $2)"

	var favorite int
	err := r.connPool.QueryRow(ctx, query, memberID, threadID).Scan(&favorite)
	if err != nil {
		return false, err
	}

	return favorite == 1, nil
}

func (r ThreadRepo) ToggleIgnore(ctx context.Context, memberID, threadID int) (bool, error) {
	query := "UPDATE thread_viewer SET ignored=NOT ignored WHERE thread_id=$1 AND member_id=$2 RETURNING ignored"

	var ignored bool
	err := r.connPool.QueryRow(ctx, query, threadID, memberID).Scan(&ignored)
	if err != nil {
		return false, err
	}

	return ignored, nil
}

func (r ThreadRepo) ViewThread(ctx context.Context, memberID, threadID int) (int, error) {
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

	var views int
	err = tx.QueryRow(ctx, "UPDATE thread SET views = views + 1 WHERE id = $1 RETURNING views", threadID).Scan(&views)
	if err != nil {
		return 0, err
	}

	return views, nil
}

func (r ThreadRepo) ListThreads(ctx context.Context, cursors domain.Cursors, limit, memberID int, filter domain.ThreadFilter) ([]domain.Thread, domain.Cursors, error) {
	dot := sq.Select("COALESCE(tv.dotted, false)")
	undot := sq.Select("COALESCE(tv.undot, false)")
	favorite := sq.Case().When("f.thread_id IS NOT NULL", "true").Else("false")

	innerQuery := sq.Select("t.id", "m.name as name", "t.date_last_posted", "t.subject", "t.date_posted", "t.posts", "t.views", "l.name as last_poster_name").
		Column(sq.Alias(dot, "dot")).
		Column(sq.Alias(undot, "undot")).
		Column(sq.Alias(favorite, "favorite")).
		From("thread t")

	ignoredThreadQuery := sq.Select("tv.thread_id").From("thread_viewer tv").Where("tv.ignored = true AND tv.member_id = ?", memberID).OrderBy("tv.last_viewed DESC")
	ignoredMemberQuery := sq.Select("mi.ignore_member_id").From("member_ignore mi").LeftJoin("member m ON m.id = mi.member_id").Where("mi.member_id = ?", memberID).OrderBy("m.name")

	totalMinusIgnoredQuery := sq.Select("COUNT(*)").From("thread")

	if filter == domain.ThreadFilterIgnored {
		totalMinusIgnoredQuery = totalMinusIgnoredQuery.Where(SubQueryIN("id", ignoredThreadQuery)).Where(SubQueryNOTIN("member_id", ignoredMemberQuery))
	} else if filter == domain.ThreadFilterAll {
		totalMinusIgnoredQuery = totalMinusIgnoredQuery.Where(SubQueryNOTIN("id", ignoredThreadQuery)).Where(SubQueryNOTIN("member_id", ignoredMemberQuery))
	} else if filter == domain.ThreadFilterCreated {
		innerQuery = innerQuery.Where("t.member_id = ?", memberID)
	} else if filter == domain.ThreadFilterParticipated {
		innerQuery = innerQuery.InnerJoin("thread_member tm ON tm.thread_id = t.id").Where("tm.member_id = ?", memberID)
	}

	if filter == domain.ThreadFilterFavorites {
		innerQuery = innerQuery.InnerJoin("favorite f ON f.thread_id = t.id").Where("f.member_id = ?", memberID)
	} else {
		innerQuery = innerQuery.JoinClause("LEFT OUTER JOIN favorite f ON f.thread_id = t.id AND f.member_id = ?", memberID)
	}

	rowsLeftQuery := totalMinusIgnoredQuery

	pagination := innerQuery.InnerJoin("member m ON m.id = t.member_id").
		InnerJoin("member l ON l.id = t.last_member_id").
		InnerJoin("thread_post tp ON tp.id = t.first_post_id").
		JoinClause("LEFT OUTER JOIN thread_viewer tv ON tv.thread_id = t.id AND tv.member_id = ?", memberID).
		Where(SubQueryNOTIN("t.member_id", ignoredMemberQuery))

	if filter == domain.ThreadFilterIgnored {
		pagination = pagination.Where(SubQueryIN("t.id", ignoredThreadQuery))
	} else {
		pagination = pagination.Where(SubQueryNOTIN("t.id", ignoredThreadQuery))
	}

	if cursors.Next != "" && cursors.Prev != "" {
		return nil, domain.Cursors{}, errors.New("two cursors cannot be provided at the same time")
	}

	// Going forward
	if cursors.Next != "" {
		rowsLeftQuery = rowsLeftQuery.Where("thread.date_last_posted < ?", cursors.Next)
		pagination = pagination.Where("t.date_last_posted < ?", cursors.Next).OrderBy("date_last_posted DESC").Limit(uint64(limit))
	}

	// Going backward
	if cursors.Prev != "" {
		rowsLeftQuery = rowsLeftQuery.Where("thread.date_last_posted > ?", cursors.Prev)
		pagination = pagination.Where("t.date_last_posted > ?", cursors.Prev).OrderBy("date_last_posted ASC").Limit(uint64(limit))
	}

	// No cursors: Going forward from the beginning
	if cursors.Next == "" && cursors.Prev == "" {
		pagination = pagination.OrderBy("t.date_last_posted DESC").Limit(uint64(limit))
	}

	stmt := sq.Select("id", "name", "date_last_posted", "subject", "date_posted", "posts", "views", "last_poster_name").
		Column(sq.Alias(rowsLeftQuery, "rows_left")).
		Column(sq.Alias(totalMinusIgnoredQuery, "total")).
		Column("dot").
		Column("undot").
		Column("favorite").
		FromSelect(pagination, "t").OrderBy("date_last_posted DESC").PlaceholderFormat(sq.Dollar)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, domain.Cursors{}, err
	}

	rows, err := r.connPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, domain.Cursors{}, err
	}
	defer rows.Close()

	var (
		threads  []domain.Thread
		rowsLeft int
		total    int
	)

	for rows.Next() {
		var thread domain.Thread
		err = rows.Scan(
			&thread.ID,
			&thread.MemberName,
			&thread.DateLastPosted,
			&thread.Subject,
			&thread.DatePosted,
			&thread.NumPosts,
			&thread.Views,
			&thread.LastPosterName,
			&rowsLeft,
			&total,
			&thread.Dotted,
			&thread.Undot,
			&thread.Favorite,
		)
		if err != nil {
			return nil, domain.Cursors{}, err
		}

		threads = append(threads, thread)
	}

	if len(threads) == 0 {
		return threads, domain.Cursors{}, nil
	}

	var (
		prevCursor string // cursor we return when there is a previous page
		nextCursor string // cursor we return when there is a next page
	)

	switch {

	// *If there are no results we don't have to compute the cursors
	case rowsLeft < 0:

	// *On A, direction A->E (going forward), return only next cursor
	case cursors.Prev == "" && cursors.Next == "":
		nextCursor = threads[len(threads)-1].DateLastPosted.UTC().Format(time.RFC3339Nano)

	// *On E, direction A->E (going forward), return only prev cursor
	case cursors.Next != "" && rowsLeft == len(threads):
		prevCursor = threads[0].DateLastPosted.UTC().Format(time.RFC3339Nano)

	// *On A, direction E->A (going backward), return only next cursor
	case cursors.Prev != "" && rowsLeft == len(threads):
		nextCursor = threads[len(threads)-1].DateLastPosted.UTC().Format(time.RFC3339Nano)

	// *On E, direction E->A (going backward), return only prev cursor
	case cursors.Prev != "" && total == rowsLeft:
		prevCursor = threads[0].DateLastPosted.UTC().Format(time.RFC3339Nano)

	// *Somewhere in the middle
	default:
		nextCursor = threads[len(threads)-1].DateLastPosted.UTC().Format(time.RFC3339Nano)
		prevCursor = threads[0].DateLastPosted.UTC().Format(time.RFC3339Nano)
	}

	return threads, domain.Cursors{Next: nextCursor, Prev: prevCursor}, nil
}

func SubQueryNOTIN(property string, query sq.SelectBuilder) sq.Sqlizer {
	sql, args, _ := query.ToSql()
	subQuery := fmt.Sprintf("%s NOT IN (%s)", property, sql)
	return sq.Expr(subQuery, args...)
}

func SubQueryIN(property string, query sq.SelectBuilder) sq.Sqlizer {
	sql, args, _ := query.ToSql()
	subQuery := fmt.Sprintf("%s IN (%s)", property, sql)
	return sq.Expr(subQuery, args...)
}
