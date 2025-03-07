package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockDB struct {
	db   *gorm.DB
	Mock sqlmock.Sqlmock
}

// NewMockDB creates a new mock database for testing
func NewMockDB() (*MockDB, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}))
	if err != nil {
		return nil, err
	}

	return &MockDB{
		db:   gormDB,
		Mock: mock,
	}, nil
}

func (m *MockDB) GetDB() *gorm.DB {
	return m.db
}
