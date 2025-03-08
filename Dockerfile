# Build stage
FROM golang:1.24.1-alpine AS builder

ARG SERVICE_NAME

WORKDIR /app

# Copy go mod files
COPY go.mod go.work ./
COPY account-service/go.mod account-service/
COPY transaction-service/go.mod transaction-service/
COPY common/go.mod common/

# Download dependencies
RUN go mod download

# Copy source code
COPY account-service/ account-service/
COPY transaction-service/ transaction-service/
COPY common/ common/

# Build the application
WORKDIR /app/${SERVICE_NAME}
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/main.go

# Final stage
FROM alpine:3.19

ARG SERVICE_NAME

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/${SERVICE_NAME}/service .

# Copy environment file
COPY .env.docker .env

# Expose port
EXPOSE 8080

# Run the application
CMD ["./service"]