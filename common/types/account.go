package types

type AccountID uint64

type AccountStatus string

type AccountBalance int64

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
)
