package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB struct {
		DB_PASS string `yaml:"DB_PASS"`
		DB_USER string `yaml:"DB_USER"`
		DB_NAME string `yaml:"DB_NAME"`
		DB_PORT int    `yaml:"DB_PORT"`
		DB_HOST string `yaml:"DB_HOST"`
	} `yaml:"db"`
}

func LoadConfig(log *zap.SugaredLogger) (*Config, error) {
	data, err := os.ReadFile("app/cmd/config/config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
