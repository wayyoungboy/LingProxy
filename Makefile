# LingProxy Makefile

# 变量定义
BINARY_NAME=lingproxy
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Go 相关变量
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOPROXY ?= https://goproxy.cn,direct

# 目录
DIST_DIR=dist
BIN_DIR=bin
COVERAGE_DIR=coverage

.PHONY: all build clean test lint fmt vet deps help

# 默认目标
all: clean deps fmt vet lint test build

# 构建应用
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOPROXY=$(GOPROXY) go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/main.go

# 交叉编译
build-linux:
	@GOOS=linux GOARCH=amd64 $(MAKE) build

build-windows:
	@GOOS=windows GOARCH=amd64 $(MAKE) build

build-mac:
	@GOOS=darwin GOARCH=amd64 $(MAKE) build

build-arm:
	@GOOS=linux GOARCH=arm64 $(MAKE) build

# 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(DIST_DIR)
	@rm -rf $(COVERAGE_DIR)
	@rm -f coverage.out
	@go clean -cache -testcache -modcache

# 运行测试
test:
	@echo "Running tests..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o $(COVERAGE_DIR)/coverage.html

# 运行单元测试
test-unit:
	@echo "Running unit tests..."
	@go test -v -short ./...

# 运行集成测试
test-integration:
	@echo "Running integration tests..."
	@go test -v -run Integration ./...

# 代码覆盖率
coverage: test
	@go tool cover -func=coverage.out

# 代码格式化
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 代码检查
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, please install it"; \
	fi

# 静态分析
vet:
	@echo "Running go vet..."
	@go vet ./...

# 依赖管理
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# 更新依赖
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# 启动服务（后台运行）
start:
	@echo "Starting services in background..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "Use 'make stop' to stop services"
	@(cd frontend && npm run dev > /dev/null 2>&1 &) && \
	(go run ./cmd/main.go > /dev/null 2>&1 &) && \
	echo "Services started in background. PIDs:" && \
	pgrep -f "npm run dev" && pgrep -f "go run.*cmd/main.go"

# 启动后端服务（后台运行）
start-backend:
	@echo "Starting backend service in background..."
	@go run ./cmd/main.go > /dev/null 2>&1 &
	@echo "Backend started. PID: $$(pgrep -f 'go run.*cmd/main.go')"

# 启动前端服务（后台运行）
start-frontend:
	@echo "Starting frontend service in background..."
	@cd frontend && npm run dev > /dev/null 2>&1 &
	@echo "Frontend started. PID: $$(pgrep -f 'npm run dev')"

# 停止所有服务
stop:
	@echo "Stopping all services..."
	@echo "Stopping frontend services..."
	@pkill -f "npm run dev" 2>/dev/null || true
	@pkill -f "vite" 2>/dev/null || true
	@echo "Stopping backend services..."
	@pkill -f "go run.*cmd/main.go" 2>/dev/null || true
	@pkill -f "lingproxy" 2>/dev/null || true
	@pkill -f "/tmp/go-build.*exe/main" 2>/dev/null || true
	@echo "Releasing ports..."
	@lsof -ti:8080 2>/dev/null | xargs kill -9 2>/dev/null || true
	@lsof -ti:3000 2>/dev/null | xargs kill -9 2>/dev/null | grep -v "Cursor" || true
	@sleep 1
	@echo "All services stopped"

# 重启所有服务（前台运行）
restart: stop
	@echo "Waiting for services to stop..."
	@sleep 2
	@$(MAKE) run

# 运行应用（前后端同时启动，前台运行）
run:
	@echo "Starting both frontend and backend..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "Press Ctrl+C to stop both services"
	@trap 'kill 0' EXIT; \
	(cd frontend && npm run dev &) && \
	go run ./cmd/main.go

# 仅启动后端（前台运行）
run-backend:
	@echo "Running backend service..."
	@go run ./cmd/main.go

# 仅启动前端（前台运行）
run-frontend:
	@echo "Running frontend development server..."
	@cd frontend && npm run dev

# Docker 相关
docker-build:
	@echo "Building Docker image..."
	@docker build -t lingproxy:$(VERSION) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 lingproxy:$(VERSION)

# 生成API文档
docs:
	@echo "Generating API documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/main.go; \
	else \
		echo "swag not found, please install it: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# 性能测试
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# 安全扫描
security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found, please install it: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# 发布版本
release: clean deps test build
	@echo "Creating release..."
	@mkdir -p $(DIST_DIR)
	@tar -czf $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz -C $(BIN_DIR) $(BINARY_NAME)

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  all              - 运行完整的构建流程"
	@echo "  build            - 构建应用"
	@echo "  build-linux      - 构建Linux版本"
	@echo "  build-windows    - 构建Windows版本"
	@echo "  build-mac        - 构建macOS版本"
	@echo "  build-arm        - 构建ARM版本"
	@echo "  clean            - 清理构建产物"
	@echo "  test             - 运行所有测试"
	@echo "  test-unit        - 运行单元测试"
	@echo "  test-integration - 运行集成测试"
	@echo "  coverage         - 生成代码覆盖率报告"
	@echo "  fmt              - 格式化代码"
	@echo "  lint             - 运行代码检查"
	@echo "  vet              - 运行静态分析"
	@echo "  deps             - 下载依赖"
	@echo "  deps-update      - 更新依赖"
	@echo "  start            - 启动前后端服务（后台运行）"
	@echo "  start-backend    - 启动后端服务（后台运行）"
	@echo "  start-frontend   - 启动前端服务（后台运行）"
	@echo "  run              - 启动前后端服务（前台运行，推荐）"
	@echo "  run-backend      - 仅启动后端服务（前台运行）"
	@echo "  run-frontend     - 仅启动前端服务（前台运行）"
	@echo "  stop             - 停止所有服务"
	@echo "  restart          - 重启所有服务"
	@echo "  docker-build     - 构建Docker镜像"
	@echo "  docker-run       - 运行Docker容器"
	@echo "  docs             - 生成API文档"
	@echo "  bench            - 运行性能测试"
	@echo "  security         - 运行安全扫描"
	@echo "  release          - 创建发布版本"
	@echo "  help             - 显示帮助信息"