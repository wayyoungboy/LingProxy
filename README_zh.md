# LingProxy - AI API网关

<p align="center">
  <img src="web/static/images/logo/lingproxy-logo.svg" alt="LingProxy Logo">
</p>


LingProxy 是一个高性能的AI API网关，专为管理和代理各种AI服务提供商的API调用而设计。它提供了OpenAI兼容接口、负载均衡、熔断保护等功能。

## 功能特性

### 🚀 核心功能
- **统一API接口**: 支持OpenAI兼容API
- **智能负载均衡**: 支持轮询负载均衡策略
- **熔断保护**: 自动检测服务故障并触发熔断，防止级联故障
- **日志追踪**: 完整的请求链路追踪

### 🔐 安全与认证
- **JWT认证**: 基于JWT的用户认证
- **API密钥管理**: 安全的API密钥存储和管理
- **CORS支持**: 跨域资源共享配置
- **限流保护**: 基于IP和用户的请求限流

### �️ 管理功能
- **用户管理**: 多用户支持和权限管理
- **LLM资源管理**: 支持多个AI服务提供商配置

## 快速开始

### 环境要求
- Go 1.21 或更高版本
- SQLite (用于数据存储)

### 安装与运行

1. **克隆项目**
```bash
git clone https://github.com/wayyoungboy/lingproxy.git
cd lingproxy
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置文件**
复制并编辑配置文件：
```bash
cp configs/config.yaml.example configs/config.yaml
# 编辑 configs/config.yaml 根据需要进行配置
```

4. **构建项目**
```bash
go build -o lingproxy ./cmd/main.go
```

5. **运行服务**
```bash
./lingproxy
```

服务将在 `http://localhost:8080` 启动

### Docker 运行

```bash
# 构建镜像
docker build -t lingproxy .

# 运行容器
docker run -p 8080:8080 -v $(pwd)/configs:/app/configs lingproxy
```

## API 使用指南

### 1. 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 2. 获取API密钥
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 3. 代理AI请求
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

## 配置说明

### 主要配置项

#### 应用配置
```yaml
app:
  name: "LingProxy"
  version: "1.0.0"
  environment: "development"  # development, staging, production
  port: 8080
  host: "0.0.0.0"
```

#### 存储配置
```yaml
storage:
  type: "gorm"
  gorm:
    driver: "sqlite"
    dsn: "lingproxy.db"
```

#### 安全配置
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

## 监控与运维

### 日志系统

#### 日志级别
- **DEBUG**: 详细的调试信息，仅用于开发环境
- **INFO**: 系统运行的一般信息
- **WARN**: 需要关注的警告信息
- **ERROR**: 需要立即处理的错误信息
- **FATAL**: 导致系统 shutdown 的严重错误

#### 日志配置
```yaml
log:
  level: "info"  # debug, info, warn, error, fatal
  format: "json"  # text, json
  output: "stdout"
```

#### 日志查看
```bash
# 查看实时日志
# 日志默认输出到 stdout
```

## 开发指南

### 项目结构
```
lingproxy/
├── cmd/                    # 应用入口
├── configs/               # 配置文件
├── docs/                  # API文档
├── internal/              # 内部包
│   ├── cache/             # 缓存实现
│   ├── client/            # AI服务客户端
│   │   ├── embedding/     # 嵌入客户端
│   │   └── openai/        # OpenAI客户端
│   ├── config/            # 配置管理
│   ├── handler/           # HTTP处理器
│   ├── middleware/        # HTTP中间件
│   ├── models/            # 数据模型
│   ├── pkg/               # 内部包
│   │   ├── balancer/      # 负载均衡
│   │   ├── circuitbreaker/ # 熔断器
│   │   └── monitor/       # 监控
│   ├── router/            # 路由
│   ├── service/           # 业务逻辑
│   └── storage/           # 存储实现
├── pkg/                   # 公共包
│   └── logger/            # 日志
├── scripts/               # 脚本文件
├── web/                   # Web界面
│   ├── static/            # 静态文件
│   └── templates/         # HTML模板
└── docker-compose.yml     # Docker配置
```

### 添加新的AI提供商

1. **创建提供商配置**
在 `internal/models/provider.go` 中添加新的提供商类型

2. **实现负载均衡策略**
在 `internal/pkg/balancer/` 中实现新的负载均衡算法

3. **添加监控指标**
在 `internal/pkg/monitor/` 中添加新的监控指标

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./internal/pkg/balancer

# 运行带覆盖率测试
go test -cover ./...
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 支持与联系

- **Issues**: [GitHub Issues](https://github.com/wayyoungboy/lingproxy/issues)
- **Discussions**: [GitHub Discussions](https://github.com/wayyoungboy/lingproxy/discussions)
- **Email**: support@lingproxy.com

## 更新日志

### v1.0.0 (2026-01-30)
- 初始版本发布
- 支持OpenAI兼容API
- 实现轮询负载均衡和熔断保护
- 添加用户管理和LLM资源管理
- 提供完整的REST API接口
- 实现基于SQLite的数据存储
- 添加支持多级别日志的日志系统
- 创建基于Web的管理界面

## 语言

- [English](README.md)
- [中文](README_zh.md) (current)