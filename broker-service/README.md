# Broker Service

The broker service acts as a gateway between the frontend and the microservices (auth-service, vocabulary-service). It handles HTTP requests from the frontend and forwards them as gRPC calls to the appropriate services.

## Features

- HTTP to gRPC gateway
- User registration and login endpoints
- CORS support for frontend communication
- JSON request/response handling

## Endpoints

### Authentication Endpoints

#### POST /auth/register
Registers a new user.

**Request Body:**
```json
{
    "email": "user@example.com",
    "password": "password123"
}
```

**Response:**
```json
{
    "success": true,
    "message": "User registered successfully",
    "token": "jwt_token_here",
    "user": {
        "id": 1,
        "email": "user@example.com",
        "created_at": "2025-09-27T10:00:00Z"
    }
}
```

#### POST /auth/login
Authenticates a user.

**Request Body:**
```json
{
    "email": "user@example.com",
    "password": "password123"
}
```

**Response:**
```json
{
    "success": true,
    "message": "Login successful",
    "token": "jwt_token_here",
    "user": {
        "id": 1,
        "email": "user@example.com",
        "created_at": "2025-09-27T10:00:00Z"
    }
}
```

### Vocabulary Endpoints (Requires Authentication)

All vocabulary endpoints require a valid JWT token in the Authorization header: `Authorization: Bearer <token>`

**Status Values:**
- `review_needed` - Default status for new vocabulary entries
- `learned` - User has learned the word but needs occasional review  
- `mastered` - User has fully mastered the word

#### GET /vocab
Get vocabularies for the authenticated user.

**Query Parameters:**
- `date` (optional): Filter by date (YYYY-MM-DD)
- `search` (optional): Search term
- `limit` (optional): Limit results (default: 50)
- `offset` (optional): Pagination offset (default: 0)

**Response:**
```json
{
    "success": true,
    "message": "Vocabularies retrieved successfully",
    "vocabularies": [
        {
            "id": 1,
            "user_id": 1,
            "word": "serendipity",
            "meaning": "pleasant surprise or fortunate discovery",
            "example": "It was serendipity that we met at the coffee shop.",
            "date": "2025-09-27",
            "status": "review_needed",
            "created_at": "2025-09-27T10:00:00Z",
            "updated_at": "2025-09-27T10:00:00Z"
        }
    ],
    "count": 1,
    "total": 10
}
```

#### POST /vocab
Create a new vocabulary entry.

**Request Body:**
```json
{
    "word": "serendipity",
    "meaning": "pleasant surprise or fortunate discovery",
    "example": "It was serendipity that we met at the coffee shop.",
    "date": "2025-09-27",
    "status": "review_needed"
}
```

**Response:**
```json
{
    "success": true,
    "message": "Vocabulary created successfully",
    "vocabulary": {
        "id": 1,
        "user_id": 1,
        "word": "serendipity",
        "meaning": "pleasant surprise or fortunate discovery",
        "example": "It was serendipity that we met at the coffee shop.",
        "date": "2025-09-27",
        "status": "review_needed",
        "created_at": "2025-09-27T10:00:00Z",
        "updated_at": "2025-09-27T10:00:00Z"
    }
}
```

#### PUT /vocab/{id}
Update an existing vocabulary entry.

**Request Body:**
```json
{
    "word": "serendipity",
    "meaning": "updated meaning",
    "example": "updated example",
    "status": "learned"
}
```

**Response:** Same as POST /vocab

#### DELETE /vocab/{id}
Delete a vocabulary entry.

**Response:**
```json
{
    "success": true,
    "message": "Vocabulary deleted successfully"
}
```

## Running the Service

### Prerequisites
- Go 1.23.1 or higher
- Auth service running on localhost:50051

### Build and Run
```bash
# Build the service
make build

# Run the service
make run

# Or run directly
go run cmd/main.go
```

The service will start on port 8080.

### Configuration
The service connects to:
- Auth Service: `localhost:50051` (gRPC)
- Vocabulary Service: `localhost:50052` (gRPC)

## Development

### Project Structure
```
broker-service/
├── cmd/
│   └── main.go          # Main entry point
├── config/
│   └── config.go        # Configuration and gRPC connections
├── routes/
│   ├── routes.go        # HTTP router setup
│   └── auth.route.go    # Auth route handlers
├── proto/               # Protocol buffer definitions
├── bin/                 # Built binaries
└── Makefile            # Build scripts
```

### Building
```bash
make build
```

### Cleaning
```bash
make clean
```