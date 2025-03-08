#!/bin/bash

# Create test reports directory if it doesn't exist
mkdir -p test-reports

# Define services to test - check if directories exist first
services=()
if [ -d "account-service" ]; then
    services+=("account")
fi
if [ -d "transaction-service" ]; then
    services+=("transaction")
fi

# Initialize coverage files
for service in "${services[@]}"; do
    echo "mode: set" > "test-reports/${service}-coverage.out"
done

# Run unit tests for each service
for service in "${services[@]}"; do
    echo "Running ${service} service unit tests..."
    if cd "${service}-service"; then
        # Run tests and handle potential failures
        if ! go test -v -run "^TestUnit" -coverprofile="../test-reports/${service}-coverage.out" ./...; then
            echo "Warning: Tests failed for ${service} service"
        fi
        cd ..
    else
        echo "Error: ${service}-service directory not found"
    fi
done

# Generate HTML coverage reports
echo "Generating coverage reports..."
for service in "${services[@]}"; do
    if [ -f "test-reports/${service}-coverage.out" ]; then
        go tool cover -html="test-reports/${service}-coverage.out" -o "test-reports/${service}-coverage.html"
    else
        echo "Warning: Coverage file for ${service} not found"
    fi
done

echo "Coverage reports generated in test-reports directory"
