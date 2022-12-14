package task

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nikp359/ates/internal/task/internal/repository"
	"github.com/nikp359/ates/internal/task/internal/stream"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/nikp359/ates/internal/estream"
	"github.com/nikp359/ates/internal/task/internal/api"
)

type App struct {
	server *api.Server
}

func NewApp(config *Config) (*App, error) {
	db, err := newDB(config.DB.Connection)
	if err != nil {
		return nil, err
	}

	userConsumer, err := stream.NewUserStream(repository.NewUserRepository(db), estream.Config(config.Kafka))
	if err != nil {
		return nil, err
	}
	userConsumer.Start()

	producer, err := estream.NewAsyncProducer(estream.Config(config.Kafka))
	if err != nil {
		return nil, err
	}

	return &App{
		server: api.NewServer(repository.NewTaskRepository(db), producer),
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

func (a *App) Start() {
	// Start server
	go func(s *api.Server) {
		if err := s.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatal(s.Start())
		}
	}(a.server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.server.Stop(ctx); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server is stopped")
}
