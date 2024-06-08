package config

import (
	"github.com/gookit/slog"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		PG
	}

	PG struct {
		URL string ` env:"PG_URL_LOCALHOST"`
	}
)

func NewConfig() *Config {
	cfg := &Config{}
	err := godotenv.Load()
	if err != nil {
		slog.Fatal("can't load env %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		slog.Fatal("error reading env %w", err)
	}
	return cfg
}
