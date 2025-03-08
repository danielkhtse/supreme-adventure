## API Documentation

This project provides API documentation through Swagger UI. You can access the documentation in the following ways:

### Using Docker Compose

1. Start all services using Docker Compose:

    ```
    docker-compose up -d
    ```

2. Access the Swagger UI at:

    ```
    http://localhost:3000
    ```

3. You can switch between the Account Service API and Transaction Service API using the dropdown menu at the top of the page.

### Generating Documentation Locally

1. Make sure you have the Swagger CLI installed:

    ```
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

2. Run the documentation generation script:

    ```
    ./generate-docs.sh
    ```

3. The documentation will be generated in:
    - `account-service/docs/`
    - `transaction-service/docs/`

### API Endpoints

#### Account Service (Port 8080)

-   `GET /health-check` - Health check endpoint
-   `POST /accounts` - Create a new account
-   `GET /accounts/{account_id}` - Get account details
-   `PUT /accounts/{account_id}/balance/transfer` - Transfer funds between accounts

#### Transaction Service (Port 8081)

-   `GET /health-check` - Health check endpoint
-   `POST /transactions` - Create a new transaction between accounts
