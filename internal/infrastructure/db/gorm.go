package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

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
		if err := ensurePostgresDatabase(cfg.DatabaseURL); err != nil {
			return nil, fmt.Errorf("ensure postgres database: %w", err)
		}
		dialector = postgres.Open(cfg.DatabaseURL)
	case "sqlite":
		dialector = sqlite.Open(cfg.DatabaseURL)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DatabaseType)
	}

	var db *gorm.DB
	const (
		maxAttempts = 5
		baseDelay   = 2 * time.Second
	)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(dialector, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		log.Printf("failed to connect to database (attempt %d/%d): %v", attempt, maxAttempts, err)
		time.Sleep(time.Duration(attempt) * baseDelay)
	}
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
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

// ensurePostgresDatabase guarantees that the database referenced in the
// provided connection string exists by creating it when necessary.
func ensurePostgresDatabase(databaseURL string) error {
	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return fmt.Errorf("parse postgres url: %w", err)
	}

	dbName := strings.TrimPrefix(parsedURL.Path, "/")
	if dbName == "" {
		return fmt.Errorf("postgres database name not found in url")
	}

	adminURL := *parsedURL
	adminURL.Path = "/postgres"

	adminDB, err := sql.Open("pgx", adminURL.String())
	if err != nil {
		return fmt.Errorf("open postgres admin connection: %w", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		return fmt.Errorf("ping postgres admin: %w", err)
	}

	var exists bool
	if err := adminDB.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists); err != nil {
		return fmt.Errorf("check database existence: %w", err)
	}
	if exists {
		return nil
	}

	quotedName := strings.ReplaceAll(dbName, "\"", "\"\"")
	if _, err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", quotedName)); err != nil {
		return fmt.Errorf("create database %s: %w", dbName, err)
	}

	return nil
}
