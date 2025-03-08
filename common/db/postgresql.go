package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Optimized connection pool settings
	sqlDB.SetMaxIdleConns(10)                  // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)                 // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Maximum lifetime of a connection
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Maximum idle time for connections

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) GetDB() *gorm.DB {
	return p.db
}
