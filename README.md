<div align="center">

<img src="./assets/lingproxy-logo-full.svg" alt="LingProxy Logo" width="300">

# LingProxy - AI API Gateway

LingProxy is a high-performance AI API gateway designed for managing and proxying API calls to various AI service providers. It offers OpenAI compatible interfaces, load balancing, circuit breaking, and more.

</div>

## Features

### 🚀 Core Features
- **Unified API Interface**: Supports OpenAI compatible API, seamlessly integrates with various AI services
- **Streaming Support**: Full support for Server-Sent Events (SSE) streaming responses for chat completions
- **Intelligent Load Balancing**: Round-robin load balancing strategy, automatically distributes requests to multiple resources
- **Automatic Retry**: Configurable automatic retry mechanism for failed requests with exponential backoff, supports retry for network errors, timeouts, and 5xx server errors
- **Circuit Breaking**: Automatically detects service failures and triggers circuit breaking to prevent cascading failures
- **Request Logging**: Complete request chain tracing and logging

### 🔐 Security & Authentication
- **Flexible Authentication**: Global authentication toggle, configurable authentication requirement
- **Admin Login**: Username/password login with password hash storage
- **API Key Management**: Request-side API key management with policy association and API key authentication
- **CORS Support**: Flexible cross-origin resource sharing configuration
- **Secure Storage**: Encrypted storage for API keys and passwords

### 📊 Management Features
- **Admin Dashboard**: Modern web-based management interface built with Vue 3 + Element Plus
- **Internationalization (i18n)**: Full support for Chinese and English language switching in the frontend interface
- **Admin Management**: Single admin mode with password and API key management
- **API Key Management**: Create and manage request-side API keys with policy binding, supports API key copying functionality
- **Policy Management**: Built-in routing policy templates (random, round-robin, weighted, model-match, regex-match, priority, failover), supports custom policy instances, supports LLM resource pool configuration for random selection policy
- **LLM Resource Management**: Supports configuration of AI service resources with driver-based architecture (currently supports OpenAI driver), supports model categories (chat, image, embedding, rerank, audio, video), supports batch import/export via Excel templates or JSON format, includes resource testing functionality to verify connectivity
- **Model Management**: Flexible model configuration, supports pricing, usage limits and other parameters
- **Request Management**: Complete request logging and tracking, supports request detail viewing and export
- **Usage Statistics**: Detailed usage statistics grouped by LLM resources, including token usage, request count, success rate, average tokens per request, and more, with support for time range and resource name filtering
- **System Settings**: Dynamic configuration management including basic settings, cache, rate limiting, security, logging, load balancing, and provider retry configurations
- **System Monitoring**: Real-time system information (CPU, memory, uptime, etc.)
- **Log Management**: View and manage system logs with filtering and search capabilities

### 🏗️ Architecture Design
- **Frontend-Backend Separation**: Modern architecture with Vue 3 + Element Plus frontend and Go backend
- **Internationalization**: Full i18n support with vue-i18n, supporting Chinese and English
- **Simplified Models**: Removed redundant features, core code is concise and efficient
- **Dual Storage**: Supports memory storage (development and debugging) and SQLite storage (production environment)
- **Modular Design**: Clear hierarchical structure, easy to extend and maintain
- **RESTful API**: Complete REST API interface, easy to integrate
- **Client Libraries**: Standard client implementations available in `clients/` directory (Python, JavaScript, Go)

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

#### Using Docker Compose (Recommended)

1. **Prepare Configuration**
```bash
# Copy configuration example
cp backend/configs/config.yaml.example backend/configs/config.yaml
# Edit backend/configs/config.yaml as needed
```

2. **Start Services**
```bash
# Build and start (from project root)
docker-compose -f docker/docker-compose.yml up -d

# View logs
docker-compose -f docker/docker-compose.yml logs -f

# Stop services
docker-compose -f docker/docker-compose.yml down
```

#### Using Docker Directly

```bash
# Build image (from project root)
docker build -f docker/Dockerfile -t lingproxy:latest .

# Run container
docker run -d \
  --name lingproxy \
  -p 8080:8080 \
  -v $(pwd)/backend/configs:/app/configs:ro \
  -v $(pwd)/logs:/app/logs \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/run:/app/run \
  -e GIN_MODE=release \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  lingproxy:latest
```

**Note**: 
- All Docker-related files are located in the `docker/` directory
- Make sure to create `backend/configs/config.yaml` before running the container
- See `docker/README.md` for detailed Docker deployment guide

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

### 2. Create Request-side API Key
```bash
curl -X POST http://localhost:8080/api/v1/api-keys \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My API Key",
    "status": "active"
  }'
```

Response example:
```json
{
  "data": {
    "id": "...",
    "name": "My API Key",
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

### 4. Bind Policy to API Key (Optional)
```bash
curl -X PUT http://localhost:8080/api/v1/api-keys/API_KEY_ID/policy \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "policy_id": "POLICY_ID"
  }'
```

**Note**: Replace `API_KEY_ID` with the actual API key ID from step 2.

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

#### API Key Management
- `GET /api/v1/api-keys` - Get API key list
- `GET /api/v1/api-keys/:id` - Get API key details
- `POST /api/v1/api-keys` - Create API key
- `PUT /api/v1/api-keys/:id` - Update API key
- `DELETE /api/v1/api-keys/:id` - Delete API key
- `POST /api/v1/api-keys/:id/reset` - Reset API key
- `PUT /api/v1/api-keys/:id/policy` - Bind policy to API key
- `DELETE /api/v1/api-keys/:id/policy` - Remove policy binding from API key

**Note**: The old `/api/v1/tokens` endpoints are still available for backward compatibility but are deprecated.

#### Policy Management
- `GET /api/v1/policy-templates` - Get policy template list
- `GET /api/v1/policy-templates/:id` - Get policy template details
- `GET /api/v1/policies` - Get policy list
- `GET /api/v1/policies/:id` - Get policy details
- `POST /api/v1/policies` - Create policy
- `PUT /api/v1/policies/:id` - Update policy
- `DELETE /api/v1/policies/:id` - Delete policy

#### LLM Resource Management
- `GET /api/v1/llm-resources` - Get LLM resource list (supports search filtering)
- `POST /api/v1/llm-resources` - Create LLM resource
- `GET /api/v1/llm-resources/:id` - Get LLM resource details
- `PUT /api/v1/llm-resources/:id` - Update LLM resource
- `DELETE /api/v1/llm-resources/:id` - Delete LLM resource
- `POST /api/v1/llm-resources/:id/test` - Test LLM resource connectivity
- `POST /api/v1/llm-resources/import` - Batch import LLM resources (Excel or JSON)
- `GET /api/v1/llm-resources/import/template` - Download Excel import template

**Batch Import**:
- Supports batch importing LLM resources via Excel files or JSON format
- Excel template includes fields: Name, Type, Driver, Model, BaseURL, APIKey, Status
- JSON import accepts an array of resource objects with the same fields
- Driver field currently only supports "openai", will be auto-set to "openai" if empty or invalid
- Import results return success/failure/duplicate counts and detailed error/duplicate information
- Duplicate detection: Resources with same type, model, base_url, and api_key are considered duplicates
- Automatic trimming: Leading and trailing whitespace are removed from all fields during import

**Resource Testing**:
- Test button available in the LLM Resources management interface
- Only resources with `active` status can be tested
- Supports testing for `chat` and `embedding` resource types
- Returns detailed test results including response time, model information, token usage, and response content
- Test timeout: 30 seconds

**Search Functionality**:
- Frontend supports fuzzy search on resource name, base URL, and model identifier
- Case-insensitive search with partial matching support

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
- `GET /api/v1/stats/system` - Get system statistics (total requests, total users, total LLM resources, success rate, average response time)
- `GET /api/v1/stats/llm-resources/usage` - Get LLM resource usage statistics (grouped by resource, includes token usage, request count, success rate, etc.)
- `GET /api/v1/stats/llm-resources/:id` - Get single LLM resource statistics
- `GET /api/v1/stats/users/:id` - Get user statistics

#### OpenAI Compatible API
- `GET /llm/v1/models` - List all available models
- `GET /llm/v1/models/:model` - Get model information
- `POST /llm/v1/chat/completions` - Create chat completion (supports streaming with `stream: true`)
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
  max_retries: 3   # Maximum retry count for failed requests (0 = disabled)
  retry_delay: "1s"  # Base retry delay between attempts (actual delay increases exponentially)
  max_idle_conns: 100  # Maximum idle connections
  max_conns_per_host: 100  # Maximum connections per host
  idle_conn_timeout: "90s"  # Idle connection timeout
```

**Retry Mechanism:**
- Automatically retries failed requests for network errors, timeouts, and 5xx server errors
- Uses exponential backoff: delay = retry_delay × attempt_number
- Does not retry 4xx client errors (except 429 rate limit), authentication errors, or context cancellations
- Configurable via admin interface: Settings → Provider Settings
- Applies to all request types: chat completions (streaming and non-streaming), text completions, and embeddings

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
│   │   └── balancer/      # Load balancing
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

// Token API Key model - request-side API key management
type Token struct {
    ID         string     // API Key unique identifier
    Name       string     // API Key name/description
    Token      string     // API Key value (prefixed with "ling-")
    Prefix     string     // API Key prefix (for display)
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
    Driver    string    // Driver (currently supports: openai)
    Model     string    // Model identifier (e.g., gpt-4, gpt-3.5-turbo)
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

    // API Key management
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

### Adding a New AI Driver

1. **Update LLM Resource Model**
Extend the Driver field validation in `internal/handler/provider.go` to support new driver types

2. **Implement Driver Client**
Create a new client implementation in `internal/client/` for the new driver

3. **Update Load Balancing Strategy**
Implement or update load balancing algorithms in `internal/pkg/balancer/` if needed

4. **Update Frontend**
Add the new driver option in the frontend LLM resource management interface

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

### v1.5.0 (2026-02-08)
- **Automatic Retry**: Added configurable automatic retry mechanism with exponential backoff for failed requests
- **Provider Configuration**: Added provider settings (timeout, max retries, retry delay) configurable via admin interface
- **Error Classification**: Intelligent error classification for retryable vs non-retryable errors
- **Streaming Retry**: Retry logic now applies to streaming requests before stream establishment
- **API Key Management**: Renamed "Token Management" to "API Key Management" across all documentation and UI to avoid confusion with LLM tokens
- **Documentation**: Comprehensive updates to all documentation (README, configuration guide, API reference, architecture)

### v1.4.0 (2026-02-05)
- **Internationalization**: Full frontend i18n support with Chinese and English language switching
- **Streaming Support**: Added Server-Sent Events (SSE) streaming support for chat completions
- **Policy Enhancement**: Random selection policy now supports LLM resource pool configuration
- **Client Libraries**: Standard client implementations added for Python, JavaScript, and Go
- **Code Cleanup**: Removed redundant `backend/examples` directory, unified client examples in `clients/` directory
- **Documentation**: Comprehensive documentation updates across all languages

### v1.3.0 (2026-02-03)
- **Driver Architecture**: Changed from "Provider" to "Driver" concept, currently supports OpenAI driver only
- **Batch Import/Export**: Added Excel template download and batch import functionality for LLM resources
- **Enhanced Search**: Added fuzzy search support for resource name, base URL, and model identifier
- **Frontend Improvements**: Fixed data display issues after batch import, improved search UX
- **Template Management**: Excel import template includes core fields (name, type, driver, model, base_url, api_key, status)

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