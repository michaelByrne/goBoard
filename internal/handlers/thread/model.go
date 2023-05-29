package thread

import (
	"goBoard/internal/core/domain"
	"time"
)

type ID struct {
	ID int `json:"id"`
}

type Thread struct {
	ID             int        `json:"id"`
	Subject        string     `json:"subject"`
	MemberID       int        `json:"member_id"`
	MemberName     string     `json:"member_name"`
	MemberIP       string     `json:"member_ip"`
	Views          int        `json:"views"`
	LastPosterName string     `json:"last_poster_name"`
	LastPosterID   int        `json:"last_poster_id"`
	LastPostText   string     `json:"last_post_text"`
	DateLastPosted string     `json:"date_last_posted"`
	Sticky         bool       `json:"sticky"`
	Locked         bool       `json:"locked"`
	Legendary      bool       `json:"legendary"`
	NumPosts       int        `json:"num_posts"`
	Posts          []Post     `json:"posts"`
	Timestamp      *time.Time `json:"timestamp"`
	FirstPostText  string     `json:"first_post_text"`
}

type Threads struct {
	Threads []Thread `json:"threads"`
}

func (t *Threads) FromDomain(threads []domain.Thread) {
	for _, thread := range threads {
		newThread := &Thread{}
		newThread.FromDomain(thread)
		t.Threads = append(t.Threads, *newThread)
	}
}

type Post struct {
	ID        int        `json:"id"`
	Timestamp *time.Time `json:"timestamp"`
	MemberID  int        `json:"member_id"`
	MemberIP  string     `json:"member_ip"`
	ThreadID  int        `json:"thread_id"`
	Text      string     `json:"text"`
}

func (t *Thread) FromDomain(thread domain.Thread) {
	t.ID = thread.ID
	t.Subject = thread.Subject
	t.MemberID = thread.MemberID
	t.MemberName = thread.MemberName
	t.Views = thread.Views
	t.LastPosterName = thread.LastPosterName
	t.LastPosterID = thread.LastPosterID
	t.LastPostText = thread.LastPostText
	t.DateLastPosted = thread.DateLastPosted.String()
	t.Sticky = thread.Sticky
	t.Locked = thread.Locked
	t.Legendary = thread.Legendary
	t.NumPosts = thread.NumPosts
	t.Timestamp = thread.Timestamp
	t.FirstPostText = thread.FirstPostText
	t.MemberIP = thread.MemberIP

	for _, p := range thread.Posts {
		var post Post
		post.FromDomain(p)
		t.Posts = append(t.Posts, post)
	}
}

func (p *Post) FromDomain(post domain.ThreadPost) {
	p.ID = post.ID
	p.Timestamp = post.Timestamp
	p.MemberID = post.MemberID
	p.MemberIP = post.MemberIP
	p.ThreadID = post.ParentID
	p.Text = post.Body
}

func (p *Post) ToDomain() domain.ThreadPost {
	var post domain.ThreadPost
	post.ID = p.ID
	post.Timestamp = p.Timestamp
	post.MemberID = p.MemberID
	post.MemberIP = p.MemberIP
	post.ParentID = p.ThreadID
	post.Body = p.Text
	return post
}

func (t *Thread) ToDomain() domain.Thread {
	dateLastPosted, _ := time.Parse("2006-01-02 15:04:05", t.DateLastPosted)

	var thread domain.Thread
	thread.ID = t.ID
	thread.Subject = t.Subject
	thread.MemberID = t.MemberID
	thread.MemberName = t.MemberName
	thread.Views = t.Views
	thread.LastPosterName = t.LastPosterName
	thread.LastPosterID = t.LastPosterID
	thread.LastPostText = t.LastPostText
	thread.DateLastPosted = &dateLastPosted
	thread.Sticky = t.Sticky
	thread.Locked = t.Locked
	thread.Legendary = t.Legendary
	thread.NumPosts = t.NumPosts
	thread.Timestamp = t.Timestamp
	thread.MemberIP = t.MemberIP

	for _, p := range t.Posts {
		var post domain.ThreadPost
		post.ID = p.ID
		post.Timestamp = p.Timestamp
		post.MemberID = p.MemberID
		post.MemberIP = p.MemberIP
		post.ParentID = p.ThreadID
		post.Body = p.Text
		thread.Posts = append(thread.Posts, post)
	}

	return thread
}
