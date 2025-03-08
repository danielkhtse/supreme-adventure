package service

import (
	"errors"

	"github.com/danielkhtse/supreme-adventure/common/models"
	"github.com/danielkhtse/supreme-adventure/common/types"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

// TransferFunds transfers funds between two accounts
func (s *AccountService) TransferFunds(sourceAccountID types.AccountID, destAccountID types.AccountID, amount types.AccountBalance) error {
	log.WithFields(log.Fields{
		"source_account_id": sourceAccountID,
		"dest_account_id":   destAccountID,
		"amount":            amount,
	}).Info("starting funds transfer")

	if amount <= 0 {
		log.WithError(errors.New("amount must be positive")).Error("invalid transfer amount")
		return errors.New("amount must be positive")
	}

	// Prevent self-transfers
	if sourceAccountID == destAccountID {
		log.WithError(errors.New("cannot transfer to same account")).Error("invalid transfer")
		return errors.New("cannot transfer to same account")
	}

	// Lock accounts in consistent order to prevent deadlocks
	firstAccountID := sourceAccountID
	secondAccountID := destAccountID
	if destAccountID < sourceAccountID {
		firstAccountID = destAccountID
		secondAccountID = sourceAccountID
	}

	log.Debug("attempting to begin transaction")
	tx := s.db.Begin()
	if err := tx.Error; err != nil {
		log.WithError(err).Error("failed to begin transaction")
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	log.WithField("account_id", firstAccountID).Debug("acquiring lock on first account")
	var firstAccount models.Account
	if err := tx.Set("gorm:query_option", "FOR UPDATE WAIT 5").First(&firstAccount, firstAccountID).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("failed to acquire first account")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("account not found")
		}
		// Check for lock timeout error
		if err.Error() == "lock timeout" {
			return errors.New("failed to acquire lock - timeout after 5 seconds")
		}
		return err
	}

	log.WithField("account_id", secondAccountID).Debug("acquiring lock on second account")
	var secondAccount models.Account
	if err := tx.Set("gorm:query_option", "FOR UPDATE WAIT 5").First(&secondAccount, secondAccountID).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("failed to acquire second account")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("account not found")
		}
		// Check for lock timeout error
		if err.Error() == "lock timeout" {
			return errors.New("failed to acquire lock - timeout after 5 seconds")
		}
		return err
	}

	// Map back to source and dest accounts
	sourceAccount := firstAccount
	destAccount := secondAccount
	if destAccountID < sourceAccountID {
		sourceAccount = secondAccount
		destAccount = firstAccount
	}

	// Check balance after getting locked records
	if sourceAccount.Balance < amount {
		tx.Rollback()
		log.WithFields(log.Fields{
			"available_balance": sourceAccount.Balance,
			"required_amount":   amount,
		}).Error("insufficient balance")
		return errors.New("insufficient balance")
	}

	log.WithFields(log.Fields{
		"account_id":  sourceAccount.ID,
		"old_balance": sourceAccount.Balance,
		"new_balance": sourceAccount.Balance - amount,
	}).Debug("updating source account balance")
	sourceAccount.Balance -= amount
	if err := tx.Save(&sourceAccount).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("failed to update source account")
		return err
	}

	log.WithFields(log.Fields{
		"account_id":  destAccount.ID,
		"old_balance": destAccount.Balance,
		"new_balance": destAccount.Balance + amount,
	}).Debug("updating destination account balance")
	destAccount.Balance += amount
	if err := tx.Save(&destAccount).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("failed to update destination account")
		return err
	}

	log.Debug("committing transaction")
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("failed to commit transaction")
		return err
	}

	log.WithFields(log.Fields{
		"source_balance": sourceAccount.Balance,
		"dest_balance":   destAccount.Balance,
	}).Info("successfully completed funds transfer")
	return nil
}
