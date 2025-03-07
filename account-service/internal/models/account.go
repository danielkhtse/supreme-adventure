package models

import (
	"time"

	"github.com/danielkhtse/supreme-adventure/common/validation"
	"gorm.io/gorm"
)

// Account represents a bank account in the system
type Account struct {
	ID        uint64        `json:"id" gorm:"primaryKey" validate:"required"`
	Balance   int64         `json:"balance" gorm:"default:0" validate:"required,min=0"`          //We will store the smallest units for the currency (e.g. cents for USD)
	Currency  string        `json:"currency" gorm:"default:'USD'" validate:"required,oneof=USD"` //We simply support USD for now
	Status    AccountStatus `json:"status" gorm:"type:varchar(10);default:'active';check:status IN ('active', 'inactive')" validate:"required,oneof=active inactive"`
	CreatedAt time.Time     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time     `json:"updatedAt" gorm:"autoUpdateTime"`
}

// AccountStatus represents the status of an account
type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
)

const (
	AccountTableName = "accounts"
)

func (a *Account) TableName() string {
	return AccountTableName
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	//default assignment
	if a.Currency == "" {
		a.Currency = "USD"
	}
	if a.Status == "" {
		a.Status = AccountStatusActive
	}

	if err = validation.ValidateStruct(a); err != nil {
		return err
	}

	return nil
}

func (a *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	if err = validation.ValidateStruct(a); err != nil {
		return err
	}
	return nil
}
