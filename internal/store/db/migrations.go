package db

import (
	"database/sql"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

//go:embed sqlc/migrations/*.sql
var migrations embed.FS

func RunMigrations(pool *pgxpool.Pool) error {
	goose.SetBaseFS(migrations)
	goose.WithLogger(goose.NopLogger())

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	cp := pool.Config().ConnConfig.ConnString()
	db, err := sql.Open("pgx/v5", cp)
	if err != nil {
		return err
	}

	err = goose.Up(db, "sqlc/migrations")
	return err
}
