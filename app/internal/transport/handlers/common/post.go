package common

import (
	"goBoard/internal/core/domain"
)

type Post struct {
	MemberName string
	MemberID   int
	Date       string
	Body       string
	ParentID   int
	Preview    bool
	RowNumber  int
}

func ThreadToPosts(thread domain.Thread) []Post {
	posts := make([]Post, len(thread.Posts))
	for i, p := range thread.Posts {
		posts[i] = Post{
			MemberName: p.MemberName,
			Date:       p.Timestamp.Format("Mon Jan 2, 2006 03:04 pm"),
			Body:       p.Body,
			ParentID:   thread.ID,
			MemberID:   p.MemberID,
			RowNumber:  p.RowNumber + 1,
		}
	}
	return posts
}

func MessageToPosts(message domain.Message) []Post {
	posts := make([]Post, len(message.Posts))
	for i, p := range message.Posts {
		posts[i] = Post{
			MemberName: p.MemberName,
			Date:       p.Timestamp.Format("Mon Jan 2, 2006 03:04 pm"),
			Body:       p.Body,
			ParentID:   message.ID,
			MemberID:   p.MemberID,
			RowNumber:  p.Position + 1,
		}
	}
	return posts
}
