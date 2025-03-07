package types

type TransactionID uint64
type TransactionStatus string
type TransactionType string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

const (
	TypeCredit TransactionType = "credit"
	TypeDebit  TransactionType = "debit"
)
