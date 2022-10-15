package auth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/nikp359/ates/internal/auth/internal/api"
	"github.com/nikp359/ates/internal/auth/internal/model"
	"github.com/nikp359/ates/internal/auth/internal/repository"
	"github.com/nikp359/ates/internal/estream"
)

type App struct {
	userRepository *repository.UserRepository
	producer       *estream.Producer
}

func NewApp(config *Config) (*App, error) {
	db, err := newDB(config.DB.Connection)
	if err != nil {
		return nil, err
	}

	producer, err := estream.NewSyncProducer(estream.Config{
		Addresses: []string{"localhost:29092"},
	})
	if err != nil {
		return nil, err
	}

	return &App{
		userRepository: repository.NewUserRepository(db),
		producer:       producer,
	}, nil
}

func newDB(dataSource string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func (a *App) AddUser() {
	err := a.userRepository.AddUser(&model.User{
		PublicID: "abc-123",
		Email:    "test@example.com",
		Role:     "manager",
	})

	if err != nil {
		logrus.WithError(err).Errorf("add user")
	}
}

func (a *App) SendMsg() {
	err := a.producer.SendSync(estream.UserCreated, &estream.UserEvent{
		PublicID: "abc-123",
		Email:    "",
	})

	logrus.Infof("error: %v", err)
}

func (a *App) Start() {
	srv := api.NewServer()

	// Start server
	go func(s *api.Server) {
		if err := s.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatal(s.Start())
		}
	}(srv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server is stopped")
}
