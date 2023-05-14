package member

import (
	"goBoard/internal/core/domain"
	"time"
)

type Member struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Pass        string `json:"pass"`
	Secret      string `json:"secret"`
	PostalCode  string `json:"postal_code"`
	Banned      bool   `json:"banned"`
	FirstPosted string `json:"first_posted"`
	LastPosted  string `json:"last_posted"`
	IP          string `json:"ip"`
}

func (m *Member) FromDomain(member domain.Member) {
	m.ID = member.ID
	m.Name = member.Name
	m.Email = member.Email
	m.Pass = member.Pass
	m.Secret = member.Secret
	m.PostalCode = member.PostalCode
	m.Banned = member.Banned

	if member.FirstPosted != nil {
		m.FirstPosted = member.FirstPosted.Format(time.RFC3339)
	}

	if member.LastPosted != nil {
		m.LastPosted = member.LastPosted.Format(time.RFC3339)
	}

	m.IP = member.IP
}

func (m *Member) ToDomain() domain.Member {
	firstPosted, err := time.Parse(time.RFC3339, m.FirstPosted)
	if err != nil {
		firstPosted = time.Now()
	}

	lastPosted, err := time.Parse(time.RFC3339, m.LastPosted)
	if err != nil {
		lastPosted = time.Now()
	}

	return domain.Member{
		ID:          m.ID,
		Name:        m.Name,
		Email:       m.Email,
		Pass:        m.Pass,
		Secret:      m.Secret,
		PostalCode:  m.PostalCode,
		Banned:      m.Banned,
		FirstPosted: &firstPosted,
		LastPosted:  &lastPosted,
		IP:          m.IP,
	}
}
