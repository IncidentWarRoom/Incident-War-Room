package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

type Config struct {
	BotToken string `env:"BOT_TOKEN" env-required:"true"`

	PostgresHost     string `env:"POSTGRES_HOST"     env-default:"localhost"`
	PostgresPort     int    `env:"POSTGRES_PORT"     env-default:"5432"`
	PostgresDB       string `env:"POSTGRES_DB"       env-required:"true"`
	PostgresUser     string `env:"POSTGRES_USER"     env-required:"true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
}

func Load() (*Config, error) {
	_ = godotenv.Load() // load .env if present, ignore error in production

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, errs.Wrapf(errs.KindInternal, "config.Load", err, "read environment")
	}
	return &cfg, nil
}
