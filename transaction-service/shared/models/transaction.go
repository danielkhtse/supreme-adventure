package models

import (
	account "github.com/danielkhtse/supreme-adventure/account-service/shared/models"
)

type TransactionID uint64

type TransactionStatus string
type TransactionType string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"

	TypeCredit TransactionType = "credit"
	TypeDebit  TransactionType = "debit"
)

type Transaction struct {
	ID              TransactionID     `gorm:"primaryKey" json:"id" validate:"required"`
	SourceAccountID account.AccountID `gorm:"index" json:"source_account_id" validate:"required"`
	DestAccountID   account.AccountID `gorm:"index" json:"destination_account_id" validate:"required"`
	Amount          uint64            `json:"amount" validate:"required,min=0.01"`                         //We will store the smallest units for the currency (e.g. cents for USD)
	Currency        string            `json:"currency" gorm:"default:'USD'" validate:"required,oneof=USD"` //We simply support USD for now
	Status          TransactionStatus `gorm:"type:varchar(20)" json:"status" validate:"required,transaction_status"`
	Type            TransactionType   `gorm:"type:varchar(20)" json:"type" validate:"required,oneof=credit debit"` //TODO: Add validation for type
}
