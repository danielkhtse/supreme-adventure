package service

import (
	"errors"

	"github.com/danielkhtse/supreme-adventure/common/db"
	"github.com/danielkhtse/supreme-adventure/shared/models"
)

// AccountService handles business logic for account operations
type AccountService struct {
	db *db.PostgresDB
}

// NewAccountService creates a new AccountService instance
func NewAccountService(db *db.PostgresDB) *AccountService {
	return &AccountService{
		db: db,
	}
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(account *models.Account) error {
	if account == nil {
		return errors.New("account cannot be nil")
	}
	return s.db.GetDB().Create(account).Error
}

// GetAccount retrieves an account by ID
func (s *AccountService) GetAccount(id uint) (*models.Account, error) {
	var account models.Account
	if err := s.db.GetDB().First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// UpdateAccount updates an existing account
func (s *AccountService) UpdateAccount(account *models.Account) error {
	if account == nil {
		return errors.New("account cannot be nil")
	}
	return s.db.GetDB().Save(account).Error
}

// DeleteAccount deletes an account by ID
func (s *AccountService) DeleteAccount(id uint) error {
	return s.db.GetDB().Delete(&models.Account{}, id).Error
}
