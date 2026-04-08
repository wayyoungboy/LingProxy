# LingProxy Makefile
# 用于构建、测试、运行和管理 LingProxy 项目

# ============================================================================
# 变量定义
# ============================================================================

# 项目信息
BINARY_NAME=lingproxy
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Go 相关变量
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOPROXY ?= https://goproxy.cn,direct
GO_CMD=go

# 目录定义
DIST_DIR=dist
BIN_DIR=bin
COVERAGE_DIR=coverage
RUN_DIR=run
LOGS_DIR=logs
FRONTEND_DIST=frontend/dist

# 端口配置
BACKEND_PORT ?= 8080
FRONTEND_PORT ?= 3000

# Docker Compose 命令适配（优先使用 V2，回退到 V1）
DOCKER_COMPOSE := $(shell command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1 && echo "docker compose" || echo "docker-compose")

# ============================================================================
# 默认目标
# ============================================================================

.DEFAULT_GOAL := help
.PHONY: all build build-backend build-frontend build-all clean test test-unit test-integration coverage fmt lint vet deps deps-update help

# ============================================================================
# 构建相关
# ============================================================================

## 完整构建流程（清理、依赖、格式化、检查、测试、构建）
all: clean deps fmt vet lint test build-all

## 构建后端应用
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	@cd backend && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOPROXY=$(GOPROXY) \
		$(GO_CMD) build $(LDFLAGS) -o ../$(BIN_DIR)/$(BINARY_NAME) ./cmd/main.go
	@echo "Backend built successfully: $(BIN_DIR)/$(BINARY_NAME)"

## 构建前端应用
build-frontend:
	@echo "Building frontend..."
	@if [ ! -d "frontend" ]; then \
		echo "Error: Frontend directory not found"; \
		exit 1; \
	fi
	@cd frontend && \
		if [ ! -f "package.json" ]; then \
			echo "Installing frontend dependencies..."; \
			npm install; \
		fi && \
		npm run build
	@echo "Frontend built successfully"

## 构建前后端（完整构建）
build-all: build build-frontend
	@echo "All components built successfully"

## 交叉编译 - Linux AMD64
build-linux:
	@GOOS=linux GOARCH=amd64 $(MAKE) build

## 交叉编译 - Windows AMD64
build-windows:
	@GOOS=windows GOARCH=amd64 $(MAKE) build

## 交叉编译 - macOS AMD64
build-mac:
	@GOOS=darwin GOARCH=amd64 $(MAKE) build

## 交叉编译 - Linux ARM64
build-arm:
	@GOOS=linux GOARCH=arm64 $(MAKE) build

## 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR) $(DIST_DIR) $(COVERAGE_DIR) $(RUN_DIR)
	@rm -f coverage.out
	@cd backend && $(GO_CMD) clean -cache -testcache -modcache 2>/dev/null || true
	@cd frontend && rm -rf dist node_modules/.vite 2>/dev/null || true
	@echo "Clean completed"

# ============================================================================
# 测试相关
# ============================================================================

## 运行所有测试
test:
	@echo "Running tests..."
	@mkdir -p $(COVERAGE_DIR)
	@cd backend && $(GO_CMD) test -v -race -coverprofile=../coverage.out ./...
	@$(GO_CMD) tool cover -html=coverage.out -o $(COVERAGE_DIR)/coverage.html 2>/dev/null || true
	@echo "Tests completed"

## 运行单元测试
test-unit:
	@echo "Running unit tests..."
	@cd backend && $(GO_CMD) test -v -short ./...

## 运行集成测试
test-integration:
	@echo "Running integration tests..."
	@cd backend && $(GO_CMD) test -v -run Integration ./...

## 生成代码覆盖率报告
coverage: test
	@echo "Code coverage report:"
	@$(GO_CMD) tool cover -func=coverage.out
	@echo ""
	@echo "HTML report: $(COVERAGE_DIR)/coverage.html"

## 运行性能基准测试
bench:
	@echo "Running benchmarks..."
	@cd backend && $(GO_CMD) test -bench=. -benchmem ./...

# ============================================================================
# 代码质量相关
# ============================================================================

## 格式化代码
fmt:
	@echo "Formatting code..."
	@cd backend && $(GO_CMD) fmt ./...
	@echo "Code formatted"

## 运行代码检查（需要 golangci-lint）
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd backend && golangci-lint run; \
	else \
		echo "Warning: golangci-lint not found, skipping..."; \
		echo "Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## 运行静态分析
vet:
	@echo "Running go vet..."
	@cd backend && $(GO_CMD) vet ./...
	@echo "Static analysis completed"

## 运行安全扫描（需要 gosec）
security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		cd backend && gosec ./...; \
	else \
		echo "Warning: gosec not found, skipping..."; \
		echo "Install: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# ============================================================================
# 依赖管理
# ============================================================================

## 下载并整理依赖
deps:
	@echo "Downloading dependencies..."
	@cd backend && $(GO_CMD) mod download && $(GO_CMD) mod tidy
	@if [ -d "frontend" ] && [ -f "frontend/package.json" ]; then \
		echo "Installing frontend dependencies..."; \
		cd frontend && npm install; \
	fi
	@echo "Dependencies installed"

## 更新依赖到最新版本
deps-update:
	@echo "Updating dependencies..."
	@cd backend && $(GO_CMD) get -u ./... && $(GO_CMD) mod tidy
	@if [ -d "frontend" ] && [ -f "frontend/package.json" ]; then \
		cd frontend && npm update; \
	fi
	@echo "Dependencies updated"

# ============================================================================
# 服务管理（后台运行）
# ============================================================================

## 启动所有服务（后台运行）
start: start-backend start-frontend
	@echo "All services started"
	@echo "Backend: http://localhost:$(BACKEND_PORT)"
	@echo "Frontend: http://localhost:$(FRONTEND_PORT)"
	@echo "Use 'make stop' to stop services"

## 启动后端服务（后台运行）
start-backend:
	@echo "Starting backend service..."
	@mkdir -p $(RUN_DIR) $(LOGS_DIR)
	@if [ -f $(RUN_DIR)/backend.pid ]; then \
		OLD_PID=$$(cat $(RUN_DIR)/backend.pid 2>/dev/null); \
		if [ -n "$$OLD_PID" ] && ps -p $$OLD_PID > /dev/null 2>&1; then \
			echo "Warning: Backend already running (PID: $$OLD_PID)"; \
			exit 0; \
		else \
			rm -f $(RUN_DIR)/backend.pid; \
		fi; \
	fi
	@cd backend && $(GO_CMD) run ./cmd/main.go > ../$(LOGS_DIR)/backend.log 2>&1 &
	@sleep 2
	@BACKEND_PID=$$(pgrep -f "go run.*cmd/main.go" | head -1); \
	if [ -z "$$BACKEND_PID" ]; then \
		BACKEND_PID=$$(pgrep -f "lingproxy" | head -1); \
	fi; \
	if [ -n "$$BACKEND_PID" ]; then \
		echo $$BACKEND_PID > $(RUN_DIR)/backend.pid; \
		echo "Backend started (PID: $$BACKEND_PID)"; \
	else \
		echo "Error: Failed to start backend"; \
		echo "Check logs: tail -f $(LOGS_DIR)/backend.log"; \
		exit 1; \
	fi

## 启动前端服务（后台运行）
start-frontend:
	@echo "Starting frontend service..."
	@mkdir -p $(RUN_DIR) $(LOGS_DIR)
	@if [ -f $(RUN_DIR)/frontend.pid ]; then \
		OLD_PID=$$(cat $(RUN_DIR)/frontend.pid 2>/dev/null); \
		if [ -n "$$OLD_PID" ] && ps -p $$OLD_PID > /dev/null 2>&1; then \
			echo "Warning: Frontend already running (PID: $$OLD_PID)"; \
			exit 0; \
		else \
			rm -f $(RUN_DIR)/frontend.pid; \
		fi; \
	fi
	@if [ ! -d "frontend" ]; then \
		echo "Error: Frontend directory not found"; \
		exit 1; \
	fi
	@cd frontend && npm run dev > ../$(LOGS_DIR)/frontend.log 2>&1 &
	@sleep 2
	@FRONTEND_PID=$$(pgrep -f "vite" | grep -v "grep" | head -1); \
	if [ -n "$$FRONTEND_PID" ]; then \
		echo $$FRONTEND_PID > $(RUN_DIR)/frontend.pid; \
		echo "Frontend started (PID: $$FRONTEND_PID)"; \
	else \
		echo "Error: Failed to start frontend"; \
		echo "Check logs: tail -f $(LOGS_DIR)/frontend.log"; \
		exit 1; \
	fi

## 停止所有服务
stop:
	@echo "Stopping services..."
	@$(MAKE) stop-frontend
	@$(MAKE) stop-backend
	@echo "All services stopped"

## 停止后端服务
stop-backend:
	@if [ -f $(RUN_DIR)/backend.pid ]; then \
		BACKEND_PID=$$(cat $(RUN_DIR)/backend.pid 2>/dev/null); \
		if [ -n "$$BACKEND_PID" ] && ps -p $$BACKEND_PID > /dev/null 2>&1; then \
			echo "Stopping backend (PID: $$BACKEND_PID)..."; \
			kill $$BACKEND_PID 2>/dev/null || true; \
			sleep 1; \
			if ps -p $$BACKEND_PID > /dev/null 2>&1; then \
				kill -9 $$BACKEND_PID 2>/dev/null || true; \
			fi; \
		fi; \
		rm -f $(RUN_DIR)/backend.pid; \
	fi
	@pkill -f "go run.*cmd/main.go" 2>/dev/null || true
	@pkill -f "lingproxy" 2>/dev/null || true
	@lsof -ti:$(BACKEND_PORT) 2>/dev/null | xargs kill -9 2>/dev/null || true

## 停止前端服务
stop-frontend:
	@if [ -f $(RUN_DIR)/frontend.pid ]; then \
		FRONTEND_PID=$$(cat $(RUN_DIR)/frontend.pid 2>/dev/null); \
		if [ -n "$$FRONTEND_PID" ] && ps -p $$FRONTEND_PID > /dev/null 2>&1; then \
			echo "Stopping frontend (PID: $$FRONTEND_PID)..."; \
			kill $$FRONTEND_PID 2>/dev/null || true; \
			sleep 1; \
			if ps -p $$FRONTEND_PID > /dev/null 2>&1; then \
				kill -9 $$FRONTEND_PID 2>/dev/null || true; \
			fi; \
		fi; \
		rm -f $(RUN_DIR)/frontend.pid; \
	fi
	@pkill -f "vite" 2>/dev/null || true
	@lsof -ti:$(FRONTEND_PORT) 2>/dev/null | xargs kill -9 2>/dev/null || true

## 重启所有服务
restart: stop
	@sleep 2
	@$(MAKE) start

## 查看服务状态
status:
	@echo "Service Status:"
	@if [ -f $(RUN_DIR)/backend.pid ]; then \
		BACKEND_PID=$$(cat $(RUN_DIR)/backend.pid 2>/dev/null); \
		if [ -n "$$BACKEND_PID" ] && ps -p $$BACKEND_PID > /dev/null 2>&1; then \
			echo "Backend: Running (PID: $$BACKEND_PID)"; \
		else \
			echo "Backend: Not running"; \
		fi; \
	else \
		echo "Backend: Not running"; \
	fi
	@if [ -f $(RUN_DIR)/frontend.pid ]; then \
		FRONTEND_PID=$$(cat $(RUN_DIR)/frontend.pid 2>/dev/null); \
		if [ -n "$$FRONTEND_PID" ] && ps -p $$FRONTEND_PID > /dev/null 2>&1; then \
			echo "Frontend: Running (PID: $$FRONTEND_PID)"; \
		else \
			echo "Frontend: Not running"; \
		fi; \
	else \
		echo "Frontend: Not running"; \
	fi

# ============================================================================
# 开发运行（前台运行）
# ============================================================================

## 运行前后端（前台运行，推荐开发使用）
run:
	@echo "Starting services in foreground..."
	@echo "Backend: http://localhost:$(BACKEND_PORT)"
	@echo "Frontend: http://localhost:$(FRONTEND_PORT)"
	@echo "Press Ctrl+C to stop"
	@trap 'kill 0' EXIT; \
	(cd frontend && npm run dev &) && \
	cd backend && $(GO_CMD) run ./cmd/main.go

## 仅运行后端（前台运行）
run-backend:
	@echo "Running backend service..."
	@echo "Note: For local development, set LINGPROXY_STORAGE_GORM_DSN=root:@tcp(localhost:2881)/lingproxy?charset=utf8mb4&parseTime=True&loc=Local"
	@cd backend && $(GO_CMD) run ./cmd/main.go

## 仅运行前端（前台运行）
run-frontend:
	@echo "Running frontend development server..."
	@cd frontend && npm run dev

# ============================================================================
# 日志管理
# ============================================================================

## 查看后端日志
logs-backend:
	@tail -f $(LOGS_DIR)/backend.log 2>/dev/null || echo "No backend log file found"

## 查看前端日志
logs-frontend:
	@tail -f $(LOGS_DIR)/frontend.log 2>/dev/null || echo "No frontend log file found"

## 查看所有日志
logs:
	@echo "Backend logs:"
	@tail -20 $(LOGS_DIR)/backend.log 2>/dev/null || echo "No backend log"
	@echo ""
	@echo "Frontend logs:"
	@tail -20 $(LOGS_DIR)/frontend.log 2>/dev/null || echo "No frontend log"

# ============================================================================
# Docker 相关
# ============================================================================

## 构建 Docker 镜像（使用 Docker Compose）
docker-build:
	@echo "Building Docker images with Docker Compose..."
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml build
	@echo "Docker images built successfully"

## 运行 Docker 容器（使用 Docker Compose）
docker-run: docker-compose-up
	@echo "Services are running"

## 检查 Docker Compose 配置文件
docker-compose-check:
	@echo "Checking Docker Compose configuration..."
	@if [ ! -f "backend/configs/config.yaml" ]; then \
		echo "Warning: Config file not found, creating from example..."; \
		cp backend/configs/config.yaml.example backend/configs/config.yaml; \
		echo "Warning: Please edit backend/configs/config.yaml and configure SeekDB connection"; \
	fi
	@if [ ! -d "logs" ]; then mkdir -p logs; fi
	@if [ ! -d "run" ]; then mkdir -p run; fi
	@echo "Configuration check completed"

## 创建 SeekDB 数据库
docker-compose-init-db:
	@echo "Initializing SeekDB database..."
	@echo "Waiting for SeekDB to be ready..."
	@sleep 10
	@docker exec seekdb mysql -h127.0.0.1 -uroot -P2881 -e "CREATE DATABASE IF NOT EXISTS lingproxy CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null && \
		echo "Database 'lingproxy' created successfully" || \
		echo "Warning: Database creation failed or already exists"

## 使用 Docker Compose 启动（完整流程）
docker-compose-up: docker-compose-check
	@echo "Starting services with Docker Compose..."
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml up -d --build
	@echo "Services started"
	@echo "Waiting for SeekDB to be ready..."
	@sleep 15
	@$(MAKE) docker-compose-init-db
	@echo "Services are running"
	@echo "Access URLs:"
	@echo "  Backend API: http://localhost:8080/api/v1"
	@echo "  Health Check: http://localhost:8080/api/v1/health"
	@echo "Note: Frontend should be run separately with 'cd frontend && npm run dev'"

## 停止 Docker Compose 服务
docker-compose-down:
	@echo "Stopping Docker Compose services..."
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml down
	@echo "Services stopped"

## 查看 Docker Compose 日志
docker-compose-logs:
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml logs -f

## 查看 Docker Compose 服务状态
docker-compose-ps:
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml ps

## 重启 Docker Compose 服务
docker-compose-restart:
	@echo "Restarting Docker Compose services..."
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml restart
	@echo "Services restarted"

## 重新构建并启动 Docker Compose 服务（停止 → 构建 → 启动）
docker-rebuild: docker-compose-down
	@echo "Rebuilding and starting Docker Compose services..."
	@$(MAKE) docker-compose-up
	@echo "Services rebuilt and started"

## 完全清理 Docker Compose 服务和挂载卷
docker-clean:
	@echo "Cleaning Docker Compose services and volumes..."
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml down -v --remove-orphans
	@echo "Docker Compose services and volumes cleaned"

## 清理、重新构建并启动 Docker Compose 服务（完全清理 → 构建 → 启动）
docker-clean-rebuild: docker-clean docker-build docker-compose-up
	@echo "Services cleaned, rebuilt and started"

# ============================================================================
# 文档生成
# ============================================================================

## 生成 API 文档（Swagger）
docs:
	@echo "Generating API documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		cd backend && swag init -g cmd/main.go -o swagger; \
		echo "API documentation generated"; \
	else \
		echo "Warning: swag not found"; \
		echo "Install: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# ============================================================================
# 官网构建相关
# ============================================================================

## 安装官网依赖
website-install:
	@echo "Installing website dependencies..."
	@cd website && npm install
	@echo "Dependencies installed"

## 启动官网开发服务器
website-start:
	@echo "Starting website development server..."
	@cd website && npm start

## 构建官网
website-build:
	@echo "Building website..."
	@cd website && npm run build
	@echo "Website built successfully. Files are in website/build/"

## 预览构建后的官网
website-serve:
	@echo "Serving built website..."
	@cd website && npm run serve

## 清理官网构建缓存
website-clean:
	@echo "Cleaning website cache..."
	@cd website && npm run clear
	@echo "Cache cleaned"

## 部署官网到 GitHub Pages
website-deploy:
	@echo "Deploying website to GitHub Pages..."
	@if [ -z "$(GIT_USER)" ]; then \
		echo "Error: GIT_USER environment variable is required"; \
		echo "Example: GIT_USER=your-username make website-deploy"; \
		exit 1; \
	fi
	@cd website && GIT_USER=$(GIT_USER) USE_SSH=true npm run deploy
	@echo "Website deployed successfully"

# ============================================================================
# 发布相关
# ============================================================================

## 创建发布版本
release: clean deps test build-all
	@echo "Creating release..."
	@mkdir -p $(DIST_DIR)
	@tar -czf $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz \
		-C $(BIN_DIR) $(BINARY_NAME) 2>/dev/null || true
	@if [ -d "$(FRONTEND_DIST)" ]; then \
		tar -czf $(DIST_DIR)/frontend-$(VERSION).tar.gz -C frontend dist; \
	fi
	@echo "Release created in $(DIST_DIR)/"

# ============================================================================
# 工具和实用功能
# ============================================================================

## 检查开发环境
check-env:
	@echo "Checking development environment..."
	@echo "Go version: $$($(GO_CMD) version 2>/dev/null || echo 'Not installed')"
	@echo "Node version: $$(node --version 2>/dev/null || echo 'Not installed')"
	@echo "NPM version: $$(npm --version 2>/dev/null || echo 'Not installed')"
	@echo "Docker version: $$(docker --version 2>/dev/null || echo 'Not installed')"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint: Installed"; \
	else \
		echo "golangci-lint: Not installed"; \
	fi
	@if command -v swag >/dev/null 2>&1; then \
		echo "swag: Installed"; \
	else \
		echo "swag: Not installed"; \
	fi

## 初始化开发环境
init: deps
	@echo "Initializing development environment..."
	@mkdir -p $(LOGS_DIR) $(RUN_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	@if [ ! -f "backend/configs/config.yaml" ]; then \
		echo "Warning: Creating config.yaml from example..."; \
		cp backend/configs/config.yaml.example backend/configs/config.yaml; \
	fi
	@echo "Development environment initialized"

# ============================================================================
# 帮助信息
# ============================================================================

## 显示帮助信息
help:
	@echo "LingProxy Makefile Commands"
	@echo ""
	@echo "构建相关:"
	@echo "  make all              - 完整构建流程（清理、依赖、格式化、检查、测试、构建）"
	@echo "  make build            - 构建后端应用"
	@echo "  make build-frontend   - 构建前端应用"
	@echo "  make build-all        - 构建前后端"
	@echo "  make build-linux      - 构建 Linux 版本"
	@echo "  make build-windows    - 构建 Windows 版本"
	@echo "  make build-mac        - 构建 macOS 版本"
	@echo "  make build-arm        - 构建 ARM 版本"
	@echo "  make clean            - 清理构建产物"
	@echo ""
	@echo "测试相关:"
	@echo "  make test             - 运行所有测试"
	@echo "  make test-unit        - 运行单元测试"
	@echo "  make test-integration - 运行集成测试"
	@echo "  make coverage         - 生成代码覆盖率报告"
	@echo "  make bench            - 运行性能基准测试"
	@echo ""
	@echo "代码质量:"
	@echo "  make fmt              - 格式化代码"
	@echo "  make lint             - 运行代码检查（需要 golangci-lint）"
	@echo "  make vet              - 运行静态分析"
	@echo "  make security         - 运行安全扫描（需要 gosec）"
	@echo ""
	@echo "依赖管理:"
	@echo "  make deps             - 下载并整理依赖"
	@echo "  make deps-update      - 更新依赖到最新版本"
	@echo ""
	@echo "服务管理（后台）:"
	@echo "  make start            - 启动所有服务（后台运行）"
	@echo "  make start-backend    - 启动后端服务（后台运行）"
	@echo "  make start-frontend   - 启动前端服务（后台运行）"
	@echo "  make stop             - 停止所有服务"
	@echo "  make stop-backend     - 停止后端服务"
	@echo "  make stop-frontend    - 停止前端服务"
	@echo "  make restart          - 重启所有服务"
	@echo "  make status           - 查看服务状态"
	@echo ""
	@echo "开发运行（前台）:"
	@echo "  make run              - 运行前后端（前台运行，推荐开发使用）"
	@echo "  make run-backend      - 仅运行后端（前台运行）"
	@echo "  make run-frontend     - 仅运行前端（前台运行）"
	@echo ""
	@echo "日志管理:"
	@echo "  make logs             - 查看所有日志（最后20行）"
	@echo "  make logs-backend     - 实时查看后端日志"
	@echo "  make logs-frontend    - 实时查看前端日志"
	@echo ""
	@echo "Docker 相关:"
	@echo "  make docker-build            - 构建 Docker 镜像"
	@echo "  make docker-run              - 运行 Docker 容器"
	@echo "  make docker-compose-up      - 启动服务（包含检查配置、创建数据库）"
	@echo "  make docker-compose-down     - 停止 Docker Compose 服务"
	@echo "  make docker-compose-logs     - 查看 Docker Compose 日志"
	@echo "  make docker-compose-ps       - 查看服务状态"
	@echo "  make docker-compose-restart  - 重启服务（不重新构建）"
	@echo "  make docker-rebuild          - 重新构建并启动服务（停止 → 构建 → 启动）"
	@echo "  make docker-clean            - 完全清理 Docker Compose 服务和挂载卷"
	@echo "  make docker-clean-rebuild    - 清理、重新构建并启动服务（完全清理 → 构建 → 启动）"
	@echo ""
	@echo "文档和工具:"
	@echo "  make docs              - 生成 API 文档（Swagger）"
	@echo "  make website-install   - 安装官网依赖"
	@echo "  make website-start     - 启动官网开发服务器"
	@echo "  make website-build     - 构建官网"
	@echo "  make website-serve     - 预览构建后的官网"
	@echo "  make website-deploy    - 部署官网到 GitHub Pages"
	@echo "  make release           - 创建发布版本"
	@echo "  make check-env         - 检查开发环境"
	@echo "  make init              - 初始化开发环境"
	@echo "  make help              - 显示帮助信息"
	@echo ""
	@echo "环境变量:"
	@echo "  BACKEND_PORT=$(BACKEND_PORT)  - 后端服务端口"
	@echo "  FRONTEND_PORT=$(FRONTEND_PORT) - 前端服务端口"
	@echo "  GOOS=$(GOOS)                  - Go 目标操作系统"
	@echo "  GOARCH=$(GOARCH)              - Go 目标架构"
	@echo "  VERSION=$(VERSION)            - 版本号"
