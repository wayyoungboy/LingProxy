# LingProxy - AI API Gateway

LingProxy is a high-performance AI API gateway designed for managing and proxying API calls to various AI service providers. It offers OpenAI compatible interfaces, load balancing, circuit breaking, and more.

## Features

### 🚀 Core Features
- **Unified API Interface**: Supports OpenAI compatible API, seamlessly integrates with various AI services
- **Intelligent Load Balancing**: Round-robin load balancing strategy, automatically distributes requests to multiple resources
- **Circuit Breaking**: Automatically detects service failures and triggers circuit breaking to prevent cascading failures
- **Request Logging**: Complete request chain tracing and logging

### 🔐 Security & Authentication
- **Flexible Authentication**: Global authentication toggle, configurable authentication requirement
- **Admin Login**: Username/password login with password hash storage
- **Token Management**: Request-side token management with policy association and API key authentication
- **CORS Support**: Flexible cross-origin resource sharing configuration
- **Secure Storage**: Encrypted storage for API keys and passwords

### 📊 Management Features
- **Admin Management**: Single admin mode with password and API key management
- **Token Management**: Create and manage request-side tokens with policy binding
- **Policy Management**: Built-in routing policy templates (random, round-robin, weighted, model-match, regex-match, priority, failover), supports custom policy instances
- **LLM Resource Management**: Supports configuration of multiple AI service providers (OpenAI, Zai, Anthropic, Google, Azure, etc.), supports model categories (chat, image, embedding, rerank, audio, video)
- **Model Management**: Flexible model configuration, supports pricing, usage limits and other parameters
- **System Settings**: Dynamic configuration management including basic settings, cache, rate limiting, security, logging, load balancing, circuit breaker configurations
- **System Monitoring**: Real-time system information (CPU, memory, uptime, etc.)

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
# ⚠️ IMPORTANT: Change the admin password in config.yaml before starting!
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

The frontend will be available at `http://localhost:3000`

### Docker Run

```bash
# Build image
docker build -t lingproxy .

# Run container
docker run -p 8080:8080 -v $(pwd)/configs:/app/configs lingproxy
```

## API Usage Guide

### 1. Admin Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "YOUR_PASSWORD"
  }'
```

Response example:
```json
{
  "token": "your_jwt_token_here",
  "user": {
    "id": "...",
    "username": "admin",
    "api_key": "..."
  }
}
```

### 2. Create Request-side Token
```bash
curl -X POST http://localhost:8080/api/v1/tokens \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My API Token",
    "status": "active"
  }'
```

Response example:
```json
{
  "data": {
    "id": "...",
    "name": "My API Token",
    "token": "ling-xxxxxxxxxxxxx",
    "status": "active"
  }
}
```

### 3. Create Routing Policy (Optional)
```bash
curl -X POST http://localhost:8080/api/v1/policies \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Random Policy",
    "template_id": "random_template_id",
    "type": "random",
    "parameters": "{\"filter_by_status\": true}",
    "enabled": true
  }'
```

### 4. Bind Policy to Token (Optional)
```bash
curl -X PUT http://localhost:8080/api/v1/tokens/TOKEN_ID/policy \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "policy_id": "POLICY_ID"
  }'
```

### 5. Proxy AI Request
```bash
curl -X POST http://localhost:8080/llm/v1/chat/completions \
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

#### Authentication & Admin
- `POST /api/v1/auth/login` - Admin login (username/password)
- `GET /api/v1/admin/info` - Get admin information
- `PUT /api/v1/admin/api-key` - Reset admin API key

#### Token Management
- `GET /api/v1/tokens` - Get token list
- `GET /api/v1/tokens/:id` - Get token details
- `POST /api/v1/tokens` - Create token
- `PUT /api/v1/tokens/:id` - Update token
- `DELETE /api/v1/tokens/:id` - Delete token
- `POST /api/v1/tokens/:id/reset` - Reset token
- `PUT /api/v1/tokens/:id/policy` - Bind policy to token
- `DELETE /api/v1/tokens/:id/policy` - Remove policy binding from token

#### Policy Management
- `GET /api/v1/policy-templates` - Get policy template list
- `GET /api/v1/policy-templates/:id` - Get policy template details
- `GET /api/v1/policies` - Get policy list
- `GET /api/v1/policies/:id` - Get policy details
- `POST /api/v1/policies` - Create policy
- `PUT /api/v1/policies/:id` - Update policy
- `DELETE /api/v1/policies/:id` - Delete policy

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
- `GET /api/v1/models/types` - Get model type list
- `GET /api/v1/models/:id/pricing` - Get model pricing information
- `GET /api/v1/llm-resources/:id/models` - Get models under specified LLM resource

#### Request Logging
- `GET /api/v1/requests` - Get request log list
- `GET /api/v1/requests/:id` - Get request details
- `POST /api/v1/requests` - Create request record

#### System Settings & Monitoring
- `GET /api/v1/settings` - Get system settings
- `PUT /api/v1/settings` - Update system settings
- `GET /api/v1/system/info` - Get system information (CPU, memory, uptime, etc.)

#### Statistics
- `GET /api/v1/stats/system` - Get system statistics
- `GET /api/v1/stats/llm-resources/:id` - Get LLM resource statistics
- `GET /api/v1/stats/users/:id` - Get user statistics

#### OpenAI Compatible API
- `GET /llm/v1/models` - List all available models
- `GET /llm/v1/models/:model` - Get model information
- `POST /llm/v1/chat/completions` - Create chat completion
- `POST /llm/v1/completions` - Create text completion

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
  auth:
    enabled: true  # Whether to enable authentication, when false all APIs (except login) don't require authentication
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

#### Admin Configuration
```yaml
admin:
  username: "admin"
  # ⚠️ Set a strong password! Recommended to set immediately after first startup
  # password: "YOUR_STRONG_PASSWORD_HERE"
  password: ""  # Leave empty to skip password setup
  api_key: ""  # Leave empty to auto-generate, check logs after first startup
  auto_create: true
```

#### Logging Configuration
```yaml
log:
  level: "info"  # debug, info, warn, error, fatal
  format: "json"  # text, json
  output: "stdout"
```

#### Load Balancer Configuration
```yaml
load_balancer:
  default_strategy: "round_robin"  # Default load balancing strategy
```

#### Provider Configuration
```yaml
provider:
  timeout: "30s"  # Request timeout
  max_retries: 3   # Maximum retry count
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
└── docker-compose.yml     # Docker configuration
```

### Data Models

The system adopts a streamlined storage model design with core models including:

```go
// User user model - admin user
type User struct {
    ID           string     // User unique identifier
    Username     string     // Username
    PasswordHash string     // Password hash
    APIKey       string     // API key
    Role         string     // Role (admin)
    Status       string     // Status (active, inactive, suspended)
    LastLoginAt  *time.Time // Last login time
    CreatedAt    time.Time  // Created at
    UpdatedAt    time.Time  // Updated at
}

// Token token model - request-side token management
type Token struct {
    ID         string     // Token unique identifier
    Name       string     // Token name/description
    Token      string     // Token value (API Key, prefixed with "ling-")
    Prefix     string     // Token prefix (for display)
    Status     string     // Status (active/inactive)
    PolicyID   string     // Associated policy ID (optional)
    LastUsedAt *time.Time // Last used time
    ExpiresAt  *time.Time // Expiration time (optional)
    CreatedAt  time.Time  // Created at
    UpdatedAt  time.Time  // Updated at
}

// PolicyTemplate policy template model - built-in policy templates
type PolicyTemplate struct {
    ID                string    // Template unique identifier
    Name              string    // Template name
    Type              string    // Type (random, round_robin, weighted, model_match, regex_match, priority, failover)
    Description       string    // Description
    ParametersSchema  string    // Parameter JSON Schema
    DefaultParameters string    // Default parameters JSON
    Builtin           bool      // Whether built-in
    CreatedAt         time.Time // Created at
    UpdatedAt         time.Time // Updated at
}

// Policy policy instance model - routing policy configuration
type Policy struct {
    ID         string    // Policy unique identifier
    Name       string    // Policy name
    TemplateID string    // Associated template ID
    Type       string    // Type
    Parameters string    // Parameters JSON
    Enabled    bool      // Whether enabled
    CreatedAt  time.Time // Created at
    UpdatedAt  time.Time // Updated at
}

// LLMResource LLM resource model - AI service provider configuration
type LLMResource struct {
    ID        string    // Resource unique identifier
    Name      string    // Resource name
    Type      string    // Model category (chat, image, embedding, rerank, audio, video)
    Provider  string    // Service provider (openai, zai, anthropic, google, azure, custom, etc.)
    Model     string    // Default model
    BaseURL   string    // API base URL
    APIKey    string    // API key
    Status    string    // Status (active/inactive)
    CreatedAt time.Time // Created at
    UpdatedAt time.Time // Updated at
}

// Model model configuration - AI model management
type Model struct {
    ID            string    // Model unique identifier
    Name          string    // Model name
    LLMResourceID string    // Associated LLM resource
    ModelID       string    // Provider's internal model identifier
    Type          string    // Model type (chat, completion, embedding, image)
    Category      string    // Model category (gpt, claude, gemini, llama, etc.)
    Version       string    // Model version
    Description   string    // Description
    Capabilities  string    // Model capabilities (JSON string)
    Pricing       string    // Pricing information (JSON string)
    Limits        string    // Usage limits (JSON string)
    Parameters    string    // Default parameters (JSON string)
    Features      string    // Features (JSON string)
    Status        string    // Status (active, inactive, deprecated)
    Metadata      string    // Extended metadata (JSON string)
    CreatedAt     time.Time // Created at
    UpdatedAt     time.Time // Updated at
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
    CreatedAt time.Time // Created at
}
```

### Storage Layer Design

The storage layer adopts a clean interface design, supporting both memory storage and GORM storage implementations:

```go
type Storage interface {
    // User management
    CreateUser(user *User) error
    GetUser(id string) (*User, error)
    GetUserByUsername(username string) (*User, error)
    GetUserByAPIKey(apiKey string) (*User, error)
    UpdateUser(user *User) error
    DeleteUser(id string) error
    ListUsers() ([]*User, error)

    // Token management
    CreateToken(token *Token) error
    GetToken(id string) (*Token, error)
    GetTokenByToken(token string) (*Token, error)
    UpdateToken(token *Token) error
    DeleteToken(id string) error
    ListTokens() ([]*Token, error)

    // Policy template management
    CreatePolicyTemplate(template *PolicyTemplate) error
    GetPolicyTemplate(id string) (*PolicyTemplate, error)
    GetPolicyTemplateByType(type string) (*PolicyTemplate, error)
    UpdatePolicyTemplate(template *PolicyTemplate) error
    DeletePolicyTemplate(id string) error
    ListPolicyTemplates() ([]*PolicyTemplate, error)

    // Policy management
    CreatePolicy(policy *Policy) error
    GetPolicy(id string) (*Policy, error)
    UpdatePolicy(policy *Policy) error
    DeletePolicy(id string) error
    ListPolicies() ([]*Policy, error)

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