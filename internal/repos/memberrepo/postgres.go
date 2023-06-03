package memberrepo

import (
	"context"
	"goBoard/internal/core/domain"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	saveMemberQuery = "INSERT INTO member (name, pass, secret, email_signup, postalcode, ip) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	getMemberPrefs  = "SELECT p.name, mp.value FROM member_pref mp JOIN pref p ON mp.pref_id = p.id WHERE mp.member_id = $1"
)

type MemberRepo struct {
	connPool *pgxpool.Pool
}

func NewMemberRepo(connPool *pgxpool.Pool) *MemberRepo {
	return &MemberRepo{connPool: connPool}
}

func (m MemberRepo) SaveMember(member domain.Member) (int, error) {
	var id int
	err := m.connPool.QueryRow(context.Background(), saveMemberQuery, member.Name, member.Pass, member.Secret, member.Email, member.PostalCode, member.IP).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m MemberRepo) GetMemberByID(id int) (*domain.Member, error) {
	var member domain.Member
	var ip pgtype.CIDR
	err := m.connPool.QueryRow(context.Background(), "SELECT id, name, pass, secret, email_signup, postalcode, ip FROM member WHERE id = $1", id).Scan(&member.ID, &member.Name, &member.Pass, &member.Secret, &member.Email, &member.PostalCode, &ip)
	if err != nil {
		return nil, err
	}

	member.IP = ip.IPNet.String()

	return &member, nil
}

func (m MemberRepo) GetMemberIDByUsername(username string) (int, error) {
	var id int
	err := m.connPool.QueryRow(context.Background(), "SELECT id FROM member WHERE name = $1", username).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m MemberRepo) GetMemberByUsername(username string) (*domain.Member, error) {
	var member domain.Member
	var ip pgtype.CIDR
	err := m.connPool.QueryRow(
		context.Background(),
		"SELECT id, name, pass, secret, email_signup, postalcode, ip, date_joined, date_first_post, last_post, last_view, total_threads, total_thread_posts, is_admin FROM member WHERE name = $1",
		username,
	).Scan(
		&member.ID,
		&member.Name,
		&member.Pass,
		&member.Secret,
		&member.Email,
		&member.PostalCode,
		&ip,
		&member.DateJoined,
		&member.FirstPosted,
		&member.LastPosted,
		&member.LastView,
		&member.TotalThreads,
		&member.TotalThreadPosts,
		&member.IsAdmin,
	)
	if err != nil {
		return nil, err
	}

	member.IP = ip.IPNet.String()

	return &member, nil
}

func (m MemberRepo) GetMemberPrefs(memberID int) (*domain.MemberPrefs, error) {
	rows, err := m.connPool.Query(context.Background(), getMemberPrefs, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prefs := make(domain.MemberPrefs)
	for rows.Next() {
		var name, value string
		err := rows.Scan(&name, &value)
		if err != nil {
			return nil, err
		}

		prefs[name] = value
	}

	return &prefs, nil
}
