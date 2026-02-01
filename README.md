# LingProxy - AI API Gateway

<p align="center">
  <img src="web/static/images/logo/lingproxy-logo.svg" alt="LingProxy Logo">
</p>


LingProxy is a high-performance AI API gateway designed for managing and proxying API calls to various AI service providers. It offers OpenAI compatible interfaces, load balancing, circuit breaking, and more.

## Features

### 🚀 Core Features
- **Unified API Interface**: Supports OpenAI compatible API, seamlessly integrates with various AI services
- **Intelligent Load Balancing**: Round-robin load balancing strategy, automatically distributes requests to multiple resources
- **Circuit Breaking**: Automatically detects service failures and triggers circuit breaking to prevent cascading failures
- **Request Logging**: Complete request chain tracing and logging

### 🔐 Security & Authentication
- **API Key Authentication**: API key-based user authentication mechanism
- **CORS Support**: Flexible cross-origin resource sharing configuration
- **Secure Storage**: Encrypted API key storage

### 📊 Management Features
- **User Management**: Multi-user support, API key management
- **LLM Resource Management**: Supports configuration of multiple AI service providers (OpenAI, Anthropic, Google, etc.)
- **Model Management**: Flexible model configuration, supports pricing, usage limits and other parameters
- **Endpoint Management**: Custom API endpoint routing configuration

### 🏗️ Architecture Design
- **Frontend-Backend Separation**: Modern architecture with Vue 3 + Element Plus frontend and Go backend
- **Simplified Models**: Removed redundant features, core code is concise and efficient
- **Dual Storage**: Supports memory storage (development and debugging) and SQLite storage (production environment)
- **Modular Design**: Clear hierarchical structure, easy to extend and maintain
- **RESTful API**: Complete REST API interface, easy to integrate

## Quick Start

### Requirements
- **Backend**: Go 1.21 or higher, SQLite (for data storage)
- **Frontend**: Node.js 18+, npm or yarn

### Installation & Running

#### Backend Setup

1. **Clone the Project**
```bash
git clone https://github.com/wayyoungboy/lingproxy.git
cd lingproxy
```

2. **Install Go Dependencies**
```bash
go mod tidy
```

3. **Configuration File**
Copy and edit the configuration file:
```bash
cp configs/config.yaml.example configs/config.yaml
# Edit configs/config.yaml to configure as needed
```

4. **Build and Run the Backend**
```bash
go run cmd/main.go
```

The backend service will start at `http://localhost:8080`

#### Frontend Setup

1. **Install Node.js Dependencies**
```bash
cd frontend
npm install
```

2. **Run the Frontend Development Server**
```bash
npm run dev
```

The frontend will be available at `http://localhost:3002`

### Docker Run

```bash
# Build image
docker build -t lingproxy .

# Run container
docker run -p 8080:8080 -v $(pwd)/configs:/app/configs lingproxy
```

## API Usage Guide

### 1. User Registration
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 2. Get API Key
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 3. Proxy AI Request
```bash
curl -X POST http://localhost:8080/api/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ]
  }'
```

### API Endpoint Reference

#### User Management
- `GET /api/v1/users` - Get user list
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user details
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

#### LLM Resource Management
- `GET /api/v1/llm-resources` - Get LLM resource list
- `POST /api/v1/llm-resources` - Create LLM resource
- `GET /api/v1/llm-resources/:id` - Get LLM resource details
- `PUT /api/v1/llm-resources/:id` - Update LLM resource
- `DELETE /api/v1/llm-resources/:id` - Delete LLM resource

#### Model Management
- `GET /api/v1/models` - Get model list
- `POST /api/v1/models` - Create model
- `GET /api/v1/models/:id` - Get model details
- `PUT /api/v1/models/:id` - Update model
- `DELETE /api/v1/models/:id` - Delete model

#### Endpoint Management
- `GET /api/v1/endpoints` - Get endpoint list
- `POST /api/v1/endpoints` - Create endpoint
- `GET /api/v1/endpoints/:id` - Get endpoint details
- `PUT /api/v1/endpoints/:id` - Update endpoint
- `DELETE /api/v1/endpoints/:id` - Delete endpoint

#### Request Logging
- `GET /api/v1/requests` - Get request log list
- `GET /api/v1/requests/:id` - Get request details

## Configuration

### Main Configuration Items

#### Application Configuration
```yaml
app:
  name: "LingProxy"
  version: "1.0.0"
  environment: "development"  # development, staging, production
  port: 8080
  host: "0.0.0.0"
```

#### Storage Configuration
```yaml
storage:
  type: "gorm"
  gorm:
    driver: "sqlite"
    dsn: "lingproxy.db"
```

#### Security Configuration
```yaml
security:
  cors:
    enabled: true
    allow_origins:
      - "*"
    allow_methods:
      - "GET"
      - "POST"
      - "PUT"
      - "DELETE"
      - "OPTIONS"
    allow_headers:
      - "*"
```

## Monitoring & Operations

### Logging System

#### Log Levels
- **DEBUG**: Detailed debug information, only for development
- **INFO**: General information about system operations
- **WARN**: Warning messages that need attention
- **ERROR**: Error messages that require immediate action
- **FATAL**: Critical errors that cause system shutdown

#### Log Configuration
```yaml
log:
  level: "info"  # debug, info, warn, error, fatal
  format: "json"  # text, json
  output: "stdout"
```

#### Log Viewing
```bash
# View real-time logs
# Logs are output to stdout by default
```

## Development Guide

### Project Structure
```
lingproxy/
├── cmd/                    # Application entry
├── configs/               # Configuration files
├── docs/                  # API documentation
├── frontend/              # Frontend application
│   ├── public/             # Public assets
│   ├── src/                # Source code
│   │   ├── api/            # API client
│   │   ├── assets/         # Static assets
│   │   ├── components/     # Vue components
│   │   ├── router/         # Vue router
│   │   ├── views/           # Vue views
│   │   ├── App.vue         # Root component
│   │   └── main.js         # Entry point
│   ├── package.json        # npm configuration
│   └── vite.config.js      # Vite configuration
├── internal/              # Internal packages
│   ├── cache/             # Caching implementation
│   ├── client/            # AI service clients
│   │   ├── embedding/     # Embedding clients
│   │   └── openai/        # OpenAI clients
│   ├── config/            # Configuration management
│   ├── handler/           # HTTP handlers
│   ├── middleware/        # HTTP middleware
│   ├── pkg/               # Internal packages
│   │   ├── balancer/      # Load balancing
│   │   └── circuitbreaker/ # Circuit breaker
│   ├── router/            # Routing
│   ├── service/           # Business logic
│   └── storage/           # Storage implementation
├── pkg/                   # Public packages
│   └── logger/            # Logging
├── scripts/               # Script files
├── web/                   # Legacy web interface (deprecated)
└── docker-compose.yml     # Docker configuration
```

### Data Models

The system adopts a streamlined storage model design with core models including:

```go
// User user model - manages API users
type User struct {
    ID        string    // User unique identifier
    Username  string    // Username
    Email     string    // Email
    APIKey    string    // API key
    Status    string    // Status
}

// LLMResource LLM resource model - AI service provider configuration
type LLMResource struct {
    ID        string    // Resource unique identifier
    Name      string    // Resource name
    Type      string    // Type (openai, anthropic, google, etc.)
    Model     string    // Default model
    BaseURL   string    // API base URL
    APIKey    string    // API key
    Status    string    // Status
}

// Model model configuration - AI model management
type Model struct {
    ID             string         // Model unique identifier
    Name           string         // Model name
    LLMResourceID  string         // Associated LLM resource
    ModelID        string         // Provider's internal model identifier
    Type           string         // Model type (chat, completion, embedding, image)
    Category       string         // Model category
    Pricing        ModelPricing   // Pricing information
    Limits         ModelLimits    // Usage limits
    Parameters     ModelParameters // Default parameters
    Status         string         // Status
}

// Endpoint endpoint model - API route configuration
type Endpoint struct {
    ID            string    // Endpoint unique identifier
    LLMResourceID string    // Associated LLM resource
    Path          string    // API path
    Method        string    // HTTP method
    Status        string    // Status
}

// Request request model - request logging
type Request struct {
    ID        string    // Request unique identifier
    UserID    string    // User ID
    Endpoint  string    // Request endpoint
    Method    string    // HTTP method
    Status    string    // Status
    Duration  int64     // Duration (milliseconds)
    Tokens    int       // Consumed tokens
}
```

### Storage Layer Design

The storage layer adopts a clean interface design, supporting both memory storage and GORM storage implementations:

```go
type Storage interface {
    // User management
    CreateUser(user *User) error
    GetUser(id string) (*User, error)
    GetUserByAPIKey(apiKey string) (*User, error)
    UpdateUser(user *User) error
    DeleteUser(id string) error
    ListUsers() ([]*User, error)

    // LLMResource management
    CreateLLMResource(resource *LLMResource) error
    GetLLMResource(id string) (*LLMResource, error)
    UpdateLLMResource(resource *LLMResource) error
    DeleteLLMResource(id string) error
    ListLLMResources() ([]*LLMResource, error)

    // Model management
    CreateModel(model *Model) error
    GetModel(id string) (*Model, error)
    UpdateModel(model *Model) error
    DeleteModel(id string) error
    ListModels() ([]*Model, error)
    ListModelsByLLMResource(llmResourceID string) ([]*Model, error)

    // Endpoint management
    CreateEndpoint(endpoint *Endpoint) error
    GetEndpoint(id string) (*Endpoint, error)
    UpdateEndpoint(endpoint *Endpoint) error
    DeleteEndpoint(id string) error
    ListEndpoints() ([]*Endpoint, error)

    // Request logging
    CreateRequest(request *Request) error
    GetRequest(id string) (*Request, error)
    ListRequests(limit int) ([]*Request, error)
}
```

### Adding a New AI Provider

1. **Create LLM Resource Configuration**
Add new provider type handling logic in `internal/handler/provider.go`

2. **Implement Load Balancing Strategy**
Implement a new load balancing algorithm in `internal/pkg/balancer/`

3. **Update Model Support**
Extend the Model struct in `internal/storage/model.go` to support new model features

### Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/pkg/balancer

# Run tests with coverage
go test -cover ./...
```

## Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support & Contact

- **Issues**: [GitHub Issues](https://github.com/wayyoungboy/lingproxy/issues)
- **Discussions**: [GitHub Discussions](https://github.com/wayyoungboy/lingproxy/discussions)
- **Email**: support@lingproxy.com

## Changelog

### v1.2.0 (2026-02-02)
- **Frontend-Backend Separation**: Implemented modern Vue 3 + Element Plus frontend
- **New Frontend Interface**: Complete rewrite with Vue 3 Composition API and Script Setup
- **Enhanced UI**: Responsive design with Element Plus components
- **Improved API Integration**: Axios-based API client with proper error handling
- **Backend API Updates**: Added missing endpoint management APIs
- **Web Interface Removal**: Removed legacy web interface routes
- **Documentation Updates**: Added frontend development guide and architecture documentation

### v1.1.0 (2026-02-01)
- **Architecture Optimization**: Simplified core code, improved code quality and maintainability
- **Model Simplification**: Removed unused ModelEndpoint and ModelVersion structs
- **Monitoring Module Optimization**: Simplified to lightweight quota manager
- **Storage Layer Refactoring**: Optimized storage interface, removed redundant methods
- **Dependency Fixes**: Fixed embedding client dependency issues
- **Documentation Updates**: Improved development guide and data model documentation

### v1.0.0 (2026-01-30)
- Initial release
- Support for OpenAI compatible API
- Implementation of round-robin load balancing and circuit breaking
- Added user management and LLM resource management
- Provided complete REST API interface
- Implemented SQLite-based data storage
- Added logging system with multiple log levels
- Created web-based admin interface

## Language

- [English](README.md) (current)
- [中文](README_zh.md)