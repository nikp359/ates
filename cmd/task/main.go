package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/nikp359/ates/internal/task"
)

func main() {
	configPath := parseFlags()
	config, err := task.NewConfig(configPath)
	if err != nil {
		logrus.WithError(err).Fatal("parse config file")
	}

	app, err := task.NewApp(config)
	if err != nil {
		logrus.WithError(err).Fatal("init App")
	}

	app.Start()
}

func parseFlags() string {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")
	flag.Parse()

	return configPath
}
