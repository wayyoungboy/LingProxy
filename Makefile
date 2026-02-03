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
RUN_DIR=run

.PHONY: all build clean test lint fmt vet deps help

# 默认目标
all: clean deps fmt vet lint test build

# 构建应用
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	@cd backend && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOPROXY=$(GOPROXY) go build $(LDFLAGS) -o ../$(BIN_DIR)/$(BINARY_NAME) ./cmd/main.go

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
	@rm -rf $(RUN_DIR)
	@rm -f coverage.out
	@cd backend && go clean -cache -testcache -modcache

# 运行测试
test:
	@echo "Running tests..."
	@mkdir -p $(COVERAGE_DIR)
	@cd backend && go test -v -race -coverprofile=../coverage.out ./...
	@go tool cover -html=coverage.out -o $(COVERAGE_DIR)/coverage.html

# 运行单元测试
test-unit:
	@echo "Running unit tests..."
	@cd backend && go test -v -short ./...

# 运行集成测试
test-integration:
	@echo "Running integration tests..."
	@cd backend && go test -v -run Integration ./...

# 代码覆盖率
coverage: test
	@go tool cover -func=coverage.out

# 代码格式化
fmt:
	@echo "Formatting code..."
	@cd backend && go fmt ./...

# 代码检查
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd backend && golangci-lint run; \
	else \
		echo "golangci-lint not found, please install it"; \
	fi

# 静态分析
vet:
	@echo "Running go vet..."
	@cd backend && go vet ./...

# 依赖管理
deps:
	@echo "Downloading dependencies..."
	@cd backend && go mod download
	@cd backend && go mod tidy

# 更新依赖
deps-update:
	@echo "Updating dependencies..."
	@cd backend && go get -u ./...
	@cd backend && go mod tidy

# 启动服务（后台运行）
start:
	@echo "Starting services in background..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "Use 'make stop' to stop services"
	@$(MAKE) start-frontend
	@$(MAKE) start-backend
	@echo "Services started. PIDs saved to $(RUN_DIR)/"

# 启动后端服务（后台运行）
start-backend:
	@echo "Starting backend service in background..."
	@mkdir -p $(RUN_DIR)
	@if [ -f $(RUN_DIR)/backend.pid ]; then \
		OLD_PID=$$(cat $(RUN_DIR)/backend.pid); \
		if ps -p $$OLD_PID > /dev/null 2>&1; then \
			echo "Backend already running (PID: $$OLD_PID)"; \
			exit 0; \
		else \
			echo "Removing stale PID file"; \
			rm -f $(RUN_DIR)/backend.pid; \
		fi; \
	fi
	@cd backend && go run ./cmd/main.go > /dev/null 2>&1 &
	@sleep 2
	@BACKEND_PID=$$(pgrep -f "go run.*cmd/main.go" | head -1); \
	if [ -z "$$BACKEND_PID" ]; then \
		BACKEND_PID=$$(pgrep -f "lingproxy" | head -1); \
	fi; \
	if [ -n "$$BACKEND_PID" ]; then \
		echo $$BACKEND_PID > $(RUN_DIR)/backend.pid; \
		echo "Backend started. PID: $$BACKEND_PID"; \
	else \
		echo "Warning: Failed to get backend PID, but service may be running"; \
		echo "Check with: lsof -ti:8080"; \
	fi

# 启动前端服务（后台运行）
start-frontend:
	@echo "Starting frontend service in background..."
	@mkdir -p $(RUN_DIR)
	@if [ -f $(RUN_DIR)/frontend.pid ]; then \
		OLD_PID=$$(cat $(RUN_DIR)/frontend.pid); \
		if ps -p $$OLD_PID > /dev/null 2>&1; then \
			echo "Frontend already running (PID: $$OLD_PID)"; \
			exit 0; \
		else \
			echo "Removing stale PID file"; \
			rm -f $(RUN_DIR)/frontend.pid; \
		fi; \
	fi
	@cd frontend && npm run dev > /dev/null 2>&1 &
	@sleep 1
	@FRONTEND_PID=$$(pgrep -f "npm run dev" | head -1); \
	if [ -n "$$FRONTEND_PID" ]; then \
		echo $$FRONTEND_PID > $(RUN_DIR)/frontend.pid; \
		echo "Frontend started. PID: $$FRONTEND_PID"; \
	else \
		echo "Failed to start frontend or get PID"; \
		exit 1; \
	fi

# 停止所有服务
stop:
	@echo "Stopping all services..."
	@if [ -f $(RUN_DIR)/frontend.pid ]; then \
		FRONTEND_PID=$$(cat $(RUN_DIR)/frontend.pid); \
		if ps -p $$FRONTEND_PID > /dev/null 2>&1; then \
			echo "Stopping frontend (PID: $$FRONTEND_PID)..."; \
			kill $$FRONTEND_PID 2>/dev/null || true; \
			sleep 1; \
			if ps -p $$FRONTEND_PID > /dev/null 2>&1; then \
				echo "Force killing frontend (PID: $$FRONTEND_PID)..."; \
				kill -9 $$FRONTEND_PID 2>/dev/null || true; \
			fi; \
		else \
			echo "Frontend process (PID: $$FRONTEND_PID) not running"; \
		fi; \
		rm -f $(RUN_DIR)/frontend.pid; \
	else \
		echo "No frontend PID file found, trying to kill by process name..."; \
		pkill -f "npm run dev" 2>/dev/null || true; \
		pkill -f "vite" 2>/dev/null || true; \
	fi
	@if [ -f $(RUN_DIR)/backend.pid ]; then \
		BACKEND_PID=$$(cat $(RUN_DIR)/backend.pid); \
		if ps -p $$BACKEND_PID > /dev/null 2>&1; then \
			echo "Stopping backend (PID: $$BACKEND_PID)..."; \
			kill $$BACKEND_PID 2>/dev/null || true; \
			sleep 1; \
			if ps -p $$BACKEND_PID > /dev/null 2>&1; then \
				echo "Force killing backend (PID: $$BACKEND_PID)..."; \
				kill -9 $$BACKEND_PID 2>/dev/null || true; \
			fi; \
		else \
			echo "Backend process (PID: $$BACKEND_PID) not running"; \
		fi; \
		rm -f $(RUN_DIR)/backend.pid; \
	else \
		echo "No backend PID file found, trying to kill by process name..."; \
		pkill -f "go run.*backend/cmd/main.go" 2>/dev/null || true; \
		pkill -f "go run.*cmd/main.go" 2>/dev/null || true; \
		pkill -f "lingproxy" 2>/dev/null || true; \
		pkill -f "/tmp/go-build.*exe/main" 2>/dev/null || true; \
	fi
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
	cd backend && go run ./cmd/main.go

# 仅启动后端（前台运行）
run-backend:
	@echo "Running backend service..."
	@cd backend && go run ./cmd/main.go

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
		cd backend && swag init -g cmd/main.go -o swagger; \
	else \
		echo "swag not found, please install it: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# 性能测试
bench:
	@echo "Running benchmarks..."
	@cd backend && go test -bench=. -benchmem ./...

# 安全扫描
security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		cd backend && gosec ./...; \
	else \
		echo "gosec not found, please install it: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# 发布版本
release: clean deps test build
	@echo "Creating release..."
	@mkdir -p $(DIST_DIR)
	@tar -czf $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz -C $(BIN_DIR) $(BINARY_NAME)

# 生成Excel导入模板
generate-template:
	@echo "Generating Excel import template..."
	@cd backend && go run cmd/generate_template.go

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
	@echo "  generate-template - 生成Excel导入模板文件"
	@echo "  docker-build     - 构建Docker镜像"
	@echo "  docker-run       - 运行Docker容器"
	@echo "  docs             - 生成API文档"
	@echo "  bench            - 运行性能测试"
	@echo "  security         - 运行安全扫描"
	@echo "  release          - 创建发布版本"
	@echo "  help             - 显示帮助信息"