# Fiber API Server

A clean, maintainable Go REST API server built with Fiber framework, featuring user authentication, point transfer system, comprehensive API documentation, and clean architecture.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/vorrarit/ai-coding-workshop-fiber.git
cd ai-coding-workshop-fiber

# Install dependencies
go mod tidy

# Run the server
go run .

# Access the API
curl http://localhost:3000/health

# View API documentation
open http://localhost:3000/swagger/
```

## 📋 Table of Contents

- [Features](#-features)
- [Project Structure](#-project-structure)
- [Architecture](#-architecture)
- [API Documentation](#-api-documentation)
- [API Endpoints](#-api-endpoints)
- [Getting Started](#-getting-started)
- [Configuration](#-configuration)
- [Development](#-development)
- [Deployment](#-deployment)
- [Documentation](#-documentation)

## ✨ Features

### 🔐 Security & Authentication
- ✅ JWT-based authentication with 24-hour expiry
- ✅ Password hashing with bcrypt (cost 14)
- ✅ SQL injection protection via GORM ORM
- ✅ Input validation and sanitization
- ✅ Bearer token authentication for protected routes
- ✅ Secure error handling without information leakage

### 💳 Point Transfer System
- ✅ Point balance management for users
- ✅ Secure point transfers between users via LBK codes
- ✅ Database transactions for atomic operations
- ✅ Insufficient balance validation
- ✅ Self-transfer prevention
- ✅ Complete transfer audit trail and history

### 🏗️ Architecture & Design
- ✅ Clean architecture with dependency injection
- ✅ Separation of concerns across layers
- ✅ Comprehensive error handling
- ✅ Environment-based configuration
- ✅ Structured logging and health checks
- ✅ Scalable and maintainable codebase

### 📚 Documentation & Testing
- ✅ Interactive Swagger/OpenAPI documentation
- ✅ Comprehensive API examples with curl commands
- ✅ Entity-Relationship diagrams (PlantUML)
- ✅ C4 Model architecture documentation
- ✅ Clean code with extensive comments

## 📁 Project Structure

```
.
├── main.go                          # Application entry point with dependency injection
├── docs/                            # Documentation and diagrams
│   ├── docs.go                     # Generated Swagger documentation
│   ├── swagger.json                # OpenAPI 2.0 JSON specification
│   ├── swagger.yaml                # OpenAPI 2.0 YAML specification
│   ├── er-diagram.md               # Entity-Relationship diagram (PlantUML)
│   └── architecture.md             # C4 Model architecture documentation
├── internal/                        # Private application code
│   ├── config/                      # Configuration management
│   │   └── config.go               # Environment and config loader
│   ├── database/                    # Database connection and setup
│   │   └── database.go             # Database initialization and migrations
│   ├── handlers/                    # HTTP request handlers (Presentation Layer)
│   │   ├── auth_handler.go         # Authentication endpoints
│   │   ├── health_handler.go       # Health and monitoring endpoints
│   │   ├── transfer_handler.go     # Point transfer endpoints
│   │   └── user_handler.go         # User management endpoints
│   ├── middleware/                  # Custom middleware
│   │   └── auth.go                 # JWT authentication middleware
│   ├── models/                      # Data models and DTOs
│   │   ├── requests.go             # Request DTOs with validation
│   │   ├── responses.go            # Response DTOs
│   │   └── user.go                 # Database models (User, Transfer)
│   ├── services/                    # Business logic layer
│   │   ├── transfer_service.go     # Point transfer business logic
│   │   └── user_service.go         # User management business logic
│   └── utils/                       # Utility functions
│       ├── auth.go                 # Password hashing utilities
│       └── jwt.go                  # JWT token utilities
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── .gitignore                      # Git ignore rules
├── POINT_TRANSFER_API.md           # API usage examples
└── README.md                       # This file
```

## 🏗️ Architecture

### Clean Architecture Principles

The application follows Clean Architecture principles with clear separation of concerns:

- **Presentation Layer** (`handlers/`): HTTP request/response handling with comprehensive input validation
- **Business Logic Layer** (`services/`): Core business rules and transaction management
- **Data Access Layer** (`database/`, `models/`): Database operations and entity definitions
- **Infrastructure Layer** (`config/`, `middleware/`, `utils/`): Cross-cutting concerns and utilities

### 📊 System Documentation

- **[ER Diagram](./docs/er-diagram.md)**: Database schema and entity relationships
- **[Architecture Diagrams](./docs/architecture.md)**: C4 Model system architecture documentation
- **[API Documentation](http://localhost:3000/swagger/)**: Interactive Swagger UI (when server is running)

## 📡 API Endpoints

| Method | Endpoint | Description | Auth Required | Documentation |
|--------|----------|-------------|---------------|---------------|
| POST | `/api/register` | User registration with email validation | ❌ | [Swagger](http://localhost:3000/swagger/) |
| POST | `/api/login` | User authentication with JWT token | ❌ | [Swagger](http://localhost:3000/swagger/) |
| GET | `/api/me` | Get current user profile from JWT | ✅ | [Swagger](http://localhost:3000/swagger/) |
| PUT | `/api/me` | Update current user profile | ✅ | [Swagger](http://localhost:3000/swagger/) |
| GET | `/api/users` | Search users by name or phone | ✅ | [Swagger](http://localhost:3000/swagger/) |
| POST | `/api/transfer` | Transfer points between users | ✅ | [Point Transfer Guide](./POINT_TRANSFER_API.md) |
| GET | `/api/transfer/history` | Get transfer history | ✅ | [Swagger](http://localhost:3000/swagger/) |
| GET | `/health` | Health check endpoint | ❌ | [Swagger](http://localhost:3000/swagger/) |

> 📚 **Complete API documentation** is available at `/swagger/` when the server is running

## 🚀 Quick Start

#### Middleware (`internal/middleware/`)
- **JWT Middleware**: Authentication and authorization
- Request/response logging (can be added)
- Rate limiting (can be added)

### Prerequisites

- Go 1.19 or higher
- SQLite (embedded database)

### Installation & Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd workshop5
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

4. **Access the services**
   - API Server: `http://localhost:3000`
   - Swagger Documentation: `http://localhost:3000/swagger/`
   - Health Check: `http://localhost:3000/health`

### Example Usage

1. **Register a new user**
   ```bash
   curl -X POST http://localhost:3000/api/register \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "password123",
       "first_name": "John",
       "last_name": "Doe",
       "phone_number": "0123456789",
       "dob": "1990-01-01"
     }'
   ```

2. **Login and get JWT token**
   ```bash
   curl -X POST http://localhost:3000/api/login \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "password123"
     }'
   ```

3. **Access protected endpoints**
   ```bash
   curl -X GET http://localhost:3000/api/me \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

## ⚙️ Configuration

### Environment Variables

```bash
# Database Configuration
DATABASE_PATH=users.db                    # SQLite database file path

# Security Configuration  
JWT_SECRET=your-super-secret-jwt-key      # JWT signing secret (change in production!)

# Server Configuration
```

### Default Configuration

The application uses sensible defaults if environment variables are not set:
- **JWT_SECRET**: `super-secret-jwt-key` (⚠️ **Change in production!**)
- **PORT**: `3000`
- **DATABASE_PATH**: `users.db` (auto-created in project root)

## 🛠️ Development

### Project Architecture

This project implements **Clean Architecture** with the following benefits:
- ✅ **Testability**: Each layer can be unit tested independently
- ✅ **Maintainability**: Clear separation of concerns and responsibilities
- ✅ **Scalability**: Easy to add new features without affecting existing code
- ✅ **Flexibility**: Database and frameworks can be easily swapped

### Key Design Patterns

- **Dependency Injection**: Services are injected into handlers for loose coupling
- **Repository Pattern**: Database operations are abstracted via GORM
- **DTO Pattern**: Separate request/response models from domain entities
- **Middleware Pattern**: Authentication and logging via Fiber middleware

### Building and Testing

```bash
# Install dependencies
go mod download

# Generate Swagger documentation
swag init

# Build the application
go build -o app main.go

# Run with development settings
go run main.go

# Run tests (when implemented)
go test ./...
```

## 🚀 Deployment

### Production Considerations

1. **Environment Variables**: Set secure values for production
   ```bash
   export JWT_SECRET="your-production-secret-key-with-at-least-32-characters"
   export DATABASE_PATH="/var/lib/app/users.db"
   export PORT="8080"
   ```

2. **Database**: Consider upgrading to PostgreSQL or MySQL for production
3. **Security**: Implement rate limiting, CORS, and HTTPS
4. **Monitoring**: Add structured logging and health monitoring
5. **Load Balancing**: Deploy multiple instances behind a load balancer

### Docker Deployment (Optional)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 3000
CMD ["./main"]
```

## 📚 Documentation

### Available Documentation

1. **[API Documentation (Swagger)](http://localhost:3000/swagger/)** - Interactive API documentation
2. **[Point Transfer Guide](./POINT_TRANSFER_API.md)** - Detailed API usage examples with curl commands
3. **[Entity-Relationship Diagram](./docs/er-diagram.md)** - Database schema visualization
4. **[Architecture Documentation](./docs/architecture.md)** - C4 Model system architecture
5. **[Generated API Specs](./docs/)** - OpenAPI 2.0 JSON/YAML specifications

### Development Guidelines

#### Code Organization
- Follow Clean Architecture principles
- Keep handlers thin (presentation logic only)
- Put business logic in services
- Use DTOs for API contracts
- Maintain clear separation between layers

#### Adding New Features
1. Define DTOs in `models/requests.go` and `models/responses.go`
2. Add business logic to appropriate service in `services/`
3. Create handler function in `handlers/`
4. Add Swagger annotations for documentation
5. Register routes in `main.go`
6. Update this README and API documentation

#### Security Guidelines
- Always validate input data
- Use parameterized queries (GORM handles this)
- Implement proper error handling
- Don't leak sensitive information in error messages
- Use strong JWT secrets in production

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Fast, Express-inspired web framework
- [GORM](https://gorm.io/) - Fantastic ORM for Go
- [Swagger](https://swagger.io/) - API documentation standard

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
