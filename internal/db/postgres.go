package db

import (
	"context"
	"fmt"
	"time"

	"github.com/BabichevDima/subManager/internal/config"
	"github.com/BabichevDima/subManager/internal/models"
	"github.com/BabichevDima/subManager/pkg/logger"

	"go.uber.org/zap"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func InitPostgres(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	logger.Info(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Fatal("db ping failed", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&models.Subscription{}); err != nil {
		logger.Fatal("db migration failed", zap.Error(err))
	}

	return db, nil
}

func ShutdownDB(db *gorm.DB) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
}
