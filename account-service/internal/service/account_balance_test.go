package service

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUnitTransferFunds(t *testing.T) {
	mockDB, mock, db := setupMockDB(t)
	defer mockDB.Close()

	service := &AccountService{
		db: db,
	}

	t.Run("Successful transfer", func(t *testing.T) {
		t.Log("Testing successful transfer between accounts")
		sourceID := types.AccountID(1)
		destID := types.AccountID(2)
		amount := types.AccountBalance(50)

		// Begin transaction
		mock.ExpectBegin()

		// Query source account with FOR UPDATE NOWAIT
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE "accounts"."id" = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "initial_balance", "currency", "status", "created_at", "updated_at"}).
				AddRow(sourceID, float64(100), float64(100), "USD", "active", time.Time{}, time.Time{}))

		t.Log("Source account found with balance: 100")

		// Query destination account with FOR UPDATE NOWAIT
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE "accounts"."id" = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(2, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "initial_balance", "currency", "status", "created_at", "updated_at"}).
				AddRow(destID, float64(0), float64(0), "USD", "active", time.Time{}, time.Time{}))

		t.Log("Destination account found with balance: 0")

		// Update source account
		mock.ExpectExec(`UPDATE "accounts" SET "balance"=\$1,"initial_balance"=\$2,"currency"=\$3,"status"=\$4,"created_at"=\$5,"updated_at"=\$6 WHERE "id" = \$7`).
			WithArgs(50, 100, "USD", "active", time.Time{}, sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		t.Log("Updated source account balance to: 50")

		// Update destination account
		mock.ExpectExec(`UPDATE "accounts" SET "balance"=\$1,"initial_balance"=\$2,"currency"=\$3,"status"=\$4,"created_at"=\$5,"updated_at"=\$6 WHERE "id" = \$7`).
			WithArgs(50, 0, "USD", "active", time.Time{}, sqlmock.AnyArg(), 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		t.Log("Updated destination account balance to: 50")

		// Commit transaction
		mock.ExpectCommit()

		err := service.TransferFunds(sourceID, destID, amount)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		t.Log("Transfer completed successfully")
	})

	t.Run("Insufficient balance", func(t *testing.T) {
		t.Log("Testing transfer with insufficient balance")
		sourceID := types.AccountID(1)
		destID := types.AccountID(2)
		amount := types.AccountBalance(150)

		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		// Begin transaction
		mock.ExpectBegin()

		// Query source account with FOR UPDATE NOWAIT
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE "accounts"."id" = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "initial_balance", "currency", "status", "created_at", "updated_at"}).
				AddRow(1, 100, 0, "", "", createdAt, updatedAt))

		t.Log("Source account found with balance: 100")

		// Query destination account with FOR UPDATE NOWAIT
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE "accounts"."id" = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(2, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "initial_balance", "currency", "status", "created_at", "updated_at"}).
				AddRow(2, 0, 0, "", "", createdAt, updatedAt))

		t.Log("Destination account found with balance: 0")

		// Expect rollback since balance is insufficient
		mock.ExpectRollback()

		err := service.TransferFunds(sourceID, destID, amount)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient balance")
		assert.NoError(t, mock.ExpectationsWereMet())
		t.Log("Transfer failed as expected due to insufficient balance")
	})

	t.Run("Invalid amount", func(t *testing.T) {
		t.Log("Testing transfer with invalid (zero) amount")
		sourceID := types.AccountID(1)
		destID := types.AccountID(2)
		amount := types.AccountBalance(0)

		err := service.TransferFunds(sourceID, destID, amount)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "amount must be positive")
		assert.NoError(t, mock.ExpectationsWereMet())
		t.Log("Transfer failed as expected due to invalid amount")
	})

	t.Run("Source account not found", func(t *testing.T) {
		t.Log("Testing transfer with non-existent source account")
		sourceID := types.AccountID(999)
		destID := types.AccountID(2)
		amount := types.AccountBalance(50)

		// Begin transaction since TransferFunds starts a transaction
		mock.ExpectBegin()

		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE "accounts"."id" = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(2, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Expect rollback since source account not found
		mock.ExpectRollback()

		err := service.TransferFunds(sourceID, destID, amount)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "account not found")
		assert.NoError(t, mock.ExpectationsWereMet())
		t.Log("Transfer failed as expected due to source account not found")
	})

	t.Run("Destination account not found", func(t *testing.T) {
		t.Log("Testing transfer with non-existent destination account")
		sourceID := types.AccountID(1)
		destID := types.AccountID(999)
		amount := types.AccountBalance(50)

		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		// Begin transaction since TransferFunds starts a transaction
		mock.ExpectBegin()

		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE "accounts"."id" = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "initial_balance", "currency", "status", "created_at", "updated_at"}).
				AddRow(1, 100, 0, "", "", createdAt, updatedAt))

		t.Log("Source account found with balance: 100")

		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE "accounts"."id" = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Expect rollback since destination account not found
		mock.ExpectRollback()

		err := service.TransferFunds(sourceID, destID, amount)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "account not found")
		assert.NoError(t, mock.ExpectationsWereMet())
		t.Log("Transfer failed as expected due to destination account not found")
	})

}
