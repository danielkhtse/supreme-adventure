package service

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/danielkhtse/supreme-adventure/account-service/internal/models"
	"github.com/danielkhtse/supreme-adventure/common/db"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// AccountService handles business logic for account operations
type AccountService struct {
	db *gorm.DB
}

// NewAccountService creates a new AccountService instance
func NewAccountService() *AccountService {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error loading .env file")
	}

	dbDSN := os.Getenv("ACCOUNT_DATABASE_DSN")
	if dbDSN == "" {
		log.Fatal("ACCOUNT_DATABASE_DSN environment variable not set")
	}

	db, err := db.NewPostgresDB(dbDSN)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: use migration script to replace AutoMigrate
	if err := db.GetDB().AutoMigrate(&models.Account{}); err != nil {
		log.Fatal(err)
	}

	return &AccountService{
		db: db.GetDB(),
	}
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(account *models.Account) error {
	if account == nil {
		return errors.New("account cannot be nil")
	}

	// Check if account already exists
	var existingAccount models.Account
	if err := s.db.Model(&models.Account{}).First(&existingAccount, "id = ?", account.ID).Error; err == nil {
		return fmt.Errorf("account with ID %d already exists", account.ID)
	}

	return s.db.Model(&models.Account{}).Create(account).Error
}

// GetAccount retrieves an account by ID
func (s *AccountService) GetAccount(id uint64) (*models.Account, error) {
	var account models.Account
	if err := s.db.First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
