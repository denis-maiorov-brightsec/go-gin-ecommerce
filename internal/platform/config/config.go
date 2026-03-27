package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                 string
	Port                   string
	LogLevel               string
	DatabaseURL            string
	WriteRateLimitRequests int
	WriteRateLimitWindow   time.Duration
}

const (
	defaultWriteRateLimitWindow       = time.Minute
	defaultTestWriteRateLimitRequests = 1000
	defaultWriteRateLimitRequests     = 60
)

func Load() (Config, error) {
	_ = godotenv.Load(".env.local")

	writeRateLimitRequests, err := getEnvInt("WRITE_RATE_LIMIT_REQUESTS")
	if err != nil {
		return Config{}, err
	}

	writeRateLimitWindow, err := getEnvDuration("WRITE_RATE_LIMIT_WINDOW")
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppEnv:                 getEnv("APP_ENV", "development"),
		Port:                   getEnv("PORT", "8080"),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		WriteRateLimitRequests: writeRateLimitRequests,
		WriteRateLimitWindow:   writeRateLimitWindow,
	}

	if cfg.Port == "" {
		return Config{}, errors.New("PORT must not be empty")
	}

	return cfg, nil
}

func (c Config) HTTPAddr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func (c Config) EffectiveWriteRateLimitRequests() int {
	if c.WriteRateLimitRequests > 0 {
		return c.WriteRateLimitRequests
	}
	if c.AppEnv == "test" {
		return defaultTestWriteRateLimitRequests
	}

	return defaultWriteRateLimitRequests
}

func (c Config) EffectiveWriteRateLimitWindow() time.Duration {
	if c.WriteRateLimitWindow > 0 {
		return c.WriteRateLimitWindow
	}

	return defaultWriteRateLimitWindow
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvInt(key string) (int, error) {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer: %w", key, err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}

	return value, nil
}

func getEnvDuration(key string) (time.Duration, error) {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return 0, nil
	}

	value, err := time.ParseDuration(rawValue)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", key, err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}

	return value, nil
}
