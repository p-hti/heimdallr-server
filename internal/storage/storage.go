package storage

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	Db *pgxpool.Pool
}
