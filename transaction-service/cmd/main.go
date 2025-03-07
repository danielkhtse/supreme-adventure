package main

import (
	"github.com/danielkhtse/supreme-adventure/transaction-service/internal/api"
	"github.com/danielkhtse/supreme-adventure/transaction-service/internal/service"
)

func main() {

	// Initial Transaction Service, handle migrations
	transactionService := service.NewTransactionService()

	// Initialize Transactions API server
	var server api.Server
	server.Initialize(transactionService)
	server.Run()

	// TODO: Transaction Service GRPC Server
	// grpcServer := grpc.NewServer()
	// grpcServer.Serve(lis)
}
