package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danielkhtse/supreme-adventure/common/db"
	"github.com/danielkhtse/supreme-adventure/common/models"
	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/danielkhtse/supreme-adventure/transaction-service/internal/client"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// TransactionService handles business logic for transaction operations
type TransactionService struct {
	db            *gorm.DB
	accountClient *client.AccountClient
}

// NewTransactionService creates a new TransactionService instance
func NewTransactionService() *TransactionService {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error loading .env file")
	}

	dbDSN := os.Getenv("TRANSACTION_DATABASE_DSN")
	if dbDSN == "" {
		log.Fatal("TRANSACTION_DATABASE_DSN environment variable not set")
	}

	db, err := db.NewPostgresDB(dbDSN)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: use migration script to replace AutoMigrate
	if err := db.GetDB().AutoMigrate(&models.Transaction{}); err != nil {
		log.Fatal(err)
	}

	accountClient := client.NewAccountClient(os.Getenv("ACCOUNT_SERVICE_URL"))

	return &TransactionService{
		db:            db.GetDB(),
		accountClient: accountClient,
	}
}

func (s *TransactionService) CreateTransaction(transaction *models.Transaction) error {

	if transaction == nil {
		return fmt.Errorf("transaction cannot be nil")
	}

	if transaction.SourceAccountID == transaction.DestAccountID {
		return fmt.Errorf("source and destination accounts cannot be the same")
	}

	sourceAccount, err := s.accountClient.GetAccount(transaction.SourceAccountID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("source account not found")
		}
		return fmt.Errorf("failed to fetch source account: %w", err)
	}

	if sourceAccount.Balance < types.AccountBalance(transaction.Amount) {
		return fmt.Errorf("insufficient balance in source account %d", transaction.SourceAccountID)
	}

	_, err = s.accountClient.GetAccount(transaction.DestAccountID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("destination account not found")
		}
		return fmt.Errorf("failed to fetch destination account: %w", err)
	}

	//create trasnaction as pending
	transaction.Status = types.TransactionStatusPending

	//save transaction to db
	if err := s.db.Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	//perform transfer
	if err := s.TransferFunds(transaction); err != nil {
		return err
	}

	return nil
}
