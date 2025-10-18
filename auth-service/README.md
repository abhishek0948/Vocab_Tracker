# Auth Service

A gRPC-based authentication microservice for the Vocabulary Tracker application.

## Features

- User registration and login
- JWT token generation and validation
- Password hashing with bcrypt
- PostgreSQL database integration
- gRPC API

## API Methods

### gRPC Service: AuthService

1. **Register** - Register a new user
   - Request: `RegisterRequest` (email, password)
   - Response: `AuthResponse` (success, message, token, user)

2. **Login** - Authenticate existing user
   - Request: `LoginRequest` (email, password)
   - Response: `AuthResponse` (success, message, token, user)

3. **ValidateToken** - Validate JWT token
   - Request: `ValidateTokenRequest` (token)
   - Response: `ValidateTokenResponse` (valid, message, user_id, email)

4. **GetProfile** - Get user profile by ID
   - Request: `GetProfileRequest` (user_id)
   - Response: `UserResponse` (success, message, user)

## Configuration

The service uses environment variables for configuration:

- `DB_HOST` - Database host (default: localhost)
- `DB_USER` - Database user (default: vocab_user)
- `DB_PASSWORD` - Database password (default: vocab_password)
- `DB_NAME` - Database name (default: vocab_tracker)
- `DB_PORT` - Database port (default: 5432)
- `JWT_SECRET` - JWT signing secret (default: your-secret-key-change-in-production)

## Running the Service

### Prerequisites

1. PostgreSQL database running
2. Go 1.23.1 or later
3. Protocol Buffers compiler (protoc)

### Build and Run

```bash
# Install dependencies
go mod tidy

# Generate protobuf files (if needed)
make proto

# Build the service
go build ./cmd/main.go

# Run the service
./main
```

The service will listen on port `:50051` for gRPC connections.

## Development

### Regenerating Protobuf Files

```bash
make proto
```

### Building

```bash
make build
```

### Running

```bash
make run
```

## Project Structure

```
auth-service/
├── cmd/
│   └── main.go          # Service entry point
├── config/
│   └── config.go        # Configuration management
├── controllers/
│   └── auth.controller.go # HTTP controllers (legacy)
├── database/
│   └── auth.database.go # Database connection and migrations
├── middleware/
│   └── auth.middleware.go # JWT middleware
├── models/
│   └── auth.model.go    # Data models
├── proto/
│   ├── auth.proto       # Protocol buffer definition
│   ├── auth.pb.go       # Generated protobuf code
│   └── auth_grpc.pb.go  # Generated gRPC code
├── routes/
│   └── auth.route.go    # HTTP routes (legacy)
├── services/
│   └── auth_service.go  # gRPC service implementation
├── go.mod
├── go.sum
├── Makefile
└── README.md
```