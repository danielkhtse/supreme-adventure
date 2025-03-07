package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/danielkhtse/supreme-adventure/common/models"
	"github.com/danielkhtse/supreme-adventure/common/response"
	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/danielkhtse/supreme-adventure/common/validation"
	"github.com/gorilla/mux"
)

type AccountResponse struct {
	ID      types.AccountID `json:"account_id"`
	Balance string          `json:"balance"`
}

// GetAccountHandler handles getting a single account
func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	requestAccountId, err := strconv.ParseUint(vars["account_id"], 10, 64)
	if err != nil {
		response.SendError(w, response.StatusBadRequest, "Invalid account ID format")
		return
	}

	account, err := s.AccountService.GetAccount(requestAccountId)
	if err != nil {
		response.SendError(w, response.StatusNotFound, "Account not found")
		return
	}

	response.SendSuccess[AccountResponse](w, response.StatusOK, "", &AccountResponse{
		ID:      account.ID,
		Balance: fmt.Sprintf("%.2f", float64(account.Balance)/100), //TODO: simplify the balance to 2 decimal places for now
	})
}

// CreateAccountRequest represents the request body for creating an account
type CreateAccountRequest struct {
	AccountID      types.AccountID `json:"account_id" validate:"required,uuid"`
	InitialBalance string          `json:"initial_balance" validate:"required"` //allow upto 2dp
}

// CreateAccountHandler handles creating a new account
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

	// Convert balance from decimal string to int64 (cents)
	balanceFloat, err := strconv.ParseFloat(request.InitialBalance, 64)
	if err != nil {
		response.SendError(w, response.StatusBadRequest, "Invalid balance format")
		return
	}
	balanceInt := int64(balanceFloat * 100) //depends on the currency, we simply it and assumpt it's USD for now and convert it to cents

	if err := s.AccountService.CreateAccount(&models.Account{
		ID:      request.AccountID,
		Balance: balanceInt,
	}); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			response.SendError(w, response.StatusBadRequest, "Account already exists")
		} else {
			response.SendError(w, response.StatusInternalServerError, "Failed to create account")
		}
		return
	}

	response.SendSuccess[struct{}](w, response.StatusCreated, "", nil)
}
