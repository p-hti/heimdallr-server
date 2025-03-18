package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/p-hti/heimdallr-server/internal/config"
)

type Postgres struct {
	Db  *pgxpool.Pool
	log *slog.Logger
}

func NewPostgres(cfg *config.Postgres, log *slog.Logger) (*Postgres, error) {
	conf := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBPass, cfg.SSLMode)

	poolConfig, err := pgxpool.ParseConfig(conf)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(
		context.Background(),
		poolConfig,
	)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Postgres{
		Db:  db,
		log: log,
	}, nil
}
