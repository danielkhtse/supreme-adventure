package models

import (
	"github.com/danielkhtse/supreme-adventure/common/types"
)

type Transaction struct {
	ID              types.TransactionID     `gorm:"primaryKey" json:"id" validate:"required"`
	SourceAccountID types.AccountID         `gorm:"index" json:"source_account_id" validate:"required"`
	DestAccountID   types.AccountID         `gorm:"index" json:"destination_account_id" validate:"required"`
	Amount          uint64                  `json:"amount" validate:"required,min=0.01"`                         //We will store the smallest units for the currency (e.g. cents for USD)
	Currency        string                  `json:"currency" gorm:"default:'USD'" validate:"required,oneof=USD"` //We simply support USD for now
	Status          types.TransactionStatus `gorm:"type:varchar(20)" json:"status" validate:"required,transaction_status"`
	Type            types.TransactionType   `gorm:"type:varchar(20)" json:"type" validate:"required,oneof=credit debit"` //TODO: Add validation for type
}
