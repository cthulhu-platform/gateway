# Builder stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy workspace files
COPY go.work go.work.sum ./

# Copy common module (required dependency)
COPY common ./common

# Copy gateway module
COPY gateway/MAIN ./gateway/MAIN

# Set working directory to gateway
WORKDIR /build/gateway/MAIN

# Download dependencies
RUN go mod download

# Build the binary (modernc.org/sqlite doesn't require CGO)
RUN GOOS=linux go build -o /build/api ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Create db directory for SQLite databases
RUN mkdir -p /app/db

# Copy the binary from builder
COPY --from=builder /build/api .

# Expose port
EXPOSE 7777

# Run the binary
CMD ["./api"]

