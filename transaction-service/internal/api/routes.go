package api

import (
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
	accounts.HandleFunc("", s.CreateTransactionHandler).Methods("POST")

	return r
}
