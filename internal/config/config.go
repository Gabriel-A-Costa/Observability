package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port string `envconfig:"APP_PORT" default:"8080"`
}

func Load() (Config, error) {
	godotenv.Load()

	var cfg Config
	err := envconfig.Process("", &cfg)
	return cfg, err
}
