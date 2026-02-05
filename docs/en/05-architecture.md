# Architecture Guide

## System Architecture

LingProxy follows a modern microservices architecture with clear separation of concerns:

```
┌─────────────┐
│   Frontend  │  Vue 3 + Element Plus
│  (Port 3000)│
└──────┬──────┘
       │ HTTP/REST API
┌──────▼──────────────────┐
│      Backend API        │  Go + Gin
│     (Port 8080)         │
├──────────────────────────┤
│  ┌────────────────────┐ │
│  │   HTTP Handlers    │ │  Request processing
│  └──────────┬─────────┘ │
│  ┌──────────▼─────────┐ │
│  │   Middleware       │ │  Auth, CORS, Logging
│  └──────────┬─────────┘ │
│  ┌──────────▼─────────┐ │
│  │   Services         │ │  Business logic
│  └──────────┬─────────┘ │
│  ┌──────────▼─────────┐ │
│  │   Storage Layer    │ │  Data persistence
│  └────────────────────┘ │
└──────────────────────────┘
       │
┌──────▼──────┐
│   Database  │  SQLite/MySQL/PostgreSQL
└─────────────┘
```

## Frontend Architecture

### Technology Stack
- **Framework**: Vue 3 (Composition API)
- **UI Component Library**: Element Plus
- **Build Tool**: Vite
- **Internationalization**: vue-i18n
- **Routing**: Vue Router
- **HTTP Client**: Axios

### Directory Structure
```
frontend/
├── src/
│   ├── api/              # API client
│   ├── assets/           # Static assets
│   ├── components/       # Vue components
│   │   ├── Layout.vue    # Layout component (includes language switcher)
│   │   └── Sidebar.vue   # Sidebar component
│   ├── config/           # Configuration files
│   │   └── menu.js       # Menu configuration
│   ├── locales/          # Internationalization language packs
│   │   ├── zh/           # Chinese language pack
│   │   ├── en/           # English language pack
│   │   └── index.js      # i18n configuration
│   ├── router/           # Route configuration
│   ├── views/            # Page views
│   │   ├── Login.vue
│   │   ├── Dashboard.vue
│   │   ├── Tokens.vue
│   │   ├── LLMResources.vue
│   │   ├── LLMResourceUsage.vue
│   │   ├── Requests.vue
│   │   ├── Policies.vue
│   │   ├── Settings.vue
│   │   ├── Logs.vue
│   │   ├── Models.vue
│   │   ├── Users.vue
│   │   └── Endpoints.vue
│   ├── App.vue           # Root component
│   └── main.js           # Entry point
├── package.json
└── vite.config.js
```

### Internationalization Support
- **Language Packs**: Complete Chinese and English language packs
- **Language Switching**: Support for runtime language switching, settings saved in localStorage
- **Element Plus Integration**: Element Plus component language automatically follows system language settings
- **Coverage**: All user interface text, error messages, and form validation messages are internationalized

### Core Feature Modules
- **Authentication**: Login page, JWT Token management
- **Dashboard**: System overview and statistics
- **Resource Management**: LLM resources, models, endpoints management
- **Policy Management**: Routing policy configuration and management
- **Request Management**: Request log viewing and export
- **Usage Statistics**: Detailed statistics grouped by resource
- **System Settings**: Dynamic configuration management
- **Log Management**: System log viewing and management

## Backend Architecture

### Directory Structure

```
backend/
├── cmd/
│   └── main.go              # Application entry point
├── configs/
│   └── config.yaml.example  # Configuration template
├── internal/
│   ├── cache/               # Cache implementation
│   ├── client/              # AI service clients
│   │   ├── embedding/       # Embedding client
│   │   └── openai/          # OpenAI client
│   ├── config/              # Configuration management
│   ├── handler/             # HTTP handlers
│   ├── middleware/          # HTTP middleware
│   ├── pkg/                 # Internal utilities
│   │   ├── balancer/        # Load balancer
│   │   ├── logger/         # Logging utilities
│   │   ├── monitor/        # Monitoring utilities
│   │   └── password/       # Password utilities
│   ├── router/             # Route configuration
│   ├── service/            # Business logic services
│   └── storage/            # Storage layer
│       ├── models.go       # Data models
│       ├── storage.go      # Storage interface
│       ├── storage_facade.go # Storage facade
│       ├── memory_storage.go # Memory implementation
│       └── gorm_storage.go   # GORM implementation
└── swagger/                # API documentation
```

### Layer Architecture

#### 1. Handler Layer
- **Purpose**: HTTP request/response handling
- **Responsibilities**:
  - Parse HTTP requests
  - Validate input data
  - Call service layer
  - Format HTTP responses
- **Files**: `internal/handler/*.go`

#### 2. Service Layer
- **Purpose**: Business logic implementation
- **Responsibilities**:
  - Implement business rules
  - Coordinate between handlers and storage
  - Handle complex operations
- **Files**: `internal/service/*.go`

#### 3. Storage Layer
- **Purpose**: Data persistence abstraction
- **Responsibilities**:
  - Define storage interface
  - Implement storage backends (memory, GORM)
  - Handle data operations
- **Files**: `internal/storage/*.go`

#### 4. Middleware Layer
- **Purpose**: Cross-cutting concerns
- **Responsibilities**:
  - Authentication
  - CORS handling
  - Request logging
  - Rate limiting
- **Files**: `internal/middleware/*.go`

## Data Models

### Core Models

#### User
```go
type User struct {
    ID           string
    Username     string
    PasswordHash string
    APIKey       string
    Role         string    // admin
    Status       string    // active, inactive, suspended
    LastLoginAt  *time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

#### LLMResource
```go
type LLMResource struct {
    ID        string
    Name      string
    Type      string    // chat, image, embedding, etc.
    Driver    string    // openai (currently only openai)
    Model     string
    BaseURL   string
    APIKey    string
    Status    string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### Token
```go
type Token struct {
    ID         string
    Name       string
    Token      string
    Prefix     string
    Status     string
    PolicyID   string
    LastUsedAt *time.Time
    ExpiresAt  *time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

#### Policy
```go
type Policy struct {
    ID         string
    Name       string
    TemplateID string
    Type       string    // round_robin, random, etc.
    Parameters string    // JSON
    Enabled    bool
    Builtin    bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

## Request Flow

### OpenAI-Compatible API Request

```
1. Client Request
   ↓
2. Authentication Middleware
   - Validate API key/token
   ↓
3. Request Logger Middleware
   - Log request details
   ↓
4. OpenAI Handler
   - Parse request
   - Select resource (via policy)
   ↓
5. Policy Executor
   - Execute routing policy
   - Select LLM resource
   ↓
6. Client Manager
   - Create/Get client
   - Forward request to AI service
   ↓
7. Response Processing
   - Format response
   - Log response
   ↓
8. Return to Client
```

### Management API Request

```
1. Client Request
   ↓
2. Authentication Middleware
   - Validate admin credentials
   ↓
3. Handler
   - Parse request
   - Validate input
   ↓
4. Service Layer
   - Execute business logic
   - Update storage
   ↓
5. Storage Layer
   - Persist changes
   ↓
6. Response
   - Return result
```

## Storage Backends

### Memory Storage
- **Use Case**: Development and testing
- **Characteristics**: Fast, ephemeral, no persistence
- **Implementation**: `memory_storage.go`

### GORM Storage
- **Use Case**: Production environments
- **Characteristics**: Persistent, supports SQLite/MySQL/PostgreSQL
- **Implementation**: `gorm_storage.go`

## Security Architecture

### Authentication Flow

```
1. Client sends request with API key/token
   ↓
2. Auth Middleware extracts credentials
   ↓
3. Validate credentials:
   - Check token in TokenService
   - Or check API key in User storage
   ↓
4. Set user context
   ↓
5. Continue to handler
```

### Password Security
- Passwords are hashed using bcrypt
- Never stored in plain text
- Password verification uses constant-time comparison

## Load Balancing

### Supported Strategies
- **Round Robin**: Distribute requests sequentially
- **Random**: Random selection
- **Weighted**: Weighted distribution
- **Model Match**: Match by model name
- **Regex Match**: Match by pattern
- **Priority**: Priority-based selection
- **Failover**: Automatic failover

## Configuration Management

### Configuration Sources (Priority Order)
1. Environment variables (`LINGPROXY_*`)
2. Configuration file (`config.yaml`)
3. Default values (in code)

### Configuration Structure
- Application settings
- Storage configuration
- Log configuration
- Security settings

## Error Handling

### Error Types
- **Validation Errors**: 400 Bad Request
- **Authentication Errors**: 401 Unauthorized
- **Not Found Errors**: 404 Not Found
- **Server Errors**: 500 Internal Server Error

### Error Response Format
```json
{
  "error": "Error message"
}
```

## Logging

### Log Levels
- **Debug**: Detailed debugging information
- **Info**: General informational messages
- **Warn**: Warning messages
- **Error**: Error messages
- **Fatal**: Fatal errors

### Log Output
- Console (stdout)
- File (rotated logs)
- Both (recommended)

## Performance Considerations

### Caching
- In-memory caching for frequently accessed data
- Configurable TTL

### Connection Pooling
- HTTP client connection pooling
- Configurable pool size

### Database Optimization
- Indexed queries
- Efficient data models
- Connection pooling
