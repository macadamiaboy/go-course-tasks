package helpers

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB  DBConfig
	App AppConfig
}

type DBConfig struct {
	ConnString string
}

type AppConfig struct {
	SecretSalt string
}

func LoadConfig(logger *slog.Logger) *Config {
	if err := godotenv.Load(); err != nil {
		logger.Error("failed to load the config", "error", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PWD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	secretSalt := os.Getenv("SECRET_SALT")

	return &Config{
		DB: DBConfig{
			ConnString: dsn,
		},
		App: AppConfig{
			SecretSalt: secretSalt,
		},
	}
}
