package service

import (
	"context"
	"time"

	"github.com/p-hti/heimdallr-server/internal/domain/model"
)

type AuthStorage interface {
	SignIn(
		ctx context.Context,
		email string,
	) (
		model.User,
		error,
	)
	SignUp(
		ctx context.Context,
		inputStruct model.User,
	) error

	SessionStorage
}

type SessionStorage interface {
	TerminateSession(
		ctx context.Context,
		refreshToken string,
	) error

	TerminateAllSessions(
		ctx context.Context,
		uid int64,
	) error
}

func (s *Service) SaveUser(
	ctx context.Context,
	inputStruct model.SingUpStruct,
) error {
	password, err := s.Hash(inputStruct.Password)
	if err != nil {
		return err
	}

	user := model.User{
		Name:      inputStruct.Name,
		Email:     inputStruct.Email,
		PassHash:  []byte(password),
		CreatedAt: time.Now(),
	}
	return s.Storage.SignUp(ctx, user)
}

func (s *Service) SignIn(
	ctx context.Context,
	inp model.SignInStruct,
) (
	string,
	string,
	error,
) {
	password, err := s.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	return "", "", nil
}
