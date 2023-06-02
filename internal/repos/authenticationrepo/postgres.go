package authenticationrepo

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	AuthenticateQuery = `SELECT id FROM member WHERE lower(name) = lower($1) AND pass = $2 AND banned = false`
)

type AuthenticationRepo struct {
	connPool *pgxpool.Pool
}

func NewAuthenticationRepo(connPool *pgxpool.Pool) *AuthenticationRepo {
	return &AuthenticationRepo{
		connPool: connPool,
	}
}

func (r *AuthenticationRepo) Authenticate(username, password string) (int, error) {
	var id int
	err := r.connPool.QueryRow(context.Background(), AuthenticateQuery, username, password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
