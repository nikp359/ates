package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/nikp359/ates/internal/auth"
)

func main() {
	configPath := parseFlags()
	config, err := auth.NewConfig(configPath)
	if err != nil {
		logrus.WithError(err).Fatal("parse config file")
	}

	app := auth.NewApp(config)

	app.AddUser()
}

func parseFlags() string {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")
	flag.Parse()

	return configPath
}
