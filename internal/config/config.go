package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"blog-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DB         *gorm.DB
	JWTSecret  string
	ServerPort string
}

var App *Config

func Init() (*Config, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		return nil, errors.New("DATABASE_DSN is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	App = &Config{
		DB:         db,
		JWTSecret:  jwtSecret,
		ServerPort: getServerPort(),
	}

	slog.Info("Configuration loaded successfully",
		slog.String("port", App.ServerPort),
		slog.String("jwt_secret_length", "***"))

	return App, nil
}

func getServerPort() string {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return ":" + port
	}
	return ":8080"
}

func (c *Config) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *Config) DBWithContext(ctx context.Context) *gorm.DB {
	return c.DB.WithContext(ctx)
}

func RunMigrations(db *gorm.DB) error {
	slog.Info("Running database migrations...")

	// AutoMigrate - for development only
	// In production, use golang-migrate with sql-migrations
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	slog.Info("Migrations completed successfully")
	return nil
}
