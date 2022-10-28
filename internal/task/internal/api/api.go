package api

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-uuid"
	"github.com/nikp359/ates/internal/task/internal/model"

	"github.com/nikp359/ates/internal/task/internal/repository"

	"github.com/nikp359/ates/internal/estream"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	srv            *echo.Echo
	taskRepository *repository.TaskRepository
	producer       Producer
}

type Producer interface {
	Send(eventName string, payload estream.Payload) error
}

type BadRequest struct {
	Message string `json:"message"`
}

func NewServer(taskRepository *repository.TaskRepository, producer Producer) *Server {
	s := &Server{
		taskRepository: taskRepository,
		producer:       producer,
	}

	s.srv = s.routers()

	return s
}

func (s *Server) Start() error {
	return s.srv.Start(":8081") // TODO: move to config
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

	api.GET("/tasks", s.taskList)
	api.POST("/tasks", s.createTask)
	api.POST("/tasks/shuffle", s.shuffledTask)
	api.PUT("/tasks/complete", s.completeTask)

	return e
}

func (s *Server) taskList(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (s *Server) createTask(c echo.Context) error {
	var task model.Task

	if err := c.Bind(&task); err != nil {
		return err
	}

	uid, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	task.Status = repository.TaskStatusNew
	task.PublicID = uid

	if err = s.taskRepository.Add(&task); err != nil {
		return err
	}

	storedTask, err := s.taskRepository.GetByPublicID(task.PublicID)
	if err != nil {
		return err
	}

	eventTaskCreated := estream.TaskCreatedPayload(storedTask)
	if err = s.producer.Send(estream.TaskCreated, &eventTaskCreated); err != nil {
		return err
	}

	eventTaskRegistered := estream.TaskRegisteredPayload(storedTask)
	if err = s.producer.Send(estream.TaskRegistered, &eventTaskRegistered); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, storedTask)
}

func (s *Server) shuffledTask(c echo.Context) error {
	tasks, err := s.taskRepository.Shuffled(c.Request().Context())
	if err != nil {
		return err
	}

	eventTasksShuffled := make(estream.TaskShuffledPayload, 0, len(tasks))
	for _, task := range tasks {
		eventTaskUpdated := estream.TaskRegisteredPayload(task)

		if err = s.producer.Send(estream.TaskUpdated, &eventTaskUpdated); err != nil {
			return err
		}

		eventTasksShuffled = append(eventTasksShuffled, estream.TaskAssigned{
			UserID: task.AssignedUserID,
			TaskID: task.PublicID,
		})
	}

	if err = s.producer.Send(estream.TasksShuffled, &eventTasksShuffled); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TasksShuffledResponse{Shuffled: len(tasks)})
}

func (s *Server) completeTask(c echo.Context) error {
	var req TaskCompletedRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.TaskID == "" {
		return c.JSON(http.StatusBadRequest, BadRequest{Message: "task_id can not be empty"})
	}

	task, err := s.taskRepository.ChangeStatus(c.Request().Context(), req.TaskID, repository.TaskStatusCompleted)
	if err != nil {
		return err
	}

	eventTaskUpdated := estream.TaskUpdatedPayload(task)
	if err = s.producer.Send(estream.TaskUpdated, &eventTaskUpdated); err != nil {
		return err
	}

	// TODO: get userID from jwt session
	eventTaskCompleted := estream.TaskCompletedPayload{
		PublicID:  task.PublicID,
		UserID:    task.AssignedUserID,
		UpdatedAt: task.UpdatedAt,
	}
	if err = s.producer.Send(estream.TaskCompeted, &eventTaskCompleted); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, task)
}
