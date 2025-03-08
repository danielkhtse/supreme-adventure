package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/danielkhtse/supreme-adventure/common/response"
	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/danielkhtse/supreme-adventure/common/validation"
	"github.com/gorilla/mux"
)

// @Summary Transfer funds between accounts
// @Description Transfer funds from source account to destination account
// @Tags Account
// @Accept json
// @Produce json
// @Param account_id path string true "Source Account ID"
// @Param request body TransferFundsRequest true "Transfer request details"
// @Success 200
// @Failure 400 {object} response.ErrorResponse "Invalid request parameters or insufficient balance"
// @Failure 404 {object} response.ErrorResponse "Source or destination account not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /accounts/{account_id}/transfer [post]
type TransferFundsRequest struct {
	// The destination account ID to transfer funds to
	DestAccountID types.AccountID `json:"dest_account_id" validate:"required,uuid"` // @example 12345

	// The amount to transfer in smallest currency units (e.g. cents for USD)
	Amount types.AccountBalance `json:"amount" validate:"required,min=1"` // @example 1000
}

func (s *Server) TransferFundsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sourceIDStr := vars["account_id"]
	if sourceIDStr == "" {
		response.SendError(w, response.StatusBadRequest, "source account ID is required")
		return
	}

	sourceAccountID, err := strconv.ParseUint(sourceIDStr, 10, 64)
	if err != nil {
		response.SendError(w, response.StatusBadRequest, "invalid source account ID format")
		return
	}

	var req TransferFundsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, response.StatusBadRequest, "invalid request body")
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		response.SendError(w, response.StatusBadRequest, err.Error())
		return
	}

	err = s.AccountService.TransferFunds(types.AccountID(sourceAccountID), req.DestAccountID, req.Amount)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "source account not found") {
			response.SendError(w, response.StatusNotFound, errStr)
		} else if strings.Contains(errStr, "destination account not found") {
			response.SendError(w, response.StatusNotFound, errStr)
		} else if strings.Contains(errStr, "insufficient balance") {
			response.SendError(w, response.StatusBadRequest, errStr)
		} else if strings.Contains(errStr, "amount must be positive") {
			response.SendError(w, response.StatusBadRequest, errStr)
		} else {
			response.SendError(w, response.StatusInternalServerError, "failed to transfer funds")
		}
		return
	}

	response.SendSuccess[struct{}](w, response.StatusOK, nil)
}
