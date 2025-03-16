package service

import (
	"context"

	"github.com/p-hti/heimdallar-server/internal/domain/model"
)

type Service struct {
	Storage Storage
}

type Storage interface {
	AuthStorage
}

type AuthStorage interface {
	GetUser(
		ctx context.Context,
		email string,
	) (
		model.User,
		error,
	)
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (
		int64,
		error,
	)
}
