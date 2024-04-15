package messagerepo

import (
	"context"
	_ "embed"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"goBoard/internal/core/domain"
	"time"
)

//go:embed queries/get_messages_with_cursor_forward.sql
var getMessagesWithCursorForward string

//go:embed queries/list_message_posts_cursor_forward.sql
var listMessagePostsCursorForward string

//go:embed queries/get_messages_with_cursor_reverse.sql
var getMessagesWithCursorReverse string

//go:embed queries/get_message_post_by_id.sql
var getMessagePostByID string

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
