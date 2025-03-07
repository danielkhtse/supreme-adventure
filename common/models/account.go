package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/danielkhtse/supreme-adventure/common/validation"
)

// Account represents a bank account in the system
type Account struct {
	ID             types.AccountID      `json:"id" gorm:"primaryKey" validate:"required"`
	Balance        types.AccountBalance `json:"balance" gorm:"default:0" validate:"required,min=0"`          //We will store the smallest units for the currency (e.g. cents for USD)
	InitialBalance types.AccountBalance `json:"initial_balance" gorm:"default:0" validate:"required,min=0"`  //audit trail for the initial balance. TODO: discussion, reflect from transactions for audit trail
	Currency       string               `json:"currency" gorm:"default:'USD'" validate:"required,oneof=USD"` //We simply support USD for now
	Status         types.AccountStatus  `json:"status" gorm:"type:varchar(10);default:'active';check:status IN ('active', 'inactive')" validate:"required,oneof=active inactive"`
	CreatedAt      time.Time            `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time            `json:"updatedAt" gorm:"autoUpdateTime"`
}

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
		a.Status = types.AccountStatusActive
	}

	//if initial balance is set, set the balance to the initial balance
	//TODO: discussion, drop this field and use transactions for audit trail
	if a.InitialBalance > 0 {
		a.Balance = a.InitialBalance
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
