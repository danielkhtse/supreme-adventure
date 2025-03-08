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

3. Start services:

    ```
    ./start-services.sh
    ```

4. Run unit tests:
    ```
    ./run-unit-test.sh
    ```
5. Generate Swagger documentation:

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

## Performance

### Test Environment

-   Tool: k6
-   Duration: 11 seconds (based on request rate)
-   Virtual Users: 50 concurrent users
-   Environment: Docker containers on local machine 8gb ram, 4 core

### Results Summary

#### Transaction Service

Test on `POST /transactions` - Create a new transaction between accounts

| Metric                | Value    |
| --------------------- | -------- |
| Average Response Time | 50.35ms  |
| 95th Percentile       | 251.23ms |
| Requests/sec          | 84.17    |
| Error Rate            | 0%       |
| Success Rate Checks   | 100%     |

Key Findings:

-   Response time threshold (p95 < 500ms) was met with 251.23ms at 95th percentile
-   No errors observed during the test
-   Perfect check pass rate at 100%
-   Transaction duration averaged 52.28ms with 95th percentile at 136.85ms
-   System maintained stable performance under load with 50 concurrent users

### Performance Test Reports

#### Run performance test

```
./run-performance-test.sh
```

Detailed performance test reports are located in `./performance-test-reports`

# TODO and decisions

-   **UUID v6 Implementation**: Migrate account and transaction identifiers to UUID v6 for improved chronological sorting while maintaining uniqueness
-   **gRPC Communication**: Implement gRPC for inter-service communication to enhance performance
-   **Asynchronous Transaction Processing**: Implement asynchronous transaction processing using message queues (RabbitMQ/Kafka) to:

    -   Improve throughput by processing transactions in parallel
    -   Enable better scalability and fault tolerance
    -   Provide transaction status updates via webhooks
    -   Support eventual consistency model for transaction reconciliation

-   **Database Migration Strategy**: Replace GORM auto-migration with explicit migration scripts for better version control and schema management
-   **Error Handling and Logging**: Add detailed error messages and logging for better debugging
-   **API Security**: Implement API key authentication and rate limiting for secure usage
-   **Performance Optimization**: Add caching for frequently accessed data e.g. redis
-   **Testing Strategy**: Add integration tests and API tests for critical functionality
