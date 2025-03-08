package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	accountsRoute = "/accounts"
)

// NewRouter creates and configures a new router
func (s *Server) InitializeRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health-check", s.HealthCheckHandler).Methods("GET")

	accounts := r.PathPrefix(accountsRoute).Subrouter()

	//single account handlers
	accounts.HandleFunc("", s.CreateAccountHandler).Methods("POST")
	accounts.HandleFunc("/{account_id}", s.GetAccountHandler).Methods("GET")
	accounts.HandleFunc("/{account_id}/balance/transfer", s.TransferFundsHandler).Methods("PUT")

	fs := http.FileServer(http.Dir("account-service/docs"))
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))

	return r
}
