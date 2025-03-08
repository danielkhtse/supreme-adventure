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

// TransferFundsRequest represents the request body for transferring funds between accounts
// swagger:model TransferFundsRequest
type TransferFundsRequest struct {
	// The destination account ID to transfer funds to
	// required: true
	// example: 12345
	DestAccountID types.AccountID `json:"dest_account_id" validate:"required,uuid"`

	// The amount to transfer in smallest currency units (e.g. cents for USD)
	// required: true
	// minimum: 1
	// example: 1000
	Amount types.AccountBalance `json:"amount" validate:"required,min=1"`
}

// swagger:route POST /accounts/{account_id}/transfer Account transferFunds
// Transfer funds between accounts
// responses:
//
//	200: description: Funds transferred successfully
//	400: description: Invalid request parameters or insufficient balance
//	404: description: Source or destination account not found
//	500: description: Internal server error
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
