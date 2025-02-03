package messagerepo

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"goBoard/internal/core/domain"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed queries/get_messages_with_cursor_forward.sql
var getMessagesWithCursorForward string

//go:embed queries/list_message_posts_cursor_forward.sql
var listMessagePostsCursorForward string

//go:embed queries/get_messages_with_cursor_reverse.sql
var getMessagesWithCursorReverse string

//go:embed queries/get_message_post_by_id.sql
var getMessagePostByID string

//go:embed queries/list_posts_collapsible.sql
var listPostsCollapsible string

//go:embed queries/get_message_by_id.sql
var getMessageByID string

const insertMessage = `INSERT INTO message (member_id, subject, last_member_id) VALUES ($1, $2, $3) RETURNING id`
const insertMessagePost = `INSERT INTO message_post (message_id, member_id, body, member_ip) VALUES ($1, $2, $3, $4) RETURNING id`
const insertMessageMember = `INSERT INTO message_member (message_id, member_id) VALUES ($1, $2)`

type MessageRepo struct {
	connPool *pgxpool.Pool
}

func NewMessageRepo(pool *pgxpool.Pool) MessageRepo {
	return MessageRepo{
		connPool: pool,
	}
}

func (r MessageRepo) GetNewMessageCounts(ctx context.Context, memberID int) (*domain.MessageCounts, error) {
	unreadQuery := `SELECT COUNT(*) FROM message_member mm 
                	LEFT JOIN message m ON m.member_id != $1 AND mm.message_id = m.id
					WHERE mm.message_id NOT IN 
					(SELECT message_id FROM message_viewer WHERE member_id = $1) 
					AND mm.member_id = $1 AND mm.deleted IS false AND mm.date_posted IS NULL`

	var unread int
	err := r.connPool.QueryRow(ctx, unreadQuery, memberID).Scan(&unread)
	if err != nil {
		return nil, err
	}

	unreadPostsQuery := `SELECT COUNT(*) FROM message_post mp
						 LEFT JOIN message_member mm ON mp.message_id = mm.message_id
						 LEFT JOIN message_viewer mv ON mv.message_id = mp.message_id
						 WHERE mm.member_id = $1 AND mv.member_id = $1 AND mp.date_posted > mv.last_viewed
						 AND mm.deleted IS false AND mp.member_id != $1`

	var newPosts int
	err = r.connPool.QueryRow(ctx, unreadPostsQuery, memberID).Scan(&newPosts)
	if err != nil {
		return nil, err
	}

	return &domain.MessageCounts{Unread: unread, NewPosts: newPosts}, nil
}

func (r MessageRepo) DeleteMessage(ctx context.Context, memberID, messageID int) error {
	_, err := r.connPool.Exec(ctx, "UPDATE message_member SET deleted = true WHERE member_id = $1 AND message_id = $2", memberID, messageID)
	if err != nil {
		return err
	}

	return nil
}

func (r MessageRepo) ViewMessage(ctx context.Context, memberID, messageID int) (int, error) {
	var views int
	err := r.connPool.QueryRow(ctx, "UPDATE message SET views = views + 1 WHERE id = $1 RETURNING views", messageID).Scan(&views)
	if err != nil {
		return 0, err
	}

	return views, nil
}

func (r MessageRepo) ListMessages(ctx context.Context, cursors domain.Cursors, limit, memberID int) ([]domain.Message, domain.Cursors, error) {
	ignoredMemberQuery := sq.Select("mi.ignore_member_id").From("member_ignore mi").LeftJoin("member m ON m.id = mi.member_id").Where("mi.member_id = ?", memberID).OrderBy("m.name")

	totalMinusIgnoredQuery := sq.Select("COUNT(*)").From("message_member mm").
		InnerJoin("message mg ON mg.id = mm.message_id").
		Where("mm.deleted IS false AND mm.member_id = ?", memberID).
		Where(SubQueryNOTIN("mm.message_id", ignoredMemberQuery))

	innerQuery := sq.Select("mg.id", "mg.member_id", "mb.name", "mg.subject", "mp.body", "mg.date_posted", "mg.posts", "mg.views", "ml.name AS last_poster_name", "mg.date_last_posted").
		From("message_member mm").
		LeftJoin("message mg ON mg.id = mm.message_id").LeftJoin("message_post mp ON mp.id = mg.first_post_id").
		LeftJoin("member mb ON mb.id = mg.member_id").
		LeftJoin("member ml ON ml.id = mg.last_member_id").
		Where("mm.deleted IS false AND mm.member_id = ?", memberID)

	rowsLeftQuery := totalMinusIgnoredQuery
	pagination := innerQuery

	if cursors.Next != "" && cursors.Prev != "" {
		return nil, domain.Cursors{}, errors.New("two cursors cannot be provided at the same time")
	}

	// Going forward
	if cursors.Next != "" {
		rowsLeftQuery = rowsLeftQuery.Where("mg.date_last_posted < ?", cursors.Next)
		pagination = pagination.Where("mg.date_last_posted < ?", cursors.Next).OrderBy("date_last_posted DESC").Limit(uint64(limit))
	}

	// Going backward
	if cursors.Prev != "" {
		rowsLeftQuery = rowsLeftQuery.Where("mg.date_last_posted > ?", cursors.Prev)
		pagination = pagination.Where("mg.date_last_posted > ?", cursors.Prev).OrderBy("date_last_posted ASC").Limit(uint64(limit))
	}

	// No cursors: Going forward from the beginning
	if cursors.Next == "" && cursors.Prev == "" {
		pagination = pagination.OrderBy("mg.date_last_posted DESC").Limit(uint64(limit))
	}

	stmt := sq.Select("id", "member_id", "name", "subject", "body", "date_posted", "posts", "views", "last_poster_name", "date_last_posted").
		Column(sq.Alias(rowsLeftQuery, "rows_left")).
		Column(sq.Alias(totalMinusIgnoredQuery, "total")).
		FromSelect(pagination, "p").OrderBy("date_last_posted DESC").PlaceholderFormat(sq.Dollar)

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
		messages []domain.Message
		rowsLeft int
		total    int
	)

	for rows.Next() {
		var message domain.Message
		err = rows.Scan(
			&message.ID,
			&message.MemberID,
			&message.MemberName,
			&message.Subject,
			&message.Body,
			&message.DatePosted,
			&message.NumPosts,
			&message.Views,
			&message.LastPosterName,
			&message.DateLastPosted,
			&rowsLeft,
			&total,
		)
		if err != nil {
			return nil, domain.Cursors{}, err
		}

		messages = append(messages, message)
	}

	// if len(messages) == 0 {
	// 	return nil, domain.Cursors{}, nil
	// }

	var (
		prevCursor string // cursor we return when there is a previous page
		nextCursor string // cursor we return when there is a next page
	)

	switch {

	// *If there are no results we don't have to compute the cursors
	case rowsLeft <= 0:

	// *On A, direction A->E (going forward), return only next cursor
	case cursors.Prev == "" && cursors.Next == "":
		if rowsLeft == len(messages) {
			return messages, domain.Cursors{}, nil
		}

		nextCursor = messages[len(messages)-1].DateLastPosted.UTC().Format(time.RFC3339Nano)

	// *On E, direction A->E (going forward), return only prev cursor
	case cursors.Next != "" && rowsLeft == len(messages):
		prevCursor = messages[0].DateLastPosted.UTC().Format(time.RFC3339Nano)

	// *On A, direction E->A (going backward), return only next cursor
	case cursors.Prev != "" && rowsLeft == len(messages):
		nextCursor = messages[len(messages)-1].DateLastPosted.UTC().Format(time.RFC3339Nano)

	// *On E, direction E->A (going backward), return only prev cursor
	case cursors.Prev != "" && total == rowsLeft:
		prevCursor = messages[0].DateLastPosted.UTC().Format(time.RFC3339Nano)

	// *Somewhere in the middle
	default:
		nextCursor = messages[len(messages)-1].DateLastPosted.UTC().Format(time.RFC3339Nano)
		prevCursor = messages[0].DateLastPosted.UTC().Format(time.RFC3339Nano)
	}

	return messages, domain.Cursors{Next: nextCursor, Prev: prevCursor}, nil
}

func (r MessageRepo) SaveMessage(message domain.Message) (int, error) {
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

	var messageID int
	err = tx.QueryRow(context.Background(), insertMessage, message.MemberID, message.Subject, message.MemberID).Scan(&messageID)
	if err != nil {
		return 0, err
	}

	var postID int
	err = tx.QueryRow(context.Background(), insertMessagePost, messageID, message.MemberID, message.Body, message.MemberIP).Scan(&postID)
	if err != nil {
		return 0, err
	}

	for _, recipientID := range message.RecipientIDs {
		_, err = tx.Exec(context.Background(), insertMessageMember, messageID, recipientID)
		if err != nil {
			return 0, err
		}
	}

	return messageID, nil
}

func (r MessageRepo) GetMessageByID(ctx context.Context, messageID, memberID int) (*domain.Message, error) {
	mvQuery := "INSERT INTO message_viewer (message_id, member_id) VALUES ($1, $2) ON CONFLICT (message_id, member_id) DO UPDATE SET last_viewed = NOW()"

	_, err := r.connPool.Exec(ctx, mvQuery, messageID, memberID)
	if err != nil {
		return nil, err
	}

	var message domain.Message
	err = r.connPool.QueryRow(ctx, getMessageByID, memberID, messageID).Scan(
		&message.ID,
		&message.DateLastPosted,
		&message.MemberID,
		&message.MemberName,
		&message.LastPosterID,
		&message.LastPosterName,
		&message.Subject,
		&message.NumPosts,
		&message.Views,
		&message.Body,
	)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r MessageRepo) GetMessagePostsCollapsible(ctx context.Context, viewable, messageID, memberID int) ([]domain.MessagePost, int, error) {
	rows, err := r.connPool.Query(ctx, listPostsCollapsible, messageID, memberID, viewable)
	if err != nil {
		return nil, 0, err
	}

	var posts []domain.MessagePost
	var collapsed int
	var ip pgtype.CIDR
	for rows.Next() {
		var post domain.MessagePost
		err = rows.Scan(
			&post.ID,
			&post.Timestamp,
			&post.MemberID,
			&post.MemberName,
			&post.Body,
			&ip,
			&post.ParentSubject,
			&post.ParentID,
			&post.Position,
			&collapsed,
		)
		if err != nil {
			return nil, 0, err
		}

		post.MemberIP = ip.IPNet.String()

		posts = append(posts, post)
	}

	return posts, collapsed, nil
}

func (r MessageRepo) GetMessageParticipants(ctx context.Context, messageID int) ([]string, error) {
	var participants []string
	rows, err := r.connPool.Query(ctx, "SELECT m.name FROM message_member mm JOIN member m ON m.id = mm.member_id WHERE mm.message_id = $1", messageID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		participants = append(participants, name)
	}

	return participants, nil
}

func (r MessageRepo) GetMessagesWithCursorForward(memberID, limit int, cursor *time.Time) ([]domain.Message, error) {
	rows, err := r.connPool.Query(context.Background(), getMessagesWithCursorForward, memberID, cursor, limit)
	if err != nil {
		return nil, err
	}

	var messages []domain.Message
	for rows.Next() {
		var message domain.Message
		err = rows.Scan(
			&message.ID,
			&message.DateLastPosted,
			&message.MemberID,
			&message.MemberName,
			&message.LastPosterID,
			&message.LastPosterName,
			&message.Subject,
			&message.NumPosts,
			&message.Views,
			&message.Body,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r MessageRepo) GetMessagesWithCursorReverse(memberID, limit int, cursor *time.Time) ([]domain.Message, error) {
	rows, err := r.connPool.Query(context.Background(), getMessagesWithCursorReverse, memberID, cursor, limit)
	if err != nil {
		return nil, err
	}

	var messages []domain.Message
	for rows.Next() {
		var message domain.Message
		err = rows.Scan(
			&message.ID,
			&message.DateLastPosted,
			&message.MemberID,
			&message.MemberName,
			&message.LastPosterID,
			&message.LastPosterName,
			&message.Subject,
			&message.NumPosts,
			&message.Views,
			&message.Body,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r MessageRepo) PeekPrevious(timestamp *time.Time) (bool, error) {
	var exists bool
	err := r.connPool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM message WHERE date_last_posted > $1)", timestamp).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r MessageRepo) GetMessagePostsByID(memberID, messageID, limit int) ([]domain.MessagePost, error) {
	var posts []domain.MessagePost
	rows, err := r.connPool.Query(context.Background(), listMessagePostsCursorForward, messageID, memberID, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post domain.MessagePost
		var ip pgtype.CIDR
		err = rows.Scan(
			&post.ID,
			&post.Timestamp,
			&post.MemberID,
			&post.MemberName,
			&post.Body,
			&ip,
			&post.ParentSubject,
			&post.ParentID,
		)
		if err != nil {
			return nil, err
		}

		post.MemberIP = ip.IPNet.String()

		posts = append(posts, post)
	}

	return posts, nil
}

func (r MessageRepo) SavePost(post domain.MessagePost) (int, error) {
	var postID int
	err := r.connPool.QueryRow(context.Background(), insertMessagePost, post.ParentID, post.MemberID, post.Body, post.MemberIP).Scan(&postID)
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (r MessageRepo) GetMessagePostByID(id int) (*domain.MessagePost, error) {
	var messagePost domain.MessagePost
	var ip pgtype.CIDR
	err := r.connPool.QueryRow(context.Background(), getMessagePostByID, id).Scan(
		&messagePost.ID,
		&messagePost.Timestamp,
		&messagePost.MemberID,
		&messagePost.MemberName,
		&messagePost.Body,
		&ip,
		&messagePost.ParentSubject,
		&messagePost.ParentID,
	)
	if err != nil {
		return nil, err
	}

	messagePost.MemberIP = ip.IPNet.String()

	return &messagePost, nil
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
