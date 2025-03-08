package models

import (
	"time"

	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/danielkhtse/supreme-adventure/common/validation"
	"gorm.io/gorm"
)

type Transaction struct {
	ID              types.TransactionID     `gorm:"primaryKey" json:"id" validate:"required"`
	SourceAccountID types.AccountID         `gorm:"index" json:"source_account_id" validate:"required"`
	DestAccountID   types.AccountID         `gorm:"index" json:"destination_account_id" validate:"required,nefield=SourceAccountID"`
	Amount          types.AccountBalance    `json:"amount" validate:"required,min=1"`                            //We will store the smallest units for the currency (e.g. cents for USD)
	Currency        string                  `json:"currency" gorm:"default:'USD'" validate:"required,oneof=USD"` //We simply support USD for now
	Status          types.TransactionStatus `gorm:"type:varchar(20)" json:"status" validate:"required,transaction_status"`
	Description     string                  `json:"description" validate:"required"`
	CreatedAt       time.Time               `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time               `json:"updated_at" gorm:"autoUpdateTime"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if err := validation.ValidateStruct(t); err != nil {
		return err
	}

	if t.Status == "" {
		t.Status = types.TransactionStatusPending
	}

	if t.Currency == "" {
		t.Currency = "USD"
	}

	return nil
}

func (t *Transaction) BeforeUpdate(tx *gorm.DB) error {
	if err := validation.ValidateStruct(t); err != nil {
		return err
	}
	return nil
}
