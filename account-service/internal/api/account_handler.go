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
type AccountResponse struct {
	// The unique identifier of the account
	ID types.AccountID `json:"account_id"`

	// The current balance in smallest currency units (e.g. cents for USD)
	Balance types.AccountBalance `json:"balance"`
}

// @Summary Get account details by ID
// @Description Get account details by ID
// @Tags Account
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Success 200 {object} AccountResponse "Account details"
// @Failure 400 {object} response.ErrorResponse "Invalid account ID format"
// @Failure 404 {object} response.ErrorResponse "Account not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /accounts/{account_id} [get]
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

	response.SendSuccess(w, response.StatusOK, &AccountResponse{
		ID:      account.ID,
		Balance: account.Balance, // smallest units for the currency (e.g. cents for USD)
	})
}

// CreateAccountRequest represents the request body for creating an account
type CreateAccountRequest struct {
	// The unique identifier for the new account
	AccountID types.AccountID `json:"account_id" validate:"required,uuid"`

	// The initial balance in smallest currency units (e.g. cents for USD)
	InitialBalance types.AccountBalance `json:"initial_balance" validate:"required,min=0"`
}

// @Summary Create a new account
// @Description Create a new account with initial balance
// @Tags Account
// @Accept json
// @Produce json
// @Param request body CreateAccountRequest true "Account creation request"
// @Success 201
// @Failure 400 {object} response.ErrorResponse "Invalid request body or account already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /accounts [post]
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
