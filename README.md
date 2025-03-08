# Supreme Adventure

A microservices-based banking system with account management and transaction processing capabilities.

## Overview

This project implements a simple banking system with two main microservices:

-   **Account Service**: Manages bank accounts, including creation, retrieval, and balance transfers
-   **Transaction Service**: Handles transactions between accounts

The services are built using Go and communicate with each other via REST APIs. Data is stored in a PostgreSQL database.

## Assumptions

-   Support single currency USD only
-   All amounts, balances, and transaction amounts are in cents (smallest unit of USD)

## Features

-   Create and manage bank accounts
-   Transfer funds between accounts
-   View account details and transaction history
-   API documentation with Swagger UI
-   Containerized deployment with Docker

## Architecture

The system consists of the following components:

-   **Account Service**: Handles account-related operations (port 8080)
-   **Transaction Service**: Manages transactions between accounts (port 8081)
-   **PostgreSQL Database**: Stores account and transaction data
-   **Swagger UI**: Provides API documentation (port 3000)

## Getting Started

### Prerequisites

-   Docker and Docker Compose
-   Go 1.24 or later (for local development)

### Running with Docker Compose

1. Clone the repository:

    ```
    git clone https://github.com/danielkhtse/supreme-adventure.git
    cd supreme-adventure
    ```

2. Start all services:

    ```
    docker-compose up -d
    ```

3. Access the services:
    - Account Service: http://localhost:8080
    - Transaction Service: http://localhost:8081
    - Swagger UI: http://localhost:3000

### Local Development

1. Install dependencies:

    ```
    go mod download
    ```

2. Create your own .env file

   ```
   cp .env.sample .env
   ```
   
4. Start services:

    ```
    ./start-services.sh
    ```

5. Run unit tests:
    ```
    ./run-unit-test.sh
    ```
6. Generate Swagger documentation:

    ```
    ./generate-docs.sh
    ```

## API Documentation

For detailed API documentation, please refer to [README-api-docs.md](README-api-docs.md).

## Testing

Run unit tests with:

```
./run-unit-test.sh
```

Test reports located in `./test-reports`

## TODO and decisions

-   **UUID v6 Implementation**: Migrate account and transaction identifiers to UUID v6 for improved chronological sorting while maintaining uniqueness
-   **gRPC Communication**: Implement gRPC for inter-service communication to enhance performance
-   **Database Migration Strategy**: Replace GORM auto-migration with explicit migration scripts for better version control and schema management
-   **Error Handling and Logging**: Add detailed error messages and logging for better debugging
-   **API Security**: Implement API key authentication and rate limiting for secure usage
-   **Performance Optimization**: Add caching for frequently accessed data e.g. redis
-   **Testing Strategy**: Add integration tests and API tests for critical functionality
