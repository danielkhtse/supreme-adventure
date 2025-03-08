package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/danielkhtse/supreme-adventure/common/models"
	"github.com/danielkhtse/supreme-adventure/common/response"
	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/danielkhtse/supreme-adventure/common/validation"
	"github.com/sirupsen/logrus"
)

type CreateTransactionRequest struct {
	SourceAccountID types.AccountID      `json:"source_account_id" validate:"required"`
	DestAccountID   types.AccountID      `json:"destination_account_id" validate:"required"`
	Amount          types.AccountBalance `json:"amount" validate:"required,min=1"` // smallest units for the currency (e.g. cents for USD)
}

func (s *Server) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("handling create transaction request")

	var request CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logrus.WithError(err).Error("failed to decode request body")
		response.SendError(w, response.StatusBadRequest, "invalid request body")
		return
	}

	logrus.WithFields(logrus.Fields{
		"source_account_id": request.SourceAccountID,
		"dest_account_id":   request.DestAccountID,
		"amount":            request.Amount,
	}).Debug("received create transaction request")

	if err := validation.ValidateStruct(request); err != nil {
		logrus.WithError(err).Error("request validation failed")
		response.SendError(w, response.StatusBadRequest, err.Error())
		return
	}

	transaction := &models.Transaction{
		SourceAccountID: request.SourceAccountID,
		DestAccountID:   request.DestAccountID,
		Amount:          request.Amount,
	}

	if err := s.TransactionService.CreateTransaction(transaction); err != nil {
		errMsg := err.Error()
		logrus.WithError(err).WithFields(logrus.Fields{
			"source_account_id": transaction.SourceAccountID,
			"dest_account_id":   transaction.DestAccountID,
			"amount":            transaction.Amount,
		}).Error("failed to create transaction")

		if strings.Contains(errMsg, "source and destination accounts cannot be the same") {
			response.SendError(w, response.StatusBadRequest, errMsg)
		} else if strings.Contains(errMsg, "source account not found") {
			response.SendError(w, response.StatusNotFound, errMsg)
		} else if strings.Contains(errMsg, "destination account not found") {
			response.SendError(w, response.StatusNotFound, errMsg)
		} else if strings.Contains(errMsg, "insufficient balance") {
			response.SendError(w, response.StatusBadRequest, errMsg)
		} else if strings.Contains(errMsg, "amount must be positive") {
			response.SendError(w, response.StatusBadRequest, errMsg)
		} else {
			response.SendError(w, response.StatusInternalServerError, "failed to create transaction")
		}
		return
	}

	response.SendSuccess(w, response.StatusCreated, transaction)
}
