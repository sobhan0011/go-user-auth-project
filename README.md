# OTP Auth Service (Go, Clean Architecture)

A backend service implementing OTP-based authentication and user management with clean architecture principles.

## Features

- **OTP Authentication**: Phone-based login/registration with 6-digit OTPs
- **Rate Limiting**: 3 OTP requests per phone number within 10 minutes
- **User Management**: CR~~UD~~ operations with pagination and search
- **JWT Tokens**: Secure authentication with configurable TTL
- **Clean Architecture**: Domain, Usecase, and Infrastructure layers
- **API Documentation**: OpenAPI 3.0 spec with Swagger UI
- **Containerized**: Docker and docker-compose ready

## Architecture

```
internal/
├── domain/          # Business entities and interfaces
│   └── user/        # User domain models
├── usecase/         # Business logic
│   ├── auth/        # Authentication use cases
│   └── user/        # User management use cases
├── infra/           # External dependencies
│   ├── db/          # Database implementations
│   └── cache/       # Cache implementations
└── http/            # HTTP layer
    ├── handlers/    # Request handlers
    └── middleware/  # HTTP middleware
```

## Database Choice

**PostgreSQL** for user data:
- Robust indexing for search performance on phone numbers
- ACID compliance for data consistency
- Mature ecosystem with excellent Go support
- JSON support for future extensibility

**Redis** for temporary data:
- Fast TTL-based expiration for OTPs
- Atomic counters for rate limiting
- High performance for caching

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose

### Run with Docker

```bash
docker-compose up --build
```

## API Documentation

- **Swagger UI**: http://localhost:8080/docs/
- **OpenAPI Spec**: http://localhost:8080/openapi.yaml

## API Examples

### Request OTP
```bash
curl -X POST http://localhost:8080/api/auth/request-otp \
  -H 'Content-Type: application/json' \
  -d '{"phone":"+15551234567"}'
```

**Response:**
```json
{
  "message": "otp_sent"
}
```

### Verify OTP
```bash
curl -X POST http://localhost:8080/api/auth/verify-otp \
  -H 'Content-Type: application/json' \
  -d '{"phone":"+15551234567","code":"123456"}'
```

**Response:**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "phone": "+15551234567",
      "created_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

### Get User
```bash
curl http://localhost:8080/api/users/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### List Users
```bash
curl "http://localhost:8080/api/users?phone=+1555&page=1&limit=20" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

**Response:**
```json
{
  "data": {
    "items": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "phone": "+15551234567",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 20
  }
}
```

## Testing

Run unit tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | `8080` | Server port |
| `JWT_SECRET` | Required | JWT signing secret |
| `POSTGRES_DSN` | `postgres://postgres:postgres@localhost:5432/otpapp?sslmode=disable` | PostgreSQL connection string |
| `REDIS_ADDR` | `localhost:6379` | Redis address |

## Rate Limiting

- **OTP Requests**: 3 per phone number per 10 minutes
- **Implementation**: Redis-based sliding window
- **Headers**: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`

## Security Features

- **Phone Validation**: E.164 format validation
- **JWT Tokens**: HS256 signing with configurable TTL
- **Rate Limiting**: Prevents OTP abuse
- **CORS**: Configurable cross-origin requests
- **Input Validation**: Request payload validation

## Development

### Project Structure
- **Domain Layer**: Pure business entities and interfaces
- **Usecase Layer**: Business logic and orchestration
- **Infrastructure Layer**: External dependencies (DB, Cache, HTTP)
- **Clean Dependencies**: Domain ← Usecase ← Infrastructure

### Adding New Features
1. Define domain entities and interfaces
2. Implement use cases with business logic
3. Create infrastructure implementations
4. Wire dependencies in main.go
5. Add HTTP handlers and routes
6. Write tests for each layer

## Production Considerations

- Set strong `JWT_SECRET` environment variable
- Configure proper CORS origins
- Use HTTPS in production