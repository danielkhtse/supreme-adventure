package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danielkhtse/supreme-adventure/common/models"
	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/danielkhtse/supreme-adventure/transaction-service/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestUnitTransferFunds(t *testing.T) {
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
		case "/accounts/1/balance/transfer":
			w.WriteHeader(http.StatusOK)
			return
		case "/accounts/999/balance/transfer":
			w.WriteHeader(http.StatusNotFound)

			errorResp := &struct {
				Message string `json:"message"`
			}{
				Message: "source_account_not_found",
			}
			json.NewEncoder(w).Encode(errorResp)
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
		err := mockService.TransferFunds(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction cannot be nil")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Successful transfer", func(t *testing.T) {
		transaction := &models.Transaction{
			SourceAccountID: 1,
			DestAccountID:   2,
			Amount:          100,
			Currency:        "USD",
			Status:          types.TransactionStatusPending,
			Description:     "",
		}
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"transactions\" \\(\"source_account_id\",\"dest_account_id\",\"amount\",\"currency\",\"status\",\"description\",\"created_at\",\"updated_at\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6,\\$7,\\$8\\) RETURNING \"id\"").
			WithArgs(transaction.SourceAccountID, transaction.DestAccountID, transaction.Amount, transaction.Currency, types.TransactionStatusCompleted, transaction.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := mockService.TransferFunds(transaction)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
