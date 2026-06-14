package config

import (
	"fmt"
	"net/url"

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

	ReportServiceURL string `env:"REPORT_SERVICE_URL" env-default:"http://localhost:8000"`
}

func Load() (*Config, error) {
	_ = godotenv.Load() // load .env if present, ignore error in production

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, errs.Wrapf(errs.KindInternal, "config.Load", err, "read environment")
	}
	return &cfg, nil
}

// PostgresDSN builds a connection URL for the configured database.
func (c *Config) PostgresDSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.PostgresUser, c.PostgresPassword),
		Host:   fmt.Sprintf("%s:%d", c.PostgresHost, c.PostgresPort),
		Path:   c.PostgresDB,
	}
	return u.String()
}
