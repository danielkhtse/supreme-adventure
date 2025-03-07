package api

import (
	"net/http"
	"os"

	"github.com/danielkhtse/supreme-adventure/transaction-service/internal/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// Transactions API server
type Server struct {
	TransactionService *service.TransactionService
	Router             *mux.Router
	Port               string
}

func (server *Server) Initialize(transactionService *service.TransactionService) {
	server.TransactionService = transactionService
	server.Port = os.Getenv("TRANSACTION_API_SERVER_PORT")

	if server.Port == "" {
		log.Fatal("TRANSACTION_API_SERVER_PORT environment variable not set")
	}

	server.Router = server.InitializeRoutes()
}

func (server *Server) Run() {
	log.Info("Listening to port " + server.Port)
	log.Info("Transaction API Server Started")
	allow := os.Getenv("ENV_CORS_ALLOWED_ORIGIN")
	log.Info("Allow origin: " + allow)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{allow},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
	})
	handler := c.Handler(server.Router)
	log.Fatal(http.ListenAndServe(":"+server.Port, handlers.CompressHandler(handler)))
}
