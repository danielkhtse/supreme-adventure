package service

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTransferFunds(t *testing.T) {
	mockDB, mock, db := setupMockDB(t)
	defer mockDB.Close()

	service := &AccountService{
		db: db,
	}

	t.Run("Successful transfer", func(t *testing.T) {
		sourceID := types.AccountID(1)
		destID := types.AccountID(2)
		amount := types.AccountBalance(50)

		// Query source account
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(sourceID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).
				AddRow(1, 100))

		// Query destination account
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(destID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).
				AddRow(2, 0))

		// Begin transaction
		mock.ExpectBegin()

		// Update source account
		mock.ExpectExec(`UPDATE "accounts" SET "balance"=\$1 WHERE "id" = \$2`).
			WithArgs(50, sourceID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Update destination account
		mock.ExpectExec(`UPDATE "accounts" SET "balance"=\$1 WHERE "id" = \$2`).
			WithArgs(50, destID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit transaction
		mock.ExpectCommit()

		err := service.TransferFunds(sourceID, destID, amount)
		assert.NoError(t, err)
	})

	t.Run("Insufficient balance", func(t *testing.T) {
		sourceID := types.AccountID(1)
		destID := types.AccountID(2)
		amount := types.AccountBalance(150)

		// Query source account for balance check
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(sourceID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).
				AddRow(1, 100))

		err := service.TransferFunds(sourceID, destID, amount)
		assert.Error(t, err)
		assert.Equal(t, "insufficient balance", err.Error())
	})

	t.Run("Invalid amount", func(t *testing.T) {
		sourceID := types.AccountID(1)
		destID := types.AccountID(2)
		amount := types.AccountBalance(0)

		err := service.TransferFunds(sourceID, destID, amount)
		assert.Error(t, err)
		assert.Equal(t, "amount must be positive", err.Error())
	})

	t.Run("Source account not found", func(t *testing.T) {
		sourceID := types.AccountID(999)
		destID := types.AccountID(2)
		amount := types.AccountBalance(50)

		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1`).
			WithArgs(sourceID).
			WillReturnError(gorm.ErrRecordNotFound)

		err := service.TransferFunds(sourceID, destID, amount)
		assert.Error(t, err)
	})

	t.Run("Destination account not found", func(t *testing.T) {
		sourceID := types.AccountID(1)
		destID := types.AccountID(999)
		amount := types.AccountBalance(50)

		err := service.TransferFunds(sourceID, destID, amount)
		assert.Error(t, err)
	})

	t.Run("Transaction fails expect rollback", func(t *testing.T) {
		sourceID := types.AccountID(1)
		destID := types.AccountID(2)
		amount := types.AccountBalance(50)

		// Query source account for balance check
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(sourceID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).
				AddRow(1, 100))

		// Query source account for update
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(sourceID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).
				AddRow(1, 100))

		// Query destination account
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(destID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).
				AddRow(2, 200))

		// Begin transaction
		mock.ExpectBegin()

		// Update source account - simulate failure
		mock.ExpectExec(`UPDATE "accounts" SET "balance"=\$1 WHERE "id" = \$2`).
			WithArgs(50, sourceID).
			WillReturnError(errors.New("database error"))

		// Expect rollback due to error
		mock.ExpectRollback()

		err := service.TransferFunds(sourceID, destID, amount)
		assert.Error(t, err)
	})

}
