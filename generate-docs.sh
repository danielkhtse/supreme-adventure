#!/bin/bash

# Exit on any error
set -e

# Function to generate swagger docs for a service
generate_swagger_docs() {
    local service=$1
    echo "Generating swagger documentation for ${service}..."
    
    # Create docs directory if it doesn't exist
    mkdir -p ./${service}/docs
    
    # Generate swagger spec
    cd ./${service}
    
    # Check if the directory exists and has Go files
    if [ ! -d "./internal/api" ]; then
        echo "Warning: ./internal/api directory does not exist for ${service}"
        # Try to find where the API handlers might be located
        possible_dirs=$(find . -type d -name "api" -o -name "handler" -o -name "handlers" | grep -v "vendor")
        if [ -n "$possible_dirs" ]; then
            echo "Found possible API directories:"
            echo "$possible_dirs"
            echo "Please update the script to use the correct path."
        fi
    else
        # List files to debug
        echo "Go files in ./internal/api directory:"
        find ./internal/api -name "*.go" | sort
        
        # Run swag init with the correct parameters
        echo "Running swag init..."
        swag init \
            --dir ./internal/api \
            --output ./docs \
            --generalInfo server.go \
            --parseDependency \
            --parseInternal
            
        echo "Swagger documentation generated in ./${service}/docs"
    fi
    
    # Return to the root directory
    cd ..
}

# Install swag if not already installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag CLI..."
    go install github.com/swaggo/swag/cmd/swag@latest
    # Use the full path to swag after installation
    export PATH="$PATH:$(go env GOPATH)/bin"
fi

# Generate docs for each service
generate_swagger_docs "account-service"
generate_swagger_docs "transaction-service"

# Create swagger-config.json if it doesn't exist
if [ ! -f "swagger-config.json" ]; then
    echo "Creating swagger-config.json..."
    cat > swagger-config.json << EOF
{
  "urls": [
    {
      "name": "Account Service API",
      "url": "/account-docs/swagger.json"
    },
    {
      "name": "Transaction Service API",
      "url": "/transaction-docs/swagger.json"
    }
  ],
  "deepLinking": true,
  "displayRequestDuration": true,
  "defaultModelsExpandDepth": 3,
  "defaultModelExpandDepth": 3,
  "defaultModelRendering": "example",
  "docExpansion": "list",
  "showExtensions": true,
  "showCommonExtensions": true,
  "supportedSubmitMethods": [
    "get",
    "put",
    "post",
    "delete",
    "options",
    "head",
    "patch",
    "trace"
  ]
}
EOF
fi

echo "Documentation generation complete!"
