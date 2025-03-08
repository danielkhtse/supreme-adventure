package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/danielkhtse/supreme-adventure/common/models"
	"github.com/danielkhtse/supreme-adventure/common/response"
	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/danielkhtse/supreme-adventure/common/validation"
	"github.com/gorilla/mux"
)

// AccountResponse represents the response for account operations
// swagger:model AccountResponse
type AccountResponse struct {
	// The unique identifier of the account
	// required: true
	// example: 12345
	ID types.AccountID `json:"account_id"`

	// The current balance in smallest currency units (e.g. cents for USD)
	// required: true
	// example: 10000
	Balance types.AccountBalance `json:"balance"`
}

// swagger:route GET /accounts/{account_id} Account getAccount
// Get account details by ID
// responses:
//
//	200: AccountResponse
//	400: description: Invalid account ID format
//	404: description: Account not found
//	500: description: Internal server error
func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	requestAccountId, err := strconv.ParseUint(vars["account_id"], 10, 64)
	if err != nil {
		response.SendError(w, response.StatusBadRequest, "Invalid account ID format")
		return
	}

	account, err := s.AccountService.GetAccount(types.AccountID(requestAccountId))
	if err != nil {
		response.SendError(w, response.StatusNotFound, "Account not found")
		return
	}

	response.SendSuccess[AccountResponse](w, response.StatusOK, &AccountResponse{
		ID:      account.ID,
		Balance: account.Balance, // smallest units for the currency (e.g. cents for USD)
	})
}

// CreateAccountRequest represents the request body for creating an account
// swagger:model CreateAccountRequest
type CreateAccountRequest struct {
	// The unique identifier for the new account
	// required: true
	// example: 12345
	AccountID types.AccountID `json:"account_id" validate:"required,uuid"`

	// The initial balance in smallest currency units (e.g. cents for USD)
	// required: true
	// minimum: 0
	// example: 10000
	InitialBalance types.AccountBalance `json:"initial_balance" validate:"required,min=0"`
}

// swagger:route POST /accounts Account createAccount
// Create a new account
// responses:
//
//	201: description: Account created successfully
//	400: description: Invalid request body or account already exists
//	500: description: Internal server error
func (s *Server) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var request CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.SendError(w, response.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validation.ValidateStruct(request); err != nil {
		response.SendError(w, response.StatusBadRequest, err.Error())
		return
	}

	if err := s.AccountService.CreateAccount(&models.Account{
		ID:             request.AccountID,
		InitialBalance: request.InitialBalance,
	}); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			response.SendError(w, response.StatusBadRequest, "Account already exists")
		} else {
			response.SendError(w, response.StatusInternalServerError, "Failed to create account")
		}
		return
	}

	response.SendSuccess[struct{}](w, response.StatusCreated, nil)
}
