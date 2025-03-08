#!/bin/bash

# Initialize go work if it doesn't exist
if [ ! -f go.work ]; then
    go work init account-service transaction-service common
fi

# Run go mod tidy for each service
echo "Running go mod tidy for all services..."
(cd account-service && go mod tidy)
(cd transaction-service && go mod tidy)
(cd common && go mod tidy)

# Start both services
echo "Starting services..."
(go run account-service/cmd/main.go) & 
(go run transaction-service/cmd/main.go) &

# Wait for both background processes
wait