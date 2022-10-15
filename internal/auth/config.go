package auth

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		DB ConfigDB `yaml:"db"`
	}

	ConfigDB struct {
		Connection string `yaml:"connection"`
	}

	Kafka struct {
	}
)

func NewConfig(configPath string) (*Config, error) {
	if err := validateConfigPath(configPath); err != nil {
		return nil, err
	}

	config := &Config{}

	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return nil, err
	}
	defer func() {
		if fErr := file.Close(); fErr != nil {
			logrus.WithError(err).Errorf("file close")
		}
	}()

	d := yaml.NewDecoder(file)
	if err = d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
