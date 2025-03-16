package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/p-hti/heimdallar-server/internal/domain/model"
	"github.com/p-hti/heimdallar-server/internal/http-server/middleware/logger"
)

type HTTPServer struct {
	Service Service
	Logger  *slog.Logger
}

func NewHTTPServer(service Service, logger *slog.Logger) *HTTPServer {
	return &HTTPServer{
		Service: service,
		Logger:  logger,
	}
}

type Service interface {
	AuthService
}

type AuthService interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (
		uid int64,
		err error,
	)
	GetUser(
		ctx context.Context,
		email string,
	) (
		model.User,
		error,
	)
	SaveToken(
		token string,
		uid int64,
		createdTime time.Time,
		expiresTime time.Time,
	) (
		int64,
		error,
	)
	GetToken(
		uid int64,
	) (
		string,
		error,
	)
}

func (h *HTTPServer) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(logger.New(h.Logger))
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})

	return router
}
