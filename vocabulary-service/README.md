# Vocabulary Service

A gRPC-based vocabulary management microservice for the Vocabulary Tracker application.

## Features

- Create, read, update, and delete vocabulary entries
- Search and filter vocabularies by date and text
- Vocabulary statistics and analytics
- User-based data isolation
- Pagination support
- PostgreSQL database integration
- gRPC API

## API Methods

### gRPC Service: VocabularyService

1. **GetVocabularies** - Get vocabularies with optional filtering
   - Request: `GetVocabulariesRequest` (user_id, date, search, limit, offset)
   - Response: `GetVocabulariesResponse` (vocabularies list, count, total)

2. **CreateVocabulary** - Create a new vocabulary entry
   - Request: `CreateVocabularyRequest` (user_id, word, meaning, example, date, status)
   - Response: `VocabularyResponse` (success, message, vocabulary)

3. **UpdateVocabulary** - Update an existing vocabulary entry
   - Request: `UpdateVocabularyRequest` (vocabulary_id, user_id, word, meaning, example, status)
   - Response: `VocabularyResponse` (success, message, vocabulary)

4. **DeleteVocabulary** - Delete a vocabulary entry
   - Request: `DeleteVocabularyRequest` (vocabulary_id, user_id)
   - Response: `DeleteVocabularyResponse` (success, message)

5. **GetVocabularyById** - Get vocabulary by ID
   - Request: `GetVocabularyByIdRequest` (vocabulary_id, user_id)
   - Response: `VocabularyResponse` (success, message, vocabulary)

6. **GetVocabularyStats** - Get vocabulary statistics
   - Request: `GetVocabularyStatsRequest` (user_id, date_from, date_to)
   - Response: `VocabularyStatsResponse` (total_words, words_this_week, words_this_month, status_counts, daily_counts)

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

The service will listen on port `:50052` for gRPC connections.

## Vocabulary Status Values

- `review_needed` - New vocabulary that needs to be reviewed
- `learning` - Currently being learned
- `learned` - Successfully learned
- `mastered` - Fully mastered

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

## Example Usage

See the `examples/client.go` file for sample gRPC client code demonstrating how to:
- Create vocabularies
- Search vocabularies
- Update vocabulary status
- Get vocabulary statistics

## Project Structure

```
vocabulary-service/
├── cmd/
│   └── main.go               # Service entry point
├── config/
│   └── config.go            # Configuration management
├── database/
│   └── vocab.database.go    # Database connection and migrations
├── models/
│   └── vocab.model.go       # Data models
├── proto/
│   ├── vocabulary.proto     # Protocol buffer definition
│   ├── vocabulary.pb.go     # Generated protobuf code
│   └── vocabulary_grpc.pb.go # Generated gRPC code
├── services/
│   └── vocabulary_service.go # gRPC service implementation
├── examples/
│   └── client.go            # Example gRPC client
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## API Integration

This vocabulary service is designed to work with:
- **Auth Service** (port 50051) - for user authentication
- **Main Backend** (HTTP REST API) - as a proxy to gRPC services
- **Frontend** - through the main backend

The auth service validates user tokens, and the vocabulary service handles all vocabulary-related operations using the authenticated user ID.