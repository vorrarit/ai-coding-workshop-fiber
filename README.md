# Fiber API Server - Refactored

A clean, maintainable Go REST API server built with Fiber framework, featuring user authentication, point transfer system, and clean architecture.

## Project Structure

```
.
├── main.go                          # Application entry point
├── internal/                        # Private application code
│   ├── config/                      # Configuration management
│   │   └── config.go               # Environment and config loader
│   ├── database/                    # Database connection and setup
│   │   └── database.go             # Database initialization
│   ├── handlers/                    # HTTP request handlers
│   │   ├── auth_handler.go         # Authentication endpoints
│   │   ├── health_handler.go       # Health and hello endpoints
│   │   ├── transfer_handler.go     # Point transfer endpoints
│   │   └── user_handler.go         # User management endpoints
│   ├── middleware/                  # Custom middleware
│   │   └── auth.go                 # JWT authentication middleware
│   ├── models/                      # Data models and DTOs
│   │   ├── requests.go             # Request DTOs
│   │   ├── responses.go            # Response DTOs
│   │   └── user.go                 # Database models
│   ├── services/                    # Business logic layer
│   │   ├── transfer_service.go     # Point transfer business logic
│   │   └── user_service.go         # User management business logic
│   └── utils/                       # Utility functions
│       ├── auth.go                 # Password hashing utilities
│       └── jwt.go                  # JWT token utilities
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── .gitignore                      # Git ignore rules
└── POINT_TRANSFER_API.md           # API documentation
```

## Architecture Overview

### Clean Architecture Principles

This project follows clean architecture principles with clear separation of concerns:

1. **Presentation Layer** (`handlers/`): HTTP request/response handling
2. **Business Logic Layer** (`services/`): Core business rules and operations
3. **Data Access Layer** (`database/`): Database operations and migrations
4. **Infrastructure Layer** (`middleware/`, `utils/`): Cross-cutting concerns

### Key Components

#### Configuration (`internal/config/`)
- Centralized configuration management
- Environment variable support
- Sensible defaults

#### Models (`internal/models/`)
- **Database Models**: User, Transfer entities with GORM tags
- **Request DTOs**: Input validation and parsing
- **Response DTOs**: Consistent API responses

#### Services (`internal/services/`)
- **UserService**: User registration, authentication, profile management
- **TransferService**: Point transfers with transaction safety

#### Handlers (`internal/handlers/`)
- **AuthHandler**: Registration and login endpoints
- **UserHandler**: User profile and search operations
- **TransferHandler**: Point transfer operations
- **HealthHandler**: System health checks

#### Middleware (`internal/middleware/`)
- **JWT Middleware**: Authentication and authorization
- Request/response logging (can be added)
- Rate limiting (can be added)

#### Utilities (`internal/utils/`)
- **Authentication**: Password hashing and verification
- **JWT**: Token generation and parsing
- **Code Generation**: LBK code generation

## API Endpoints

### Public Endpoints
- `GET /api/hello` - Hello world endpoint
- `GET /health` - Health check
- `POST /register` - User registration
- `POST /login` - User authentication

### Protected Endpoints (Require JWT)
- `GET /me` - Get user profile
- `GET /points/balance` - Get point balance
- `GET /users/search` - Search users by LBK code
- `POST /points/transfer` - Transfer points
- `GET /points/history` - Get transfer history

## Features

### Security
- ✅ JWT-based authentication
- ✅ Password hashing with bcrypt
- ✅ SQL injection protection via GORM
- ✅ Input validation
- ✅ Error handling without information leakage

### Data Integrity
- ✅ Database transactions for transfers
- ✅ Point balance validation
- ✅ Duplicate prevention
- ✅ Audit trail via transfer records

### Scalability
- ✅ Clean architecture for easy testing
- ✅ Dependency injection
- ✅ Separation of concerns
- ✅ Environment-based configuration

## Environment Variables

```bash
# Database
DATABASE_PATH=users.db

# Security
JWT_SECRET=your-super-secret-jwt-key

# Server
PORT=3000
```

## Getting Started

### Prerequisites
- Go 1.21 or higher
- SQLite (embedded)

### Installation

1. Clone the repository
```bash
git clone <repository-url>
cd fiber-api
```

2. Install dependencies
```bash
go mod tidy
```

3. Set environment variables (optional)
```bash
export JWT_SECRET="your-secret-key"
export PORT="3000"
```

4. Run the application
```bash
go run .
```

### Building for Production

```bash
go build -o fiber-api .
./fiber-api
```

## Testing

### Manual Testing
See `POINT_TRANSFER_API.md` for complete API documentation with curl examples.

### Unit Testing (Future Enhancement)
The clean architecture makes it easy to add unit tests:

```bash
go test ./internal/services/...
go test ./internal/handlers/...
go test ./internal/utils/...
```

## Development Guidelines

### Adding New Features

1. **Add Models**: Define request/response DTOs in `internal/models/`
2. **Business Logic**: Implement in appropriate service in `internal/services/`
3. **HTTP Layer**: Create handler in `internal/handlers/`
4. **Wire Up**: Register routes in `main.go`

### Code Style
- Follow Go conventions and formatting
- Use dependency injection
- Keep handlers thin - business logic belongs in services
- Handle errors appropriately at each layer

### Database Changes
- Add new models to `internal/models/`
- Update database migration in `internal/database/database.go`
- Test locally before deployment

## Deployment

### Docker (Future Enhancement)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o fiber-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/fiber-api .
CMD ["./fiber-api"]
```

### Environment Setup
- Use environment variables for sensitive configuration
- Set appropriate JWT secrets in production
- Configure database path and backup strategy
- Set up monitoring and logging

## Contributing

1. Follow the established architecture patterns
2. Add appropriate error handling
3. Update documentation for new endpoints
4. Test thoroughly before submitting PR

## License

This project is licensed under the MIT License.
