# LingProxy - AI API网关

<p align="center">
  <img src="web/static/images/logo/lingproxy-logo.svg" alt="LingProxy Logo">
</p>


LingProxy 是一个高性能的AI API网关，专为管理和代理各种AI服务提供商的API调用而设计。它提供了OpenAI兼容接口、负载均衡、熔断保护等功能。

## 功能特性

### 🚀 核心功能
- **统一API接口**: 支持OpenAI兼容API，可无缝对接各种AI服务
- **智能负载均衡**: 轮询负载均衡策略，自动分配请求到多个资源
- **熔断保护**: 自动检测服务故障并触发熔断，防止级联故障
- **请求日志**: 完整的请求链路追踪和日志记录

### 🔐 安全与认证
- **API密钥认证**: 基于API密钥的用户认证机制
- **CORS支持**: 灵活的跨域资源共享配置
- **安全存储**: API密钥加密存储

### 📊 管理功能
- **用户管理**: 多用户支持，API密钥管理
- **LLM资源管理**: 支持配置多个AI服务提供商（OpenAI、Anthropic、Google等）
- **模型管理**: 灵活的模型配置，支持定价、使用限制等参数
- **端点管理**: 自定义API端点路由配置

### 🏗️ 架构设计
- **前后端分离**: 现代架构，使用Vue 3 + Element Plus前端和Go后端
- **精简模型**: 去除冗余功能，核心代码简洁高效
- **双重存储**: 支持内存存储（开发调试）和SQLite存储（生产环境）
- **模块化设计**: 清晰的层次结构，易于扩展和维护
- **RESTful API**: 完整的REST API接口，便于集成

## 快速开始

### 环境要求
- **后端**: Go 1.21 或更高版本，SQLite (用于数据存储)
- **前端**: Node.js 18+，npm 或 yarn

### 安装与运行

#### 后端设置

1. **克隆项目**
```bash
git clone https://github.com/wayyoungboy/lingproxy.git
cd lingproxy
```

2. **安装 Go 依赖**
```bash
go mod tidy
```

3. **配置文件**
复制并编辑配置文件：
```bash
cp configs/config.yaml.example configs/config.yaml
# 编辑 configs/config.yaml 根据需要进行配置
```

4. **构建并运行后端**
```bash
go run cmd/main.go
```

后端服务将在 `http://localhost:8080` 启动

#### 前端设置

1. **安装 Node.js 依赖**
```bash
cd frontend
npm install
```

2. **运行前端开发服务器**
```bash
npm run dev
```

前端将在 `http://localhost:3002` 可用

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

### API 端点参考

#### 用户管理
- `GET /api/v1/users` - 获取用户列表
- `POST /api/v1/users` - 创建用户
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户
- `DELETE /api/v1/users/:id` - 删除用户

#### LLM资源管理
- `GET /api/v1/llm-resources` - 获取LLM资源列表
- `POST /api/v1/llm-resources` - 创建LLM资源
- `GET /api/v1/llm-resources/:id` - 获取LLM资源详情
- `PUT /api/v1/llm-resources/:id` - 更新LLM资源
- `DELETE /api/v1/llm-resources/:id` - 删除LLM资源

#### 模型管理
- `GET /api/v1/models` - 获取模型列表
- `POST /api/v1/models` - 创建模型
- `GET /api/v1/models/:id` - 获取模型详情
- `PUT /api/v1/models/:id` - 更新模型
- `DELETE /api/v1/models/:id` - 删除模型

#### 端点管理
- `GET /api/v1/endpoints` - 获取端点列表
- `POST /api/v1/endpoints` - 创建端点
- `GET /api/v1/endpoints/:id` - 获取端点详情
- `PUT /api/v1/endpoints/:id` - 更新端点
- `DELETE /api/v1/endpoints/:id` - 删除端点

#### 请求日志
- `GET /api/v1/requests` - 获取请求日志列表
- `GET /api/v1/requests/:id` - 获取请求详情

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
├── frontend/              # 前端应用
│   ├── public/             # 公共资源
│   ├── src/                # 源代码
│   │   ├── api/            # API客户端
│   │   ├── assets/         # 静态资源
│   │   ├── components/     # Vue组件
│   │   ├── router/         # Vue路由
│   │   ├── views/           # Vue视图
│   │   ├── App.vue         # 根组件
│   │   └── main.js         # 入口文件
│   ├── package.json        # npm配置
│   └── vite.config.js      # Vite配置
├── internal/              # 内部包
│   ├── cache/             # 缓存实现
│   ├── client/            # AI服务客户端
│   │   ├── embedding/     # 嵌入客户端
│   │   └── openai/        # OpenAI客户端
│   ├── config/            # 配置管理
│   ├── handler/           # HTTP处理器
│   ├── middleware/        # HTTP中间件
│   ├── pkg/               # 内部包
│   │   ├── balancer/      # 负载均衡
│   │   └── circuitbreaker/ # 熔断器
│   ├── router/            # 路由
│   ├── service/           # 业务逻辑
│   └── storage/           # 存储实现
├── pkg/                   # 公共包
│   └── logger/            # 日志
├── scripts/               # 脚本文件
├── web/                   # 旧版Web界面（已废弃）
└── docker-compose.yml     # Docker配置
```

### 数据模型

系统采用精简的存储模型设计，核心模型包括：

```go
// User 用户模型 - 管理API用户
type User struct {
    ID        string    // 用户唯一标识
    Username  string    // 用户名
    Email     string    // 邮箱
    APIKey    string    // API密钥
    Status    string    // 状态
}

// LLMResource LLM资源模型 - AI服务提供商配置
type LLMResource struct {
    ID        string    // 资源唯一标识
    Name      string    // 资源名称
    Type      string    // 类型 (openai, anthropic, google等)
    Model     string    // 默认模型
    BaseURL   string    // API基础URL
    APIKey    string    // API密钥
    Status    string    // 状态
}

// Model 模型配置 - AI模型管理
type Model struct {
    ID             string         // 模型唯一标识
    Name           string         // 模型名称
    LLMResourceID  string         // 关联的LLM资源
    ModelID        string         // 提供商内部的模型标识
    Type           string         // 模型类型 (chat, completion, embedding, image)
    Category       string         // 模型分类
    Pricing        ModelPricing   // 定价信息
    Limits         ModelLimits    // 使用限制
    Parameters     ModelParameters // 默认参数
    Status         string         // 状态
}

// Endpoint 端点模型 - API路由配置
type Endpoint struct {
    ID            string    // 端点唯一标识
    LLMResourceID string    // 关联的LLM资源
    Path          string    // API路径
    Method        string    // HTTP方法
    Status        string    // 状态
}

// Request 请求模型 - 请求日志
type Request struct {
    ID        string    // 请求唯一标识
    UserID    string    // 用户ID
    Endpoint  string    // 请求端点
    Method    string    // HTTP方法
    Status    string    // 状态
    Duration  int64     // 耗时(毫秒)
    Tokens    int       // 消耗的token数
}
```

### 存储层设计

存储层采用简洁的接口设计，支持内存存储和GORM存储两种实现：

```go
type Storage interface {
    // 用户管理
    CreateUser(user *User) error
    GetUser(id string) (*User, error)
    GetUserByAPIKey(apiKey string) (*User, error)
    UpdateUser(user *User) error
    DeleteUser(id string) error
    ListUsers() ([]*User, error)

    // LLM资源管理
    CreateLLMResource(resource *LLMResource) error
    GetLLMResource(id string) (*LLMResource, error)
    UpdateLLMResource(resource *LLMResource) error
    DeleteLLMResource(id string) error
    ListLLMResources() ([]*LLMResource, error)

    // 模型管理
    CreateModel(model *Model) error
    GetModel(id string) (*Model, error)
    UpdateModel(model *Model) error
    DeleteModel(id string) error
    ListModels() ([]*Model, error)
    ListModelsByLLMResource(llmResourceID string) ([]*Model, error)

    // 端点管理
    CreateEndpoint(endpoint *Endpoint) error
    GetEndpoint(id string) (*Endpoint, error)
    UpdateEndpoint(endpoint *Endpoint) error
    DeleteEndpoint(id string) error
    ListEndpoints() ([]*Endpoint, error)

    // 请求日志
    CreateRequest(request *Request) error
    GetRequest(id string) (*Request, error)
    ListRequests(limit int) ([]*Request, error)
}
```

### 添加新的AI提供商

1. **创建LLM资源配置**
在 `internal/handler/provider.go` 中添加新的提供商类型处理逻辑

2. **实现负载均衡策略**
在 `internal/pkg/balancer/` 中实现新的负载均衡算法

3. **更新模型支持**
在 `internal/storage/model.go` 中扩展 Model 结构体以支持新的模型特性

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

### v1.2.0 (2026-02-02)
- **前后端分离**: 实现现代 Vue 3 + Element Plus 前端
- **新前端界面**: 使用 Vue 3 Composition API 和 Script Setup 完全重写
- **增强UI**: 响应式设计，使用 Element Plus 组件
- **改进API集成**: 基于 Axios 的 API 客户端，具有完善的错误处理
- **后端API更新**: 添加缺失的端点管理 API
- **Web界面移除**: 移除旧版 Web 界面路由
- **文档更新**: 添加前端开发指南和架构文档

### v1.1.0 (2026-02-01)
- **架构优化**: 精简核心代码，提升代码质量和可维护性
- **模型简化**: 移除未使用的 ModelEndpoint 和 ModelVersion 结构体
- **监控模块优化**: 简化为轻量级配额管理器
- **存储层重构**: 优化存储接口，移除冗余方法
- **依赖修复**: 修复 embedding 客户端的依赖问题
- **文档更新**: 完善开发指南和数据模型文档

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