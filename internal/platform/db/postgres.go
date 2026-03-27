package db

import (
	"errors"
	"time"

	"go-gin-ecommerce/internal/platform/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL must be set")
	}

	return gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Warn),
		NowFunc:     time.Now().UTC,
	})
}
