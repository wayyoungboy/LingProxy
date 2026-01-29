# LingProxy - AI API网关

LingProxy 是一个高性能的AI API网关，专为管理和代理各种AI服务提供商的API调用而设计。它提供了统一的接口、负载均衡、熔断保护、配额管理、监控告警等功能。

## 功能特性

### 🚀 核心功能
- **统一API接口**: 支持OpenAI、Claude、Gemini等多种AI服务提供商
- **智能负载均衡**: 支持轮询、权重、最少连接等多种负载均衡策略
- **熔断保护**: 自动检测服务故障并触发熔断，防止级联故障
- **配额管理**: 支持RPM、TPM、每日Token和请求数限制
- **缓存机制**: 智能缓存响应，减少重复请求
- **请求/响应转换**: 支持多种格式的请求和响应转换

### � 监控与告警
- **实时监控**: 基于Prometheus的指标收集
- **可视化面板**: 集成Grafana监控面板
- **智能告警**: 支持邮件、Webhook等多种告警方式
- **日志追踪**: 完整的请求链路追踪

### 🔐 安全与认证
- **JWT认证**: 基于JWT的用户认证
- **API密钥管理**: 安全的API密钥存储和管理
- **CORS支持**: 跨域资源共享配置
- **限流保护**: 基于IP和用户的请求限流

### �️ 管理功能
- **用户管理**: 多用户支持和权限管理
- **提供商管理**: 支持多个AI服务提供商配置
- **端点管理**: 灵活的API端点配置
- **密钥轮换**: 自动化的API密钥轮换

## 快速开始

### 环境要求
- Go 1.21 或更高版本
- Redis (可选，用于配额管理)
- SQLite 或 MySQL (用于数据存储)

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
支持 SQLite 和 MySQL：
```yaml
storage:
  type: "sqlite"  # sqlite 或 mysql
  sqlite:
    path: "./data/lingproxy.db"
  mysql:
    host: "localhost"
    port: 3306
    user: "lingproxy"
    password: "password"
    database: "lingproxy"
```

#### Redis 配置（用于配额管理）
```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

#### 安全配置
```yaml
security:
  jwt:
    secret: "your-jwt-secret-key"
    expire_hours: "24h"
  rate_limit:
    enabled: true
    requests_per_minute: 1000
```

## 监控与运维

### 访问监控面板
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

### 关键指标
- 请求总数和成功率
- 响应时间和延迟分布
- 各提供商的使用统计
- 配额使用情况
- 错误率和熔断状态

### 日志查看
```bash
# 查看实时日志
tail -f logs/lingproxy.log

# 查看特定级别日志
grep "ERROR" logs/lingproxy.log
```

## 开发指南

### 项目结构
```
lingproxy/
├── cmd/                    # 应用入口
├── configs/               # 配置文件
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── models/           # 数据模型
│   ├── pkg/              # 内部包
│   │   ├── balancer/     # 负载均衡
│   │   ├── circuitbreaker/ # 熔断器
│   │   └── monitor/      # 监控
│   ├── router/           # 路由
│   └── storage/          # 存储
├── pkg/                  # 公共包
│   └── logger/           # 日志
└── scripts/              # 脚本文件
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

### v1.0.0 (2024-01-29)
- 初始版本发布
- 支持OpenAI、Claude、Gemini API代理
- 实现负载均衡和熔断保护
- 添加配额管理和监控功能
- 提供完整的REST API接口