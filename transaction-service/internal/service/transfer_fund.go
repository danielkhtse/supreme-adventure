package service

import (
	"fmt"
	"log"

	"github.com/danielkhtse/supreme-adventure/common/models"
	"github.com/danielkhtse/supreme-adventure/common/types"
)

// TransferFunds transfers funds between two accounts and creates a transaction record
func (s *TransactionService) TransferFunds(transaction *models.Transaction) error {

	if transaction == nil {
		return fmt.Errorf("transaction cannot be nil")
	}

	// Call account service to transfer funds
	err := s.accountClient.TransferFunds(transaction.SourceAccountID, transaction.DestAccountID, transaction.Amount)
	if err != nil {
		log.Printf("Failed to transfer funds: %v", err)
		transaction.Status = types.TransactionStatusFailed
		transaction.Description = err.Error()
		if dbErr := s.db.Save(transaction).Error; dbErr != nil {
			return fmt.Errorf("failed to update transaction status after transfer failure: %w", dbErr)
		}
		return fmt.Errorf("failed to transfer funds: %w", err)
	}

	// Update transaction status to completed
	transaction.Status = types.TransactionStatusCompleted
	if err := s.db.Save(transaction).Error; err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	return nil
}
