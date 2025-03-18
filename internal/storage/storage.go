package storage

import (
	"log/slog"

	"github.com/p-hti/heimdallr-server/internal/config"
	"github.com/p-hti/heimdallr-server/internal/storage/postgres"
)

type Storage struct {
	Postgres *postgres.Postgres
	log      *slog.Logger
}

func NewStorage(cfg *config.Config, log *slog.Logger) (*Storage, error) {
	postgresDb, err := postgres.NewPostgres(&cfg.Postgres, log)
	if err != nil {
		return nil, err
	}
	return &Storage{
		Postgres: postgresDb,
		log:      log,
	}, nil
}
