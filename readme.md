1) install Docker for Desktop 
2) install Golang Migrate (https://github.com/golang-migrate/migrate)
3) run `docker compose up --build`
4) the app will be running but you'll need to manually add some data
5) run the tests with `go test ./...`

### Seeding the database

1. Make sure the `board_db_postgres` container is running, and that you can access the database from outside the container:
```
ebernstein@ebernstein-GVH2J-MBP goBoard % psql 'postgres://boardking:test@localhost:5432/board?sslmode=disable'
psql (13.10, server 15.3 (Debian 15.3-1.pgdg110+1))
WARNING: psql major version 13, server major version 15.
         Some psql features might not work.
Type "help" for help.

board=# \q
```

2. From the root directory, run:
```
    go run db/tools/seed.go
```