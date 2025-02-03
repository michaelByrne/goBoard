package domain

import "time"

type Token struct {
	TokenStr string
	Expires  time.Time
}
