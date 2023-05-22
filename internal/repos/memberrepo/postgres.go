package memberrepo

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"goBoard/internal/core/domain"
)

const saveMemberQuery = "INSERT INTO member (name, pass, secret, email_signup, postalcode, ip) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

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
	err := m.connPool.QueryRow(context.Background(), "SELECT id, name, pass, secret, email_signup, postalcode, ip, date_joined, date_first_post, last_post, last_view, total_threads, total_thread_posts FROM member WHERE name = $1", username).Scan(&member.ID, &member.Name, &member.Pass, &member.Secret, &member.Email, &member.PostalCode, &ip, &member.DateJoined, &member.FirstPosted, &member.LastPosted, &member.LastView, &member.TotalThreads, &member.TotalThreadPosts)
	if err != nil {
		return nil, err
	}

	member.IP = ip.IPNet.String()

	return &member, nil
}