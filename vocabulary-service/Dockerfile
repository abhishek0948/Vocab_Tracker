# Build stage
FROM golang:1.23.1-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (needed for go modules)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vocabulary-service cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates and netcat for health checks
RUN apk --no-cache add ca-certificates netcat-openbsd

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/vocabulary-service .

# Expose port
EXPOSE 50052

# Command to run
CMD ["./vocabulary-service"]