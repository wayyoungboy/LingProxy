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

# 颜色输出
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_RED=\033[31m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

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
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Building $(BINARY_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BIN_DIR)
	@cd backend && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOPROXY=$(GOPROXY) \
		$(GO_CMD) build $(LDFLAGS) -o ../$(BIN_DIR)/$(BINARY_NAME) ./cmd/main.go
	@echo "$(COLOR_GREEN)✓ Backend built successfully: $(BIN_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

## 构建前端应用
build-frontend:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Building frontend...$(COLOR_RESET)"
	@if [ ! -d "frontend" ]; then \
		echo "$(COLOR_RED)✗ Frontend directory not found$(COLOR_RESET)"; \
		exit 1; \
	fi
	@cd frontend && \
		if [ ! -f "package.json" ]; then \
			echo "$(COLOR_YELLOW)Installing frontend dependencies...$(COLOR_RESET)"; \
			npm install; \
		fi && \
		npm run build
	@echo "$(COLOR_GREEN)✓ Frontend built successfully$(COLOR_RESET)"

## 构建前后端（完整构建）
build-all: build build-frontend
	@echo "$(COLOR_GREEN)✓ All components built successfully$(COLOR_RESET)"

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
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf $(BIN_DIR) $(DIST_DIR) $(COVERAGE_DIR) $(RUN_DIR)
	@rm -f coverage.out
	@cd backend && $(GO_CMD) clean -cache -testcache -modcache 2>/dev/null || true
	@cd frontend && rm -rf dist node_modules/.vite 2>/dev/null || true
	@echo "$(COLOR_GREEN)✓ Clean completed$(COLOR_RESET)"

# ============================================================================
# 测试相关
# ============================================================================

## 运行所有测试
test:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@cd backend && $(GO_CMD) test -v -race -coverprofile=../coverage.out ./...
	@$(GO_CMD) tool cover -html=coverage.out -o $(COVERAGE_DIR)/coverage.html 2>/dev/null || true
	@echo "$(COLOR_GREEN)✓ Tests completed$(COLOR_RESET)"

## 运行单元测试
test-unit:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Running unit tests...$(COLOR_RESET)"
	@cd backend && $(GO_CMD) test -v -short ./...

## 运行集成测试
test-integration:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Running integration tests...$(COLOR_RESET)"
	@cd backend && $(GO_CMD) test -v -run Integration ./...

## 生成代码覆盖率报告
coverage: test
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Code coverage report:$(COLOR_RESET)"
	@$(GO_CMD) tool cover -func=coverage.out
	@echo "\n$(COLOR_GREEN)HTML report: $(COVERAGE_DIR)/coverage.html$(COLOR_RESET)"

## 运行性能基准测试
bench:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Running benchmarks...$(COLOR_RESET)"
	@cd backend && $(GO_CMD) test -bench=. -benchmem ./...

# ============================================================================
# 代码质量相关
# ============================================================================

## 格式化代码
fmt:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	@cd backend && $(GO_CMD) fmt ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

## 运行代码检查（需要 golangci-lint）
lint:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd backend && golangci-lint run; \
	else \
		echo "$(COLOR_YELLOW)⚠ golangci-lint not found, skipping...$(COLOR_RESET)"; \
		echo "Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## 运行静态分析
vet:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	@cd backend && $(GO_CMD) vet ./...
	@echo "$(COLOR_GREEN)✓ Static analysis completed$(COLOR_RESET)"

## 运行安全扫描（需要 gosec）
security:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Running security scan...$(COLOR_RESET)"
	@if command -v gosec >/dev/null 2>&1; then \
		cd backend && gosec ./...; \
	else \
		echo "$(COLOR_YELLOW)⚠ gosec not found, skipping...$(COLOR_RESET)"; \
		echo "Install: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# ============================================================================
# 依赖管理
# ============================================================================

## 下载并整理依赖
deps:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	@cd backend && $(GO_CMD) mod download && $(GO_CMD) mod tidy
	@if [ -d "frontend" ] && [ -f "frontend/package.json" ]; then \
		echo "$(COLOR_BLUE)Installing frontend dependencies...$(COLOR_RESET)"; \
		cd frontend && npm install; \
	fi
	@echo "$(COLOR_GREEN)✓ Dependencies installed$(COLOR_RESET)"

## 更新依赖到最新版本
deps-update:
	@echo "$(COLOR_BOLD)$(COLOR_YELLOW)Updating dependencies...$(COLOR_RESET)"
	@cd backend && $(GO_CMD) get -u ./... && $(GO_CMD) mod tidy
	@if [ -d "frontend" ] && [ -f "frontend/package.json" ]; then \
		cd frontend && npm update; \
	fi
	@echo "$(COLOR_GREEN)✓ Dependencies updated$(COLOR_RESET)"

# ============================================================================
# 服务管理（后台运行）
# ============================================================================

## 启动所有服务（后台运行）
start: start-backend start-frontend
	@echo "$(COLOR_GREEN)✓ All services started$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)Backend:$(COLOR_RESET) http://localhost:$(BACKEND_PORT)"
	@echo "$(COLOR_BOLD)Frontend:$(COLOR_RESET) http://localhost:$(FRONTEND_PORT)"
	@echo "Use 'make stop' to stop services"

## 启动后端服务（后台运行）
start-backend:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Starting backend service...$(COLOR_RESET)"
	@mkdir -p $(RUN_DIR) $(LOGS_DIR)
	@if [ -f $(RUN_DIR)/backend.pid ]; then \
		OLD_PID=$$(cat $(RUN_DIR)/backend.pid 2>/dev/null); \
		if [ -n "$$OLD_PID" ] && ps -p $$OLD_PID > /dev/null 2>&1; then \
			echo "$(COLOR_YELLOW)⚠ Backend already running (PID: $$OLD_PID)$(COLOR_RESET)"; \
			exit 0; \
		else \
			rm -f $(RUN_DIR)/backend.pid; \
		fi; \
	fi
	@cd backend && $(GO_CMD) run ./cmd/main.go > $(LOGS_DIR)/backend.log 2>&1 &
	@sleep 2
	@BACKEND_PID=$$(pgrep -f "go run.*cmd/main.go" | head -1); \
	if [ -z "$$BACKEND_PID" ]; then \
		BACKEND_PID=$$(pgrep -f "lingproxy" | head -1); \
	fi; \
	if [ -n "$$BACKEND_PID" ]; then \
		echo $$BACKEND_PID > $(RUN_DIR)/backend.pid; \
		echo "$(COLOR_GREEN)✓ Backend started (PID: $$BACKEND_PID)$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_RED)✗ Failed to start backend$(COLOR_RESET)"; \
		echo "Check logs: tail -f $(LOGS_DIR)/backend.log"; \
		exit 1; \
	fi

## 启动前端服务（后台运行）
start-frontend:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Starting frontend service...$(COLOR_RESET)"
	@mkdir -p $(RUN_DIR) $(LOGS_DIR)
	@if [ -f $(RUN_DIR)/frontend.pid ]; then \
		OLD_PID=$$(cat $(RUN_DIR)/frontend.pid 2>/dev/null); \
		if [ -n "$$OLD_PID" ] && ps -p $$OLD_PID > /dev/null 2>&1; then \
			echo "$(COLOR_YELLOW)⚠ Frontend already running (PID: $$OLD_PID)$(COLOR_RESET)"; \
			exit 0; \
		else \
			rm -f $(RUN_DIR)/frontend.pid; \
		fi; \
	fi
	@if [ ! -d "frontend" ]; then \
		echo "$(COLOR_RED)✗ Frontend directory not found$(COLOR_RESET)"; \
		exit 1; \
	fi
	@cd frontend && npm run dev > ../$(LOGS_DIR)/frontend.log 2>&1 &
	@sleep 2
	@FRONTEND_PID=$$(pgrep -f "vite" | grep -v "grep" | head -1); \
	if [ -n "$$FRONTEND_PID" ]; then \
		echo $$FRONTEND_PID > $(RUN_DIR)/frontend.pid; \
		echo "$(COLOR_GREEN)✓ Frontend started (PID: $$FRONTEND_PID)$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_RED)✗ Failed to start frontend$(COLOR_RESET)"; \
		echo "Check logs: tail -f $(LOGS_DIR)/frontend.log"; \
		exit 1; \
	fi

## 停止所有服务
stop:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Stopping services...$(COLOR_RESET)"
	@$(MAKE) stop-frontend
	@$(MAKE) stop-backend
	@echo "$(COLOR_GREEN)✓ All services stopped$(COLOR_RESET)"

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
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Service Status:$(COLOR_RESET)"
	@if [ -f $(RUN_DIR)/backend.pid ]; then \
		BACKEND_PID=$$(cat $(RUN_DIR)/backend.pid 2>/dev/null); \
		if [ -n "$$BACKEND_PID" ] && ps -p $$BACKEND_PID > /dev/null 2>&1; then \
			echo "$(COLOR_GREEN)✓ Backend: Running (PID: $$BACKEND_PID)$(COLOR_RESET)"; \
		else \
			echo "$(COLOR_RED)✗ Backend: Not running$(COLOR_RESET)"; \
		fi; \
	else \
		echo "$(COLOR_RED)✗ Backend: Not running$(COLOR_RESET)"; \
	fi
	@if [ -f $(RUN_DIR)/frontend.pid ]; then \
		FRONTEND_PID=$$(cat $(RUN_DIR)/frontend.pid 2>/dev/null); \
		if [ -n "$$FRONTEND_PID" ] && ps -p $$FRONTEND_PID > /dev/null 2>&1; then \
			echo "$(COLOR_GREEN)✓ Frontend: Running (PID: $$FRONTEND_PID)$(COLOR_RESET)"; \
		else \
			echo "$(COLOR_RED)✗ Frontend: Not running$(COLOR_RESET)"; \
		fi; \
	else \
		echo "$(COLOR_RED)✗ Frontend: Not running$(COLOR_RESET)"; \
	fi

# ============================================================================
# 开发运行（前台运行）
# ============================================================================

## 运行前后端（前台运行，推荐开发使用）
run:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Starting services in foreground...$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)Backend:$(COLOR_RESET) http://localhost:$(BACKEND_PORT)"
	@echo "$(COLOR_BOLD)Frontend:$(COLOR_RESET) http://localhost:$(FRONTEND_PORT)"
	@echo "Press Ctrl+C to stop"
	@trap 'kill 0' EXIT; \
	(cd frontend && npm run dev &) && \
	cd backend && $(GO_CMD) run ./cmd/main.go

## 仅运行后端（前台运行）
run-backend:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Running backend service...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Note: For local development, set LINGPROXY_STORAGE_GORM_DSN=root:@tcp(localhost:2881)/lingproxy?charset=utf8mb4&parseTime=True&loc=Local$(COLOR_RESET)"
	@cd backend && $(GO_CMD) run ./cmd/main.go

## 仅运行前端（前台运行）
run-frontend:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Running frontend development server...$(COLOR_RESET)"
	@cd frontend && npm run dev

# ============================================================================
# 日志管理
# ============================================================================

## 查看后端日志
logs-backend:
	@tail -f $(LOGS_DIR)/backend.log 2>/dev/null || echo "$(COLOR_RED)No backend log file found$(COLOR_RESET)"

## 查看前端日志
logs-frontend:
	@tail -f $(LOGS_DIR)/frontend.log 2>/dev/null || echo "$(COLOR_RED)No frontend log file found$(COLOR_RESET)"

## 查看所有日志
logs:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Backend logs:$(COLOR_RESET)"
	@tail -20 $(LOGS_DIR)/backend.log 2>/dev/null || echo "No backend log"
	@echo "\n$(COLOR_BOLD)$(COLOR_BLUE)Frontend logs:$(COLOR_RESET)"
	@tail -20 $(LOGS_DIR)/frontend.log 2>/dev/null || echo "No frontend log"

# ============================================================================
# Docker 相关
# ============================================================================

## 构建 Docker 镜像（使用 Docker Compose）
docker-build:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Building Docker images with Docker Compose...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml build
	@echo "$(COLOR_GREEN)✓ Docker images built successfully$(COLOR_RESET)"

## 运行 Docker 容器（使用 Docker Compose）
docker-run: docker-compose-up
	@echo "$(COLOR_GREEN)✓ Services are running$(COLOR_RESET)"

## 检查 Docker Compose 配置文件
docker-compose-check:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Checking Docker Compose configuration...$(COLOR_RESET)"
	@if [ ! -f "backend/configs/config.yaml" ]; then \
		echo "$(COLOR_YELLOW)⚠ Config file not found, creating from example...$(COLOR_RESET)"; \
		cp backend/configs/config.yaml.example backend/configs/config.yaml; \
		echo "$(COLOR_YELLOW)⚠ Please edit backend/configs/config.yaml and configure SeekDB connection$(COLOR_RESET)"; \
	fi
	@if [ ! -d "logs" ]; then mkdir -p logs; fi
	@if [ ! -d "run" ]; then mkdir -p run; fi
	@echo "$(COLOR_GREEN)✓ Configuration check completed$(COLOR_RESET)"

## 创建 SeekDB 数据库
docker-compose-init-db:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Initializing SeekDB database...$(COLOR_RESET)"
	@echo "Waiting for SeekDB to be ready..."
	@sleep 10
	@docker exec seekdb mysql -h127.0.0.1 -uroot -P2881 -e "CREATE DATABASE IF NOT EXISTS lingproxy CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null && \
		echo "$(COLOR_GREEN)✓ Database 'lingproxy' created successfully$(COLOR_RESET)" || \
		echo "$(COLOR_YELLOW)⚠ Database creation failed or already exists$(COLOR_RESET)"

## 使用 Docker Compose 启动（完整流程）
docker-compose-up: docker-compose-check
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Starting services with Docker Compose...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml up -d --build
	@echo "$(COLOR_GREEN)✓ Services started$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Waiting for SeekDB to be ready...$(COLOR_RESET)"
	@sleep 15
	@$(MAKE) docker-compose-init-db
	@echo "$(COLOR_GREEN)✓ Services are running$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)Access URLs:$(COLOR_RESET)"
	@echo "  Backend API: $(COLOR_GREEN)http://localhost:8080/api/v1$(COLOR_RESET)"
	@echo "  Health Check: $(COLOR_GREEN)http://localhost:8080/api/v1/health$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Note: Frontend should be run separately with 'cd frontend && npm run dev'$(COLOR_RESET)"

## 停止 Docker Compose 服务
docker-compose-down:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Stopping Docker Compose services...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml down
	@echo "$(COLOR_GREEN)✓ Services stopped$(COLOR_RESET)"

## 查看 Docker Compose 日志
docker-compose-logs:
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml logs -f

## 查看 Docker Compose 服务状态
docker-compose-ps:
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml ps

## 重启 Docker Compose 服务
docker-compose-restart:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Restarting Docker Compose services...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -f docker/docker-compose.yml restart
	@echo "$(COLOR_GREEN)✓ Services restarted$(COLOR_RESET)"

## 重新构建并启动 Docker Compose 服务（停止 → 构建 → 启动）
docker-rebuild: docker-compose-down
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Rebuilding and starting Docker Compose services...$(COLOR_RESET)"
	@$(MAKE) docker-compose-up
	@echo "$(COLOR_GREEN)✓ Services rebuilt and started$(COLOR_RESET)"

# ============================================================================
# 文档生成
# ============================================================================

## 生成 API 文档（Swagger）
docs:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Generating API documentation...$(COLOR_RESET)"
	@if command -v swag >/dev/null 2>&1; then \
		cd backend && swag init -g cmd/main.go -o swagger; \
		echo "$(COLOR_GREEN)✓ API documentation generated$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ swag not found$(COLOR_RESET)"; \
		echo "Install: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# ============================================================================
# 官网构建相关
# ============================================================================

## 安装官网依赖
website-install:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Installing website dependencies...$(COLOR_RESET)"
	@cd website && npm install
	@echo "$(COLOR_GREEN)✓ Dependencies installed$(COLOR_RESET)"

## 启动官网开发服务器
website-start:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Starting website development server...$(COLOR_RESET)"
	@cd website && npm start

## 构建官网
website-build:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Building website...$(COLOR_RESET)"
	@cd website && npm run build
	@echo "$(COLOR_GREEN)✓ Website built successfully. Files are in website/build/$(COLOR_RESET)"

## 预览构建后的官网
website-serve:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Serving built website...$(COLOR_RESET)"
	@cd website && npm run serve

## 清理官网构建缓存
website-clean:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Cleaning website cache...$(COLOR_RESET)"
	@cd website && npm run clear
	@echo "$(COLOR_GREEN)✓ Cache cleaned$(COLOR_RESET)"

## 部署官网到 GitHub Pages
website-deploy:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Deploying website to GitHub Pages...$(COLOR_RESET)"
	@if [ -z "$(GIT_USER)" ]; then \
		echo "$(COLOR_RED)✗ Error: GIT_USER environment variable is required$(COLOR_RESET)"; \
		echo "Example: $(COLOR_YELLOW)GIT_USER=your-username make website-deploy$(COLOR_RESET)"; \
		exit 1; \
	fi
	@cd website && GIT_USER=$(GIT_USER) USE_SSH=true npm run deploy
	@echo "$(COLOR_GREEN)✓ Website deployed successfully$(COLOR_RESET)"

# ============================================================================
# 发布相关
# ============================================================================

## 创建发布版本
release: clean deps test build-all
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Creating release...$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)
	@tar -czf $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz \
		-C $(BIN_DIR) $(BINARY_NAME) 2>/dev/null || true
	@if [ -d "$(FRONTEND_DIST)" ]; then \
		tar -czf $(DIST_DIR)/frontend-$(VERSION).tar.gz -C frontend dist; \
	fi
	@echo "$(COLOR_GREEN)✓ Release created in $(DIST_DIR)/$(COLOR_RESET)"

# ============================================================================
# 工具和实用功能
# ============================================================================

## 检查开发环境
check-env:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Checking development environment...$(COLOR_RESET)"
	@echo "Go version: $$($(GO_CMD) version 2>/dev/null || echo 'Not installed')"
	@echo "Node version: $$(node --version 2>/dev/null || echo 'Not installed')"
	@echo "NPM version: $$(npm --version 2>/dev/null || echo 'Not installed')"
	@echo "Docker version: $$(docker --version 2>/dev/null || echo 'Not installed')"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint: ✓ Installed"; \
	else \
		echo "golangci-lint: ✗ Not installed"; \
	fi
	@if command -v swag >/dev/null 2>&1; then \
		echo "swag: ✓ Installed"; \
	else \
		echo "swag: ✗ Not installed"; \
	fi

## 初始化开发环境
init: deps
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)Initializing development environment...$(COLOR_RESET)"
	@mkdir -p $(LOGS_DIR) $(RUN_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	@if [ ! -f "backend/configs/config.yaml" ]; then \
		echo "$(COLOR_YELLOW)⚠ Creating config.yaml from example...$(COLOR_RESET)"; \
		cp backend/configs/config.yaml.example backend/configs/config.yaml; \
	fi
	@echo "$(COLOR_GREEN)✓ Development environment initialized$(COLOR_RESET)"

# ============================================================================
# 帮助信息
# ============================================================================

## 显示帮助信息
help:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)LingProxy Makefile Commands$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)构建相关:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make all$(COLOR_RESET)              - 完整构建流程（清理、依赖、格式化、检查、测试、构建）"
	@echo "  $(COLOR_GREEN)make build$(COLOR_RESET)            - 构建后端应用"
	@echo "  $(COLOR_GREEN)make build-frontend$(COLOR_RESET)   - 构建前端应用"
	@echo "  $(COLOR_GREEN)make build-all$(COLOR_RESET)       - 构建前后端"
	@echo "  $(COLOR_GREEN)make build-linux$(COLOR_RESET)     - 构建 Linux 版本"
	@echo "  $(COLOR_GREEN)make build-windows$(COLOR_RESET)   - 构建 Windows 版本"
	@echo "  $(COLOR_GREEN)make build-mac$(COLOR_RESET)        - 构建 macOS 版本"
	@echo "  $(COLOR_GREEN)make build-arm$(COLOR_RESET)       - 构建 ARM 版本"
	@echo "  $(COLOR_GREEN)make clean$(COLOR_RESET)            - 清理构建产物"
	@echo ""
	@echo "$(COLOR_BOLD)测试相关:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make test$(COLOR_RESET)             - 运行所有测试"
	@echo "  $(COLOR_GREEN)make test-unit$(COLOR_RESET)        - 运行单元测试"
	@echo "  $(COLOR_GREEN)make test-integration$(COLOR_RESET) - 运行集成测试"
	@echo "  $(COLOR_GREEN)make coverage$(COLOR_RESET)         - 生成代码覆盖率报告"
	@echo "  $(COLOR_GREEN)make bench$(COLOR_RESET)            - 运行性能基准测试"
	@echo ""
	@echo "$(COLOR_BOLD)代码质量:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make fmt$(COLOR_RESET)              - 格式化代码"
	@echo "  $(COLOR_GREEN)make lint$(COLOR_RESET)             - 运行代码检查（需要 golangci-lint）"
	@echo "  $(COLOR_GREEN)make vet$(COLOR_RESET)              - 运行静态分析"
	@echo "  $(COLOR_GREEN)make security$(COLOR_RESET)         - 运行安全扫描（需要 gosec）"
	@echo ""
	@echo "$(COLOR_BOLD)依赖管理:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make deps$(COLOR_RESET)             - 下载并整理依赖"
	@echo "  $(COLOR_GREEN)make deps-update$(COLOR_RESET)      - 更新依赖到最新版本"
	@echo ""
	@echo "$(COLOR_BOLD)服务管理（后台）:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make start$(COLOR_RESET)            - 启动所有服务（后台运行）"
	@echo "  $(COLOR_GREEN)make start-backend$(COLOR_RESET)   - 启动后端服务（后台运行）"
	@echo "  $(COLOR_GREEN)make start-frontend$(COLOR_RESET)  - 启动前端服务（后台运行）"
	@echo "  $(COLOR_GREEN)make stop$(COLOR_RESET)             - 停止所有服务"
	@echo "  $(COLOR_GREEN)make stop-backend$(COLOR_RESET)   - 停止后端服务"
	@echo "  $(COLOR_GREEN)make stop-frontend$(COLOR_RESET)   - 停止前端服务"
	@echo "  $(COLOR_GREEN)make restart$(COLOR_RESET)         - 重启所有服务"
	@echo "  $(COLOR_GREEN)make status$(COLOR_RESET)           - 查看服务状态"
	@echo ""
	@echo "$(COLOR_BOLD)开发运行（前台）:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make run$(COLOR_RESET)              - 运行前后端（前台运行，推荐开发使用）"
	@echo "  $(COLOR_GREEN)make run-backend$(COLOR_RESET)     - 仅运行后端（前台运行）"
	@echo "  $(COLOR_GREEN)make run-frontend$(COLOR_RESET)    - 仅运行前端（前台运行）"
	@echo ""
	@echo "$(COLOR_BOLD)日志管理:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make logs$(COLOR_RESET)            - 查看所有日志（最后20行）"
	@echo "  $(COLOR_GREEN)make logs-backend$(COLOR_RESET)   - 实时查看后端日志"
	@echo "  $(COLOR_GREEN)make logs-frontend$(COLOR_RESET)  - 实时查看前端日志"
	@echo ""
	@echo "$(COLOR_BOLD)Docker 相关:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make docker-build$(COLOR_RESET)         - 构建 Docker 镜像"
	@echo "  $(COLOR_GREEN)make docker-run$(COLOR_RESET)           - 运行 Docker 容器"
	@echo "  $(COLOR_GREEN)make docker-compose-up$(COLOR_RESET)    - 启动服务（包含检查配置、创建数据库）"
	@echo "  $(COLOR_GREEN)make docker-compose-down$(COLOR_RESET)  - 停止 Docker Compose 服务"
	@echo "  $(COLOR_GREEN)make docker-compose-logs$(COLOR_RESET)  - 查看 Docker Compose 日志"
	@echo "  $(COLOR_GREEN)make docker-compose-ps$(COLOR_RESET)    - 查看服务状态"
	@echo "  $(COLOR_GREEN)make docker-compose-restart$(COLOR_RESET) - 重启服务（不重新构建）"
	@echo "  $(COLOR_GREEN)make docker-rebuild$(COLOR_RESET)       - 重新构建并启动服务（停止 → 构建 → 启动）"
	@echo ""
	@echo "$(COLOR_BOLD)文档和工具:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make docs$(COLOR_RESET)            - 生成 API 文档（Swagger）"
	@echo "  $(COLOR_GREEN)make website-install$(COLOR_RESET) - 安装官网依赖"
	@echo "  $(COLOR_GREEN)make website-start$(COLOR_RESET)   - 启动官网开发服务器"
	@echo "  $(COLOR_GREEN)make website-build$(COLOR_RESET)   - 构建官网"
	@echo "  $(COLOR_GREEN)make website-serve$(COLOR_RESET)  - 预览构建后的官网"
	@echo "  $(COLOR_GREEN)make website-deploy$(COLOR_RESET) - 部署官网到 GitHub Pages"
	@echo "  $(COLOR_GREEN)make release$(COLOR_RESET)         - 创建发布版本"
	@echo "  $(COLOR_GREEN)make check-env$(COLOR_RESET)       - 检查开发环境"
	@echo "  $(COLOR_GREEN)make init$(COLOR_RESET)            - 初始化开发环境"
	@echo "  $(COLOR_GREEN)make help$(COLOR_RESET)            - 显示帮助信息"
	@echo ""
	@echo "$(COLOR_BOLD)环境变量:$(COLOR_RESET)"
	@echo "  BACKEND_PORT=$(BACKEND_PORT)  - 后端服务端口"
	@echo "  FRONTEND_PORT=$(FRONTEND_PORT) - 前端服务端口"
	@echo "  GOOS=$(GOOS)                  - Go 目标操作系统"
	@echo "  GOARCH=$(GOARCH)              - Go 目标架构"
	@echo "  VERSION=$(VERSION)            - 版本号"
