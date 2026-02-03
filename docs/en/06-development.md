# Development Guide

## Development Environment Setup

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm/yarn
- Git
- SQLite (for local development)

### Setup Steps

1. **Clone the repository**
```bash
git clone https://github.com/your-org/lingproxy.git
cd lingproxy
```

2. **Backend setup**
```bash
cd backend
go mod download
```

3. **Frontend setup**
```bash
cd frontend
npm install
```

4. **Configuration**
```bash
cp backend/configs/config.yaml.example backend/configs/config.yaml
# Edit config.yaml as needed
```

## Code Structure

### Backend Structure

Follow the existing structure:
- `cmd/`: Application entry point
- `internal/handler/`: HTTP handlers
- `internal/service/`: Business logic
- `internal/storage/`: Data persistence
- `internal/middleware/`: HTTP middleware
- `internal/config/`: Configuration management

### Frontend Structure

- `src/views/`: Page components
- `src/components/`: Reusable components
- `src/api/`: API client
- `src/router/`: Route configuration

## Coding Standards

### Go Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Follow naming conventions:
  - Exported functions: PascalCase
  - Unexported functions: camelCase
  - Constants: UPPER_SNAKE_CASE

### Error Handling

- Always check and handle errors
- Return errors from functions, don't ignore them
- Use descriptive error messages

### Testing

- Write unit tests for new features
- Test files should be named `*_test.go`
- Use table-driven tests where appropriate

Example:
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name string
        input string
        expected string
    }{
        {"test1", "input1", "expected1"},
        {"test2", "input2", "expected2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Function(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Adding New Features

### 1. Add a New Handler

1. Create handler file: `internal/handler/feature.go`
2. Implement handler methods
3. Register routes in `internal/router/router.go`

Example:
```go
package handler

type FeatureHandler struct {
    storage *storage.StorageFacade
}

func NewFeatureHandler(storage *storage.StorageFacade) *FeatureHandler {
    return &FeatureHandler{storage: storage}
}

func (h *FeatureHandler) GetFeature(c *gin.Context) {
    // Implementation
}
```

### 2. Add a New Service

1. Create service file: `internal/service/feature_service.go`
2. Implement business logic
3. Use storage layer for data operations

### 3. Add a New Storage Method

1. Add interface method in `internal/storage/storage.go`
2. Implement in `memory_storage.go` and `gorm_storage.go`
3. Add to `storage_facade.go`

### 4. Add a New Model

1. Add model definition in `internal/storage/models.go`
2. Update storage implementations
3. Add migration if needed

## Running Tests

### Backend Tests

```bash
cd backend
go test ./...
go test -v ./internal/handler  # Verbose output
go test -cover ./...            # Coverage report
```

### Frontend Tests

```bash
cd frontend
npm test
```

## Building

### Backend

```bash
cd backend
go build -o lingproxy cmd/main.go
```

### Frontend

```bash
cd frontend
npm run build
```

## Debugging

### Backend Debugging

1. Use logging: `logger.Debug()`, `logger.Info()`, etc.
2. Use Go debugger (delve)
3. Check logs in `./logs/lingproxy.log`

### Frontend Debugging

1. Use browser DevTools
2. Check console for errors
3. Use Vue DevTools extension

## Git Workflow

### Branch Naming

- `main`: Production-ready code
- `develop`: Development branch
- `feature/feature-name`: New features
- `bugfix/bug-name`: Bug fixes
- `hotfix/hotfix-name`: Hot fixes

### Commit Messages

Follow conventional commits:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build/tooling changes

Example:
```
feat: add batch import for LLM resources
fix: resolve authentication issue
docs: update API documentation
```

## Pull Request Process

1. Create a feature branch
2. Make changes and commit
3. Write/update tests
4. Update documentation if needed
5. Create pull request
6. Address review comments
7. Merge after approval

## Code Review Guidelines

### What to Review

- Code correctness
- Code style and formatting
- Test coverage
- Documentation updates
- Performance implications
- Security considerations

### Review Checklist

- [ ] Code follows style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No security vulnerabilities
- [ ] Performance is acceptable
- [ ] Error handling is proper

## Common Tasks

### Adding a New API Endpoint

1. Add handler method
2. Add route in router
3. Add Swagger documentation
4. Write tests
5. Update API documentation

### Adding a New Configuration Option

1. Add to `config.go` struct
2. Add default value in `setDefaults()`
3. Add validation if needed
4. Update configuration documentation

### Adding a New Database Field

1. Update model in `models.go`
2. Update storage implementations
3. Add migration if needed
4. Update API handlers/services

## Troubleshooting

### Common Issues

**Issue**: Go module errors
**Solution**: Run `go mod tidy`

**Issue**: Frontend build errors
**Solution**: Delete `node_modules` and `package-lock.json`, then `npm install`

**Issue**: Database connection errors
**Solution**: Check database configuration and ensure database is running

**Issue**: Port already in use
**Solution**: Change port in `config.yaml` or stop the process using the port

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Vue 3 Documentation](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
