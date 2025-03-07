package service

import (
	"testing"

	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/danielkhtse/supreme-adventure/account-service/internal/models"
)

// setupMockDB creates a new mock database connection and gorm DB instance
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

func TestCreateAccount(t *testing.T) {
	mockDB, mock, db := setupMockDB(t)
	defer mockDB.Close()

	service := &AccountService{
		db: db,
	}

	t.Run("Success", func(t *testing.T) {
		account := &models.Account{
			ID:       1,
			Balance:  100.0,
			Currency: "USD",
			Status:   "active",
		}

		// Expect check for existing account
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(account.ID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Expect account creation
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "accounts" \("balance","currency","status","created_at","updated_at","id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
			WithArgs(account.Balance, account.Currency, account.Status, sqlmock.AnyArg(), sqlmock.AnyArg(), account.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := service.CreateAccount(account)
		assert.NoError(t, err)
	})

	t.Run("Account already exists", func(t *testing.T) {
		account := &models.Account{
			ID:      1,
			Balance: 100.0,
		}

		// Expect check for existing account to find one
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(account.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow(1, 100.0))

		err := service.CreateAccount(account)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Nil account", func(t *testing.T) {
		err := service.CreateAccount(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})
}

func TestGetAccount(t *testing.T) {
	mockDB, mock, db := setupMockDB(t)
	defer mockDB.Close()

	service := &AccountService{
		db: db,
	}

	t.Run("Success", func(t *testing.T) {
		expectedAccount := &models.Account{
			ID:      1,
			Balance: 100.0,
		}

		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(expectedAccount.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).
				AddRow(expectedAccount.ID, expectedAccount.Balance))

		account, err := service.GetAccount(expectedAccount.ID)
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, expectedAccount.ID, account.ID)
		assert.Equal(t, expectedAccount.Balance, account.Balance)
	})

	t.Run("Account not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"."id" LIMIT \$2`).
			WithArgs(uint64(999), 1).
			WillReturnError(gorm.ErrRecordNotFound)

		account, err := service.GetAccount(999)
		assert.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
