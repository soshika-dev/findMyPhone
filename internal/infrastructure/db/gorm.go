package db

import (
	"fmt"
	"log"

	"findMyPhone/internal/domain"
	"findMyPhone/internal/infrastructure/config"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGorm opens a gorm DB based on config.
func NewGorm(cfg *config.Config) (*gorm.DB, error) {
	var (
		dialector gorm.Dialector
		err       error
	)

	switch cfg.DatabaseType {
	case "postgres", "postgresql":
		dialector = postgres.Open(cfg.DatabaseURL)
	case "sqlite":
		dialector = sqlite.Open(cfg.DatabaseURL)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DatabaseType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)

	if err := db.AutoMigrate(&domain.User{}, &domain.Device{}, &domain.Log{}); err != nil {
		log.Printf("failed to migrate: %v", err)
		return nil, err
	}

	return db, nil
}
