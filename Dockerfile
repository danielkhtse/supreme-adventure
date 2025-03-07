# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.work ./
COPY account-service/go.mod account-service/
COPY common/go.mod common/

# Download dependencies
RUN go mod download

# Copy source code
COPY account-service/ account-service/
COPY common/ common/

# Build the application
WORKDIR /app/account-service
RUN CGO_ENABLED=0 GOOS=linux go build -o account-service ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/account-service/account-service .

# Copy environment file
COPY .env .env

# Expose port
EXPOSE 8080

# Run the application
CMD ["./account-service"]