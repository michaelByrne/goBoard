package domain

import (
	"html/template"
	"time"
)

type ThreadPost struct {
	ID            int
	Timestamp     *time.Time
	MemberID      int
	MemberIP      string
	ParentID      int
	Body          string
	HtmlBody      *template.HTML
	ParentSubject string
	IsAdmin       bool
	MemberName    string
	Position      int
	Collapsed     int
	RowNumber     int
}
