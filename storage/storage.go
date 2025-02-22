package storage

import (
	"context"
	"database/sql"
	"embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/rsmanito/developstoday-test-assessment/config"
	"github.com/rsmanito/developstoday-test-assessment/storage/postgres"
)

type Storage struct {
	*postgres.Queries
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (s *Storage) Migrate(cfg *config.Config) {
	slog.Info("Migrating database")

	db, err := sql.Open(
		"postgres",
		cfg.DbConnUrl,
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("pgx"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	slog.Info("Database migrated")
}

func New(cfg *config.Config) *Storage {
	conn, err := pgx.Connect(
		context.Background(),
		cfg.DbConnUrl,
	)
	if err != nil {
		panic(err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	queries := postgres.New(conn)

	st := &Storage{queries}
	st.Migrate(cfg)

	return st
}
