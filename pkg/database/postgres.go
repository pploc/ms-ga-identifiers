package database

import (
	"fmt"
	"log"

	"github.com/gym-api/ms-ga-identifier/internal/infrastructure/persistence/gorm/model"
	"github.com/gym-api/ms-ga-identifier/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(
		&model.IdentityModel{},
		&model.RefreshTokenModel{},
		&model.LoginAttemptModel{},
		&model.PasswordResetModel{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}
