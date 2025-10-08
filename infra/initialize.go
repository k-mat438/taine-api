package infra

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// package infra

type DB struct {
	*gorm.DB
}

func NewDB(dsn string) (*DB, error) {
	logger := logger.Default.LogMode(logger.Warn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
