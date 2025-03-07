package api

import (
	"net/http"
	"os"

	"github.com/danielkhtse/supreme-adventure/account-service/internal/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// Accounts API server
type Server struct {
	AccountService *service.AccountService
	Router         *mux.Router
	Port           string
}

func (server *Server) Initialize(accountService *service.AccountService) {
	server.AccountService = accountService
	server.Port = os.Getenv("ACCOUNT_API_SERVER_PORT")

	if server.Port == "" {
		log.Fatal("ACCOUNT_API_SERVER_PORT environment variable not set")
	}

	server.Router = server.InitializeRoutes()
}

func (server *Server) Run() {
	log.Info("Listening to port " + server.Port)
	log.Info("Application Started")
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
