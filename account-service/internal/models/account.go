package models

import (
	"time"

	"github.com/google/uuid"
)

// Account represents a bank account in the system
type Account struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v6()" validate:"required"`
	AccountNumber string        `json:"accountNumber" gorm:"unique;not null" validate:"required"`
	Balance       int64         `json:"balance" gorm:"default:0" validate:"min=0"`                   //We will store the smallest units for the currency (e.g. cents for USD)
	Currency      string        `json:"currency" gorm:"default:'USD'" validate:"required,oneof=USD"` //We simply support USD for now
	Status        AccountStatus `json:"status" gorm:"type:varchar(10);default:'active';check:status IN ('active', 'inactive')" validate:"required,oneof=active inactive"`
	CreatedAt     time.Time     `json:"createdAt" gorm:"autoCreateTime" validate:"required"`
	UpdatedAt     time.Time     `json:"updatedAt" gorm:"autoUpdateTime" validate:"required"`
}

// AccountStatus represents the status of an account
type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
)
