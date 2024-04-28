package themerepo

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ThemeRepo struct {
	connPool *pgxpool.Pool
}

func NewThemeRepo(connPool *pgxpool.Pool) *ThemeRepo {
	return &ThemeRepo{
		connPool: connPool,
	}
}

func (r ThemeRepo) GetTheme(ctx context.Context, name string) (string, error) {
	var theme string
	err := r.connPool.QueryRow(ctx, "SELECT value FROM theme WHERE name = $1", name).Scan(&theme)
	if err != nil {
		return "", err
	}

	return theme, nil
}
