package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	transactionsRoute = "/transactions"
)

// NewRouter creates and configures a new router
func (s *Server) InitializeRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health-check", s.HealthCheckHandler).Methods("GET")

	transactions := r.PathPrefix(transactionsRoute).Subrouter()

	//single account handlers
	transactions.HandleFunc("", s.CreateTransactionHandler).Methods("POST")

	fs := http.FileServer(http.Dir("transaction-service/docs"))
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))

	return r
}
