package api

import (
	"github.com/gorilla/mux"
)

// NewRouter creates and configures a new router
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health-check", HealthCheckHandler).Methods("GET")

	return r
}
