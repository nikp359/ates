package main

import (
	"context"
	"github.com/nikp359/ates/internal/api"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	srv := api.NewServer()

	// Start server
	go func(s *api.Server) {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
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
