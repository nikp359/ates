package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nikp359/ates/internal/auth/internal/model"

	"github.com/nikp359/ates/internal/api"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/nikp359/ates/internal/auth/internal/repository"
)

type App struct {
	userRepository *repository.UserRepository
}

func NewApp(config *Config) *App {
	return &App{
		userRepository: repository.NewUserRepository(newDB(config.DB.Connection)),
	}
}

func newDB(dataSource string) *sqlx.DB {
	db, err := sqlx.Connect("mysql", dataSource)
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
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

func Start() {
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