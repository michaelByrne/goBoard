package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strings"
)

//go:embed migrations
var migrations embed.FS

func SetupTestDatabase() (testcontainers.Container, *pgxpool.Pool, error) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "board",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_USER":     "board",
		},
	}
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, nil, err
	}
	port, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		return nil, nil, err
	}
	host, err := dbContainer.Host(context.Background())
	if err != nil {
		return nil, nil, err
	}

	dbURI := fmt.Sprintf("postgres://board:test@%v:%v/board", host, port.Port())
	err = MigrateDBUp(dbURI)
	if err != nil {
		return nil, nil, err
	}

	connPool, err := pgxpool.Connect(context.Background(), dbURI)
	if err != nil {
		return nil, nil, err
	}

	return dbContainer, connPool, err
}

func MigrateDBUp(dbURI string) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, strings.Replace(dbURI, "postgres://", "pgx://", 1))
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func MigrateDBDown(dbURI string) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, strings.Replace(dbURI, "postgres://", "pgx://", 1))
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
