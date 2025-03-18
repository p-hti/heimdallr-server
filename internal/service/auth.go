package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/p-hti/heimdallr-server/internal/domain/model"
)

type AuthStorage interface {
	SaveUser(
		ctx context.Context,
		input model.User,
	) error

	GetUser(
		ctx context.Context,
		email string,
		hashPass string,
	) (
		model.User,
		error,
	)

	SessionStorage
}

type SessionStorage interface {
	CreateSession(
		ctx context.Context,
		refreshToken model.RefreshToken,
	) error
	GetSession(
		ctx context.Context,
		token string,
	) (
		model.RefreshToken,
		error,
	)
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
	return s.Storage.SaveUser(ctx, user)
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

	user, err := s.Storage.GetUser(ctx, inp.Email, password)
	if err != nil {
		return "", "", nil
	}

	return s.generateToken(ctx, user.ID)
}

func (s *Service) generateToken(
	ctx context.Context,
	uid int64,
) (
	string,
	string,
	error,
) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(uid)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute + 15).Unix(),
	})

	accessToken, err := t.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.Storage.CreateSession(
		ctx,
		model.RefreshToken{
			UserID:    uid,
			ExpiresAt: time.Now().Add(time.Hour * 24),
			Token:     refreshToken,
		},
	); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *Service) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (
	string,
	string,
	error,
) {
	session, err := s.Storage.GetSession(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", errors.New("error refresh token") // # TODO: package for error
	}

	return s.generateToken(ctx, session.ID)
}
