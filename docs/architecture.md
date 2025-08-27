# Architecture Documentation - C4 Model

This document describes the architecture of the Fiber API Server using the C4 model approach, providing multiple levels of architectural views from system context to implementation details.

## Overview

The Fiber API Server is a clean, maintainable Go REST API built with the Fiber framework, featuring user authentication and point transfer system. The architecture follows clean architecture principles with clear separation of concerns.

## Level 1: System Context Diagram

```plantuml
@startuml SystemContext
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Context.puml

title System Context Diagram - Fiber API Server

Person(user, "API User", "Mobile app, web client, or external system that consumes the API")

System(fiberApi, "Fiber API Server", "Go-based REST API providing user authentication and point transfer services")

System_Ext(database, "SQLite Database", "Embedded database storing user accounts and transfer transactions")

Rel(user, fiberApi, "Makes API calls", "HTTPS/JSON")
Rel(fiberApi, database, "Reads from and writes to", "SQL")

note as N1
  **Key Features:**
  - User registration and authentication
  - JWT-based security
  - Point balance management
  - Point transfer between users
  - Transfer history tracking
  - API documentation with Swagger
end note

@enduml
```

## Level 2: Container Diagram

```plantuml
@startuml ContainerDiagram
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml

title Container Diagram - Fiber API Server

Person(user, "API User", "Mobile app, web client, or external system")

Container_Boundary(apiServer, "Fiber API Server") {
    Container(webApp, "Web Application", "Go, Fiber Framework", "Provides REST API endpoints with JWT authentication")
    Container(swaggerUI, "Swagger UI", "Static Web Interface", "Interactive API documentation and testing interface")
}

ContainerDb(database, "SQLite Database", "SQLite", "Stores user accounts, authentication data, and point transfer transactions")

Rel(user, webApp, "Makes API calls", "HTTPS/JSON")
Rel(user, swaggerUI, "Views API docs & tests", "HTTPS")
Rel(webApp, database, "Reads from and writes to", "SQL/GORM")
Rel(swaggerUI, webApp, "Calls API endpoints", "HTTP/JSON")

note as N1
  **API Endpoints:**
  - Authentication: /register, /login
  - User: /me, /points/balance, /users/search
  - Transfers: /points/transfer, /points/history
  - Health: /health, /api/hello
  - Documentation: /swagger/*
end note

@enduml
```

## Level 3: Component Diagram

```plantuml
@startuml ComponentDiagram
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title Component Diagram - Fiber API Server

Container(user, "API Client", "External system making API calls")

Container_Boundary(fiberApi, "Fiber API Server") {
    Component(router, "Fiber Router", "Go Fiber", "HTTP request routing and middleware")
    Component(middleware, "Middleware Layer", "Go", "CORS, JWT authentication, request validation")
    
    Component(authHandler, "Auth Handler", "Go", "User registration and login endpoints")
    Component(userHandler, "User Handler", "Go", "User profile and search endpoints")
    Component(transferHandler, "Transfer Handler", "Go", "Point transfer and history endpoints")
    Component(healthHandler, "Health Handler", "Go", "Health check and monitoring endpoints")
    
    Component(userService, "User Service", "Go", "User business logic and authentication")
    Component(transferService, "Transfer Service", "Go", "Point transfer business logic with transaction safety")
    
    Component(middleware_auth, "JWT Middleware", "Go", "Token validation and user context")
    Component(database, "Database Layer", "GORM", "Database connection and ORM operations")
    Component(models, "Data Models", "Go Structs", "User, Transfer, and DTO definitions")
    Component(utils, "Utilities", "Go", "JWT token generation, password hashing, LBK code generation")
}

ContainerDb(sqlite, "SQLite Database", "SQLite", "User and transfer data storage")

' External connections
Rel(user, router, "HTTP requests", "JSON/REST")

' Router to handlers
Rel(router, middleware, "Request processing", "")
Rel(middleware, authHandler, "Auth requests", "")
Rel(middleware, userHandler, "User requests", "")
Rel(middleware, transferHandler, "Transfer requests", "")
Rel(middleware, healthHandler, "Health requests", "")

' Handlers to services
Rel(authHandler, userService, "User operations", "")
Rel(userHandler, userService, "User operations", "")
Rel(transferHandler, transferService, "Transfer operations", "")

' Services to database
Rel(userService, database, "User data access", "SQL")
Rel(transferService, database, "Transfer data access", "SQL")

' Cross-cutting concerns
Rel(middleware, middleware_auth, "JWT validation", "")
Rel(userService, utils, "Password hashing, JWT", "")
Rel(transferService, utils, "LBK generation", "")
Rel(database, models, "Data mapping", "")
Rel(database, sqlite, "SQL operations", "GORM")

@enduml
```

## Level 4: Code Diagram - Clean Architecture Layers

```plantuml
@startuml CodeDiagram
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title Code Diagram - Clean Architecture Implementation

Container_Boundary(presentation, "Presentation Layer") {
    Component(main, "main.go", "Application Entry Point", "Dependency injection and server startup")
    Component(handlers, "handlers/", "HTTP Handlers", "Request/response processing")
    Component(middleware, "middleware/", "Middleware", "Cross-cutting concerns")
}

Container_Boundary(business, "Business Logic Layer") {
    Component(services, "services/", "Business Services", "Core business rules and operations")
    Component(models, "models/", "Domain Models", "Business entities and DTOs")
}

Container_Boundary(infrastructure, "Infrastructure Layer") {
    Component(database, "database/", "Database Access", "GORM setup and migrations")
    Component(config, "config/", "Configuration", "Environment and app configuration")
    Component(utils, "utils/", "Utilities", "Helper functions and tools")
}

Container_Boundary(external, "External Dependencies") {
    Component(fiber, "Fiber Framework", "Web Framework", "HTTP server and routing")
    Component(gorm, "GORM", "ORM", "Database operations")
    Component(jwt, "JWT Library", "Authentication", "Token generation and validation")
    Component(sqlite, "SQLite", "Database", "Data persistence")
}

' Layer dependencies (Clean Architecture - dependency rule)
Rel_D(presentation, business, "Depends on", "Interface contracts")
Rel_D(business, infrastructure, "Depends on", "Implementation details")
Rel_D(infrastructure, external, "Depends on", "External libraries")

' Specific dependencies
Rel(main, handlers, "Initializes", "")
Rel(main, services, "Initializes", "")
Rel(main, config, "Loads", "")
Rel(handlers, services, "Uses", "Business operations")
Rel(services, models, "Uses", "Domain objects")
Rel(services, database, "Uses", "Data access")
Rel(database, gorm, "Uses", "ORM operations")
Rel(handlers, fiber, "Uses", "HTTP handling")
Rel(utils, jwt, "Uses", "Token operations")

note as N1
  **Clean Architecture Principles:**
  
  1. **Dependency Rule**: Inner layers never depend on outer layers
  2. **Separation of Concerns**: Each layer has specific responsibilities
  3. **Interface Segregation**: Layers communicate through abstractions
  4. **Dependency Injection**: Dependencies are injected at startup
  
  **Benefits:**
  - Testable business logic
  - Framework independence
  - Database independence
  - Easy to maintain and extend
end note

@enduml
```

## Architectural Patterns & Principles

### Clean Architecture Implementation

The Fiber API Server follows Clean Architecture principles:

1. **Presentation Layer** (`handlers/`, `middleware/`, `main.go`)
   - HTTP request/response handling
   - Input validation and serialization
   - Authentication middleware
   - Dependency wiring

2. **Business Logic Layer** (`services/`, `models/`)
   - Core business rules
   - Domain entities and operations
   - Business validation
   - Use case orchestration

3. **Infrastructure Layer** (`database/`, `config/`, `utils/`)
   - External system integrations
   - Database access patterns
   - Configuration management
   - Utility functions

### Design Patterns Used

- **Dependency Injection**: Services and handlers are wired at startup
- **Repository Pattern**: Database layer abstracts data access
- **Service Layer**: Business logic is encapsulated in services
- **Middleware Pattern**: Cross-cutting concerns (auth, CORS, logging)
- **DTO Pattern**: Separate request/response models from domain models

### Security Architecture

```plantuml
@startuml SecurityFlow
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Sequence.puml

title Security Flow - JWT Authentication

participant Client
participant "Fiber Router" as Router
participant "Auth Handler" as Auth
participant "User Service" as UserSvc
participant "JWT Middleware" as JWT
participant "Protected Handler" as Handler
participant Database

== Registration/Login ==
Client -> Router: POST /register or /login
Router -> Auth: Route request
Auth -> UserSvc: Validate credentials
UserSvc -> Database: Check/create user
Database --> UserSvc: User data
UserSvc --> Auth: User object
Auth -> Auth: Generate JWT token
Auth --> Client: JWT token + user data

== Protected Endpoint Access ==
Client -> Router: GET /me (with Authorization header)
Router -> JWT: Validate token
JWT -> JWT: Parse and verify JWT
JWT -> JWT: Extract user context
JWT -> Handler: Forward with user context
Handler -> UserSvc: Get user data
UserSvc -> Database: Query user
Database --> UserSvc: User data
UserSvc --> Handler: User object
Handler --> Client: User profile

@enduml
```

### Data Flow Architecture

```plantuml
@startuml DataFlow
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title Data Flow - Point Transfer Transaction

Container_Boundary(request, "Request Flow") {
    Component(client, "Client Request", "JSON", "Transfer request with to_lbk_code and amount")
    Component(validation, "Input Validation", "Handler", "Validate request format and required fields")
    Component(auth, "Authentication", "Middleware", "Extract user ID from JWT token")
}

Container_Boundary(business, "Business Logic") {
    Component(transfer, "Transfer Service", "Business Logic", "Orchestrate point transfer operation")
    Component(transaction, "Database Transaction", "GORM", "Ensure atomic operations")
    Component(balance, "Balance Validation", "Service", "Check sufficient funds")
}

Container_Boundary(persistence, "Data Persistence") {
    Component(userUpdate, "User Balance Update", "Database", "Update sender and receiver balances")
    Component(transferLog, "Transfer Record", "Database", "Create audit trail record")
    Component(commit, "Transaction Commit", "Database", "Finalize all changes")
}

Rel(client, validation, "1. Submit transfer", "")
Rel(validation, auth, "2. Validate format", "")
Rel(auth, transfer, "3. Authenticate user", "")
Rel(transfer, transaction, "4. Start transaction", "")
Rel(transaction, balance, "5. Check balance", "")
Rel(balance, userUpdate, "6. Update balances", "")
Rel(userUpdate, transferLog, "7. Log transfer", "")
Rel(transferLog, commit, "8. Commit transaction", "")

@enduml
```

## Quality Attributes

### Maintainability
- **Clean Architecture**: Clear separation of concerns
- **Dependency Injection**: Loose coupling between components
- **Single Responsibility**: Each component has one reason to change
- **Documentation**: Comprehensive Swagger API docs and architecture docs

### Security
- **JWT Authentication**: Stateless token-based authentication
- **Password Hashing**: bcrypt for secure password storage
- **Input Validation**: Request validation at handler level
- **SQL Injection Protection**: GORM ORM prevents SQL injection

### Performance
- **Lightweight Framework**: Fiber provides high-performance HTTP handling
- **Database Indexing**: Optimized queries with proper indexes
- **Connection Pooling**: GORM manages database connections efficiently
- **Stateless Design**: No server-side session storage

### Scalability
- **Stateless Architecture**: Horizontal scaling capability
- **Database Independence**: Easy to switch databases
- **Microservice Ready**: Clean boundaries for service decomposition
- **Configuration Management**: Environment-based configuration

### Reliability
- **Transaction Safety**: Atomic operations for point transfers
- **Error Handling**: Comprehensive error responses
- **Audit Trail**: Complete transaction history
- **Health Checks**: Monitoring endpoints for system health

## Deployment Architecture

```plantuml
@startuml DeploymentDiagram
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Deployment.puml

title Deployment Diagram - Production Environment

Deployment_Node(loadBalancer, "Load Balancer", "NGINX/HAProxy") {
    Container(lb, "Load Balancer", "Reverse Proxy", "Traffic distribution and SSL termination")
}

Deployment_Node(appServers, "Application Servers", "Docker/Kubernetes") {
    Deployment_Node(server1, "App Server 1", "Docker Container") {
        Container(app1, "Fiber API", "Go Binary", "API server instance 1")
    }
    Deployment_Node(server2, "App Server 2", "Docker Container") {
        Container(app2, "Fiber API", "Go Binary", "API server instance 2")
    }
}

Deployment_Node(database, "Database Server", "SQLite/PostgreSQL") {
    ContainerDb(db, "Database", "SQLite", "User and transfer data")
}

Deployment_Node(monitoring, "Monitoring", "Observability Stack") {
    Container(metrics, "Metrics", "Prometheus", "Application metrics")
    Container(logs, "Logs", "ELK Stack", "Centralized logging")
}

Rel(lb, app1, "Route requests", "HTTPS")
Rel(lb, app2, "Route requests", "HTTPS")
Rel(app1, db, "Database queries", "SQL")
Rel(app2, db, "Database queries", "SQL")
Rel(app1, metrics, "Metrics export", "HTTP")
Rel(app2, metrics, "Metrics export", "HTTP")
Rel(app1, logs, "Application logs", "JSON")
Rel(app2, logs, "Application logs", "JSON")

@enduml
```

## Future Architecture Considerations

### Microservices Evolution
- **User Service**: Extract user management to separate service
- **Point Service**: Dedicated service for point operations
- **Notification Service**: Handle transfer notifications
- **API Gateway**: Centralized routing and authentication

### Advanced Features
- **Event Sourcing**: Track all state changes as events
- **CQRS**: Separate read and write models
- **Caching Layer**: Redis for performance optimization
- **Message Queue**: Async processing for transfers

### Observability
- **Distributed Tracing**: OpenTelemetry integration
- **Structured Logging**: Standardized log format
- **Metrics Dashboard**: Grafana visualization
- **Alerting**: Proactive monitoring and alerts
