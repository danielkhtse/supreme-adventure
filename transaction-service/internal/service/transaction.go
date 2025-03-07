package service

import (
	"fmt"
	"log"
	"os"

	"github.com/danielkhtse/supreme-adventure/common/db"
	"github.com/danielkhtse/supreme-adventure/common/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// TransactionService handles business logic for transaction operations
type TransactionService struct {
	db *gorm.DB
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

	return &TransactionService{
		db: db.GetDB(),
	}
}
