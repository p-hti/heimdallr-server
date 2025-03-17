package service

import "log/slog"

type Service struct {
	Storage Storage
	PasswordHasher
	Logger     *slog.Logger
	hmacSecret []byte
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Storage interface {
	AuthStorage
}

func NewService(storage Storage, logger *slog.Logger, secret []byte) *Service {
	return &Service{
		Storage:    storage,
		Logger:     logger,
		hmacSecret: secret,
	}
}
