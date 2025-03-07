package main

import (
	"github.com/danielkhtse/supreme-adventure/account-service/internal/api"
	"github.com/danielkhtse/supreme-adventure/account-service/internal/service"
)

func main() {

	// Initial Account Service, handle migrations
	accountService := service.NewAccountService()

	// Initialize Accounts API server
	var server api.Server
	server.Initialize(accountService)
	server.Run()

	// TODO: Account Service GRPC Server
	// grpcServer := grpc.NewServer()
	// grpcServer.Serve(lis)
}
