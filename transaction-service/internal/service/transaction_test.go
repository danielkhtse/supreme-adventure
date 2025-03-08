package service

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danielkhtse/supreme-adventure/common/models"
	"github.com/danielkhtse/supreme-adventure/transaction-service/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	// Create a new SQL mock
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 mockDB,
		PreferSimpleProtocol: true,
	})

	// Open a gorm DB connection with the mock
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return mockDB, mock, db
}

func TestUnitCreateTransaction(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response interface{}
		var exists bool

		switch r.URL.Path {
		case "/accounts/1":
			response = &models.Account{
				ID:      1,
				Balance: 200,
			}
			exists = true
		case "/accounts/2":
			response = &models.Account{
				ID:      2,
				Balance: 50,
			}
			exists = true
		case "/accounts/1/transfer":
			w.WriteHeader(http.StatusOK)
			return
		default:
			exists = false
		}

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	mockDB, mock, db := setupMockDB(t)
	defer mockDB.Close()

	mockService := &TransactionService{
		db:            db,
		accountClient: client.NewAccountClient(mockServer.URL),
	}

	t.Run("Nil transaction", func(t *testing.T) {
		err := mockService.CreateTransaction(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction cannot be nil")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Same source and destination", func(t *testing.T) {
		transaction := &models.Transaction{
			SourceAccountID: 1,
			DestAccountID:   1,
			Amount:          100,
		}

		err := mockService.CreateTransaction(transaction)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "source and destination accounts cannot be the same")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Insufficient balance", func(t *testing.T) {
		transaction := &models.Transaction{
			SourceAccountID: 1,
			DestAccountID:   2,
			Amount:          300,
		}

		err := mockService.CreateTransaction(transaction)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient balance in source account 1")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Source account not found", func(t *testing.T) {
		transaction := &models.Transaction{
			SourceAccountID: 999, // Non-existent account ID
			DestAccountID:   2,
			Amount:          100,
		}

		err := mockService.CreateTransaction(transaction)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "source account not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Destination account not found", func(t *testing.T) {
		transaction := &models.Transaction{
			SourceAccountID: 1,
			DestAccountID:   999, // Non-existent account ID
			Amount:          100,
		}

		err := mockService.CreateTransaction(transaction)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "destination account not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success case", func(t *testing.T) {
		transaction := &models.Transaction{
			SourceAccountID: 1,   // Using account ID 1 with balance 200
			DestAccountID:   2,   // Using account ID 2 with balance 50
			Amount:          100, // Amount that source account can afford (200 > 100)
			Currency:        "USD",
		}

		// Mock the BEGIN transaction
		mock.ExpectBegin()
		// Mock the INSERT query for creating pending transaction
		mock.ExpectQuery("INSERT INTO \"transactions\"").
			WithArgs(transaction.SourceAccountID, transaction.DestAccountID, transaction.Amount, transaction.Currency, "pending", "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		// Mock the COMMIT transaction since we expect success
		mock.ExpectCommit()

		err := mockService.CreateTransaction(transaction)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
