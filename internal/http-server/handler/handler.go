package handler

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/p-hti/heimdallr-server/internal/http-server/middleware/logger"
)

type HTTPServer struct {
	Service Service
	Logger  *slog.Logger
	JwtKey  []byte
}

func NewHTTPServer(service Service, logger *slog.Logger, jwtKey []byte) *HTTPServer {
	return &HTTPServer{
		Service: service,
		Logger:  logger,
		JwtKey:  jwtKey,
	}
}

type Service interface {
	AuthService
}

func (h *HTTPServer) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(logger.New(h.Logger))
	router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", h.SignUp)
		r.Get("/sign-in", h.SignIn)
		r.Get("/log-out", h.LogOut)
		r.Get("/terminate-sessions", h.FullLogOut)
		r.Get("/refresh", h.RefreshToken)
	})

	return router
}
