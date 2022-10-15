package api

import (
	"context"
	"net/http"

	"github.com/nikp359/ates/internal/auth/internal/repository"
	"github.com/nikp359/ates/internal/estream"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	srv            *echo.Echo
	userRepository *repository.UserRepository
	producer       *estream.Producer
}

func NewServer(userRepository *repository.UserRepository, producer *estream.Producer) *Server {
	s := &Server{
		userRepository: userRepository,
		producer:       producer,
	}

	s.srv = s.routers()

	return s
}

func (s *Server) Start() error {
	return s.srv.Start(":8080")
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) routers() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	return e
}