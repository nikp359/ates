package api

import (
	"context"
	"net/http"

	"github.com/nikp359/ates/internal/auth/internal/model"

	"github.com/hashicorp/go-uuid"

	"github.com/nikp359/ates/internal/auth/internal/repository"
	"github.com/nikp359/ates/internal/estream"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	srv            *echo.Echo
	userRepository *repository.UserRepository
	producer       Producer
}

type Producer interface {
	Send(eventName string, payload estream.Payload) error
}

type BadRequest struct {
	Message string `json:"message"`
}

func NewServer(userRepository *repository.UserRepository, producer Producer) *Server {
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

	api := e.Group("/api")

	api.GET("/users", s.userList)
	api.POST("/users", s.createUser)
	api.PUT("/users", s.updateUser)
	api.DELETE("/users", s.deleteUser)

	return e
}

func (s *Server) userList(c echo.Context) error {
	users, err := s.userRepository.List()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

func (s *Server) createUser(c echo.Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {
		return err
	}

	uid, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	user.PublicID = uid

	if err = s.userRepository.Add(&user); err != nil {
		return err
	}

	storeUser, err := s.userRepository.GetByPublicID(user.PublicID)
	if err != nil {
		return err
	}

	if err = s.producer.Send(estream.UserCreated, &estream.UserCreatedPayload{
		PublicID:  storeUser.PublicID,
		Email:     storeUser.Email,
		Role:      storeUser.Role,
		Timestamp: storeUser.UpdatedAt,
	}); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, storeUser)
}

func (s *Server) updateUser(c echo.Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {
		return err
	}

	if user.PublicID == "" {
		return c.JSON(http.StatusBadRequest, &BadRequest{Message: "public_id is empty"})
	}

	if err := s.userRepository.Update(&user); err != nil {
		return err
	}

	storeUser, err := s.userRepository.GetByPublicID(user.PublicID)
	if err != nil {
		return err
	}

	if err = s.producer.Send(estream.UserUpdated, &estream.UserUpdatedPayload{
		PublicID:  storeUser.PublicID,
		Email:     storeUser.Email,
		Role:      storeUser.Role,
		Timestamp: storeUser.UpdatedAt,
	}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, storeUser)
}

func (s *Server) deleteUser(c echo.Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {
		return err
	}

	if user.PublicID == "" {
		return c.JSON(http.StatusBadRequest, &BadRequest{Message: "public_id is empty"})
	}

	if err := s.userRepository.Delete(user.PublicID); err != nil {
		return err
	}

	if err := s.producer.Send(estream.UserDeleted, &estream.UserDeletedPayload{
		PublicID: user.PublicID,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
