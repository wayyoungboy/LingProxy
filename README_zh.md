<div align="center">

<img src="./assets/lingproxy-logo-full.svg" alt="LingProxy Logo" width="300">

# LingProxy - AI API网关

LingProxy 是一个高性能的AI API网关，专为管理和代理各种AI服务提供商的API调用而设计。它提供了OpenAI兼容接口、负载均衡等功能。

</div>

## 功能特性

### 🚀 核心功能
- **统一API接口**: 支持OpenAI兼容API，可无缝对接各种AI服务
- **流式响应支持**: 完整支持 Server-Sent Events (SSE) 流式响应，适用于聊天补全接口
- **智能负载均衡**: 轮询负载均衡策略，自动分配请求到多个资源
- **自动重试机制**: 可配置的自动重试功能，支持指数退避，自动重试网络错误、超时和5xx服务器错误
- **请求日志**: 完整的请求链路追踪和日志记录

### 🔐 安全与认证
- **灵活认证机制**: 支持全局认证开关，可配置是否启用认证
- **管理员登录**: 用户名/密码登录，支持密码哈希存储
- **API Key管理**: 请求端API Key管理，支持策略关联和API密钥认证
- **CORS支持**: 灵活的跨域资源共享配置
- **安全存储**: API密钥和密码加密存储

### 📊 管理功能
- **管理后台**: 基于 Vue 3 + Element Plus 的现代化 Web 管理界面
- **国际化支持**: 前端界面完整支持中英文切换
- **管理员管理**: 单管理员模式，支持密码和API密钥管理
- **API Key管理**: 创建和管理请求端API Key，支持策略绑定，支持API Key复制功能
- **策略管理**: 内置多种路由策略模板（随机、轮询、加权、模型匹配、正则匹配、优先级、故障转移），支持自定义策略实例，随机选择策略支持LLM资源池配置
- **LLM资源管理**: 支持基于驱动的AI服务资源配置（目前支持OpenAI驱动），支持模型类别（对话、生图、嵌入、重排序、语音、视频），支持Excel模板或JSON格式批量导入/导出，包含资源测试功能以验证连通性
- **模型管理**: 灵活的模型配置，支持定价、使用限制等参数
- **请求管理**: 完整的请求记录和追踪，支持请求详情查看和导出
- **用量统计**: 按LLM资源分组的详细用量统计，包括Token使用量、请求数、成功率、平均Token/请求等指标，支持时间范围和资源名称筛选
- **系统设置**: 动态配置管理，包括基础设置、缓存、限流、安全、日志、负载均衡、Provider重试等配置
- **系统监控**: 实时系统信息查看（CPU、内存、运行时间等）
- **日志管理**: 查看和管理系统日志，支持过滤和搜索功能

### 🏗️ 架构设计
- **前后端分离**: 现代架构，使用Vue 3 + Element Plus前端和Go后端
- **国际化支持**: 完整的前端i18n支持，使用vue-i18n，支持中文和英文
- **精简模型**: 去除冗余功能，核心代码简洁高效
- **双重存储**: 支持内存存储（开发调试）和SQLite存储（生产环境）
- **模块化设计**: 清晰的层次结构，易于扩展和维护
- **RESTful API**: 完整的REST API接口，便于集成
- **客户端库**: `clients/` 目录提供标准客户端实现（Python、JavaScript、Go）

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
# ⚠️ 重要：启动前请修改 config.yaml 中的管理员密码！
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

前端将在 `http://localhost:3000` 可用

### Docker 部署

项目采用**前后端分离**架构。Docker 部署仅包含后端服务。

#### 使用 Docker Compose（推荐）

1. **启动后端服务**
```bash
# 构建并启动后端 + 数据库（在项目根目录执行）
docker-compose -f docker/docker-compose.yml up -d

# 查看日志
docker-compose -f docker/docker-compose.yml logs -f lingproxy-backend

# 停止服务
docker-compose -f docker/docker-compose.yml down
```

**后端 API**: http://localhost:8080/api/v1

2. **单独运行前端**（开发时）
```bash
cd frontend
npm install
npm run dev
```

**前端**: http://localhost:3000（API 请求会自动代理到后端）

**注意**: 
- Docker 部署使用 `docker/backend.Dockerfile`（仅后端）
- 数据库会在后端启动时自动创建
- 使用 `config.yaml.docker` 作为 Docker 专用配置
- 详细 Docker 部署指南请查看[快速开始指南](docs/zh/02-quick-start.md#docker-部署)

## API 使用指南

### 1. 管理员登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "YOUR_PASSWORD"
  }'
```

响应示例：
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

### 2. 创建请求端API Key
```bash
curl -X POST http://localhost:8080/api/v1/api-keys \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "我的 API Key",
    "status": "active"
  }'
```

响应示例：
```json
{
  "data": {
    "id": "...",
    "name": "我的 API Key",
    "token": "ling-xxxxxxxxxxxxx",
    "status": "active"
  }
}
```

### 3. 创建路由策略（可选）
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

### 4. 为API Key绑定策略（可选）
```bash
curl -X PUT http://localhost:8080/api/v1/api-keys/API_KEY_ID/policy \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "policy_id": "POLICY_ID"
  }'
```

**注意**: 将 `API_KEY_ID` 替换为步骤2中创建的实际API Key ID。

### 5. 代理AI请求
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

### API 端点参考

#### 认证与管理员
- `POST /api/v1/auth/login` - 管理员登录（用户名/密码）
- `GET /api/v1/admin/info` - 获取管理员信息
- `PUT /api/v1/admin/api-key` - 重置管理员API密钥

#### API Key管理
- `GET /api/v1/api-keys` - 获取API Key列表
- `GET /api/v1/api-keys/:id` - 获取API Key详情
- `POST /api/v1/api-keys` - 创建API Key
- `PUT /api/v1/api-keys/:id` - 更新API Key
- `DELETE /api/v1/api-keys/:id` - 删除API Key
- `POST /api/v1/api-keys/:id/reset` - 重置API Key
- `PUT /api/v1/api-keys/:id/policy` - 为API Key绑定策略
- `DELETE /api/v1/api-keys/:id/policy` - 移除API Key的策略绑定

**注意**: 旧的 `/api/v1/tokens` 端点仍然可用以保持向后兼容，但已弃用。

#### 策略管理
- `GET /api/v1/policy-templates` - 获取策略模板列表
- `GET /api/v1/policy-templates/:id` - 获取策略模板详情
- `GET /api/v1/policies` - 获取策略列表
- `GET /api/v1/policies/:id` - 获取策略详情
- `POST /api/v1/policies` - 创建策略
- `PUT /api/v1/policies/:id` - 更新策略
- `DELETE /api/v1/policies/:id` - 删除策略

#### LLM资源管理
- `GET /api/v1/llm-resources` - 获取LLM资源列表（支持搜索过滤）
- `POST /api/v1/llm-resources` - 创建LLM资源
- `GET /api/v1/llm-resources/:id` - 获取LLM资源详情
- `PUT /api/v1/llm-resources/:id` - 更新LLM资源
- `DELETE /api/v1/llm-resources/:id` - 删除LLM资源
- `POST /api/v1/llm-resources/:id/test` - 测试LLM资源连通性
- `POST /api/v1/llm-resources/import` - 批量导入LLM资源（Excel或JSON格式）
- `GET /api/v1/llm-resources/import/template` - 下载Excel导入模板

**批量导入说明**:
- 支持通过Excel文件或JSON格式批量导入LLM资源
- Excel模板包含字段：资源名称、模型类别、驱动、模型标识、基础URL、API密钥、状态
- JSON导入接受资源对象数组，字段与Excel相同
- 驱动字段目前仅支持"openai"，如为空或非openai值将自动设置为"openai"
- 导入结果会返回成功、失败、重复的数量及详细的错误和重复信息
- 重复检测：type、model、base_url、api_key 都相同的资源会被识别为重复
- 自动去空格：导入时会自动去除所有字段的前后空格

**资源测试功能**:
- LLM资源管理界面提供测试按钮
- 只有状态为 `active` 的资源才能被测试
- 支持测试 `chat` 和 `embedding` 类型的资源
- 返回详细的测试结果，包括响应时间、模型信息、Token使用情况和响应内容
- 测试超时时间：30秒

**搜索功能**:
- 前端支持对资源名称、基础URL和模型标识进行模糊搜索
- 搜索不区分大小写，支持部分匹配

#### 模型管理
- `GET /api/v1/models` - 获取模型列表
- `POST /api/v1/models` - 创建模型
- `GET /api/v1/models/:id` - 获取模型详情
- `PUT /api/v1/models/:id` - 更新模型
- `DELETE /api/v1/models/:id` - 删除模型
- `GET /api/v1/models/types` - 获取模型类型列表
- `GET /api/v1/models/:id/pricing` - 获取模型定价信息
- `GET /api/v1/llm-resources/:id/models` - 获取指定LLM资源下的模型列表

#### 请求日志
- `GET /api/v1/requests` - 获取请求日志列表
- `GET /api/v1/requests/:id` - 获取请求详情
- `POST /api/v1/requests` - 创建请求记录

#### 系统设置与监控
- `GET /api/v1/settings` - 获取系统设置
- `PUT /api/v1/settings` - 更新系统设置
- `GET /api/v1/system/info` - 获取系统信息（CPU、内存、运行时间等）

#### 统计信息
- `GET /api/v1/stats/system` - 获取系统统计信息（总请求数、总用户数、总LLM资源数、成功率、平均响应时间）
- `GET /api/v1/stats/llm-resources/usage` - 获取LLM资源使用统计（按资源分组，包含Token使用量、请求数、成功率等）
- `GET /api/v1/stats/llm-resources/:id` - 获取单个LLM资源统计信息
- `GET /api/v1/stats/users/:id` - 获取用户统计信息

#### OpenAI 兼容 API
- `GET /llm/v1/models` - 获取所有可用模型列表
- `GET /llm/v1/models/:model` - 获取模型信息
- `POST /llm/v1/chat/completions` - 创建聊天补全（支持流式响应，使用 `stream: true`）
- `POST /llm/v1/completions` - 创建文本补全

  **curl 示例：**
  ```bash
  # 如果认证已启用（默认情况）
  curl -X GET http://localhost:8080/llm/v1/models \
    -H "Authorization: Bearer YOUR_API_KEY_OR_TOKEN"
  
  # 如果认证已禁用（security.auth.enabled: false）
  curl -X GET http://localhost:8080/llm/v1/models
  ```

- `GET /llm/v1/models/:model` - 获取模型信息

  **curl 示例：**
  ```bash
  # 如果认证已启用
  curl -X GET http://localhost:8080/llm/v1/models/gpt-3.5-turbo \
    -H "Authorization: Bearer YOUR_API_KEY_OR_TOKEN"
  
  # 如果认证已禁用
  curl -X GET http://localhost:8080/llm/v1/models/gpt-3.5-turbo
  ```

- `POST /llm/v1/chat/completions` - 创建聊天补全

  **curl 示例：**
  ```bash
  # 如果认证已启用
  curl -X POST http://localhost:8080/llm/v1/chat/completions \
    -H "Authorization: Bearer YOUR_API_KEY_OR_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "model": "gpt-3.5-turbo",
      "messages": [
        {"role": "user", "content": "Hello!"}
      ]
    }'
  ```

- `POST /llm/v1/completions` - 创建文本补全

  **curl 示例：**
  ```bash
  # 如果认证已启用
  curl -X POST http://localhost:8080/llm/v1/completions \
    -H "Authorization: Bearer YOUR_API_KEY_OR_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "model": "gpt-3.5-turbo",
      "prompt": "Say hello"
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
  auth:
    enabled: true  # 是否启用认证，false时所有API（除登录外）都不需要认证
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

#### 管理员配置
```yaml
admin:
  username: "admin"
  # ⚠️ 请设置强密码！首次启动后建议立即修改
  # password: "YOUR_STRONG_PASSWORD_HERE"
  password: ""  # 留空则不会设置密码，需要通过其他方式设置
  api_key: ""  # 留空则自动生成，首次启动后查看日志获取
  auto_create: true
```

#### 日志配置
```yaml
log:
  level: "info"  # debug, info, warn, error, fatal
  format: "json"  # text, json
  output: "stdout"
```

#### 负载均衡配置
```yaml
load_balancer:
  default_strategy: "round_robin"  # 默认负载均衡策略
```

#### 服务提供商配置
```yaml
provider:
  timeout: "30s"  # 请求超时时间
  max_retries: 3   # 失败请求的最大重试次数（0 = 禁用重试）
  retry_delay: "1s"  # 重试之间的基础延迟时间（实际延迟会指数增长）
  max_idle_conns: 100  # 最大空闲连接数
  max_conns_per_host: 100  # 每个主机最大连接数
  idle_conn_timeout: "90s"  # 空闲连接超时时间
```

**重试机制说明：**
- 自动重试网络错误、超时和5xx服务器错误的请求
- 使用指数退避策略：延迟 = 重试延迟 × 重试次数
- 不会重试4xx客户端错误（429限流除外）、认证错误或上下文取消错误
- 可通过管理界面配置：设置 → Provider设置
- 适用于所有请求类型：聊天补全（流式和非流式）、文本补全、嵌入请求

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
│   │   └── balancer/      # 负载均衡
│   ├── router/            # 路由
│   ├── service/           # 业务逻辑
│   └── storage/           # 存储实现
├── pkg/                   # 公共包
│   └── logger/            # 日志
└── docker-compose.yml     # Docker配置
```

### 数据模型

系统采用精简的存储模型设计，核心模型包括：

```go
// User 用户模型 - 管理员用户
type User struct {
    ID           string     // 用户唯一标识
    Username     string     // 用户名
    PasswordHash string     // 密码哈希值
    APIKey       string     // API密钥
    Role         string     // 角色 (admin)
    Status       string     // 状态 (active, inactive, suspended)
    LastLoginAt  *time.Time // 最后登录时间
    CreatedAt    time.Time  // 创建时间
    UpdatedAt    time.Time  // 更新时间
}

// Token API Key模型 - 请求端API Key管理
type Token struct {
    ID         string     // API Key唯一标识
    Name       string     // API Key名称/描述
    Token      string     // API Key值（以"ling-"开头）
    Prefix     string     // API Key前缀（用于显示）
    Status     string     // 状态 (active/inactive)
    PolicyID   string     // 关联的策略ID（可选）
    LastUsedAt *time.Time // 最后使用时间
    ExpiresAt  *time.Time // 过期时间（可选）
    CreatedAt  time.Time  // 创建时间
    UpdatedAt  time.Time  // 更新时间
}

// PolicyTemplate 策略模板模型 - 内置策略模板
type PolicyTemplate struct {
    ID                string    // 模板唯一标识
    Name              string    // 模板名称
    Type              string    // 类型 (random, round_robin, weighted, model_match, regex_match, priority, failover)
    Description       string    // 描述
    ParametersSchema  string    // 参数JSON Schema
    DefaultParameters string    // 默认参数JSON
    Builtin           bool      // 是否内置
    CreatedAt         time.Time // 创建时间
    UpdatedAt         time.Time // 更新时间
}

// Policy 策略实例模型 - 路由策略配置
type Policy struct {
    ID         string    // 策略唯一标识
    Name       string    // 策略名称
    TemplateID string    // 关联的模板ID
    Type       string    // 类型
    Parameters string    // 参数JSON
    Enabled    bool      // 是否启用
    CreatedAt  time.Time // 创建时间
    UpdatedAt  time.Time // 更新时间
}

// LLMResource LLM资源模型 - AI服务资源配置
type LLMResource struct {
    ID        string    // 资源唯一标识
    Name      string    // 资源名称
    Type      string    // 模型类别 (chat, image, embedding, rerank, audio, video)
    Driver    string    // 驱动 (目前支持: openai)
    Model     string    // 模型标识 (如: gpt-4, gpt-3.5-turbo)
    BaseURL   string    // API基础URL
    APIKey    string    // API密钥
    Status    string    // 状态 (active/inactive)
    CreatedAt time.Time // 创建时间
    UpdatedAt time.Time // 更新时间
}

// Model 模型配置 - AI模型管理
type Model struct {
    ID            string    // 模型唯一标识
    Name          string    // 模型名称
    LLMResourceID string    // 关联的LLM资源
    ModelID       string    // 提供商内部的模型标识
    Type          string    // 模型类型 (chat, completion, embedding, image)
    Category      string    // 模型分类 (gpt, claude, gemini, llama等)
    Version       string    // 模型版本
    Description   string    // 描述
    Capabilities  string    // 模型能力 (JSON字符串)
    Pricing       string    // 定价信息 (JSON字符串)
    Limits        string    // 使用限制 (JSON字符串)
    Parameters    string    // 默认参数 (JSON字符串)
    Features      string    // 功能特性 (JSON字符串)
    Status        string    // 状态 (active, inactive, deprecated)
    Metadata      string    // 扩展元数据 (JSON字符串)
    CreatedAt     time.Time // 创建时间
    UpdatedAt     time.Time // 更新时间
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
    CreatedAt time.Time // 创建时间
}
```

### 存储层设计

存储层采用简洁的接口设计，支持内存存储和GORM存储两种实现：

```go
type Storage interface {
    // 用户管理
    CreateUser(user *User) error
    GetUser(id string) (*User, error)
    GetUserByUsername(username string) (*User, error)
    GetUserByAPIKey(apiKey string) (*User, error)
    UpdateUser(user *User) error
    DeleteUser(id string) error
    ListUsers() ([]*User, error)

    // API Key管理
    CreateToken(token *Token) error
    GetToken(id string) (*Token, error)
    GetTokenByToken(token string) (*Token, error)
    UpdateToken(token *Token) error
    DeleteToken(id string) error
    ListTokens() ([]*Token, error)

    // 策略模板管理
    CreatePolicyTemplate(template *PolicyTemplate) error
    GetPolicyTemplate(id string) (*PolicyTemplate, error)
    GetPolicyTemplateByType(type string) (*PolicyTemplate, error)
    UpdatePolicyTemplate(template *PolicyTemplate) error
    DeletePolicyTemplate(id string) error
    ListPolicyTemplates() ([]*PolicyTemplate, error)

    // 策略管理
    CreatePolicy(policy *Policy) error
    GetPolicy(id string) (*Policy, error)
    UpdatePolicy(policy *Policy) error
    DeletePolicy(id string) error
    ListPolicies() ([]*Policy, error)

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

    // 请求日志
    CreateRequest(request *Request) error
    GetRequest(id string) (*Request, error)
    ListRequests(limit int) ([]*Request, error)
}
```

### 添加新的AI驱动

1. **更新LLM资源模型**
在 `internal/handler/provider.go` 中扩展Driver字段验证逻辑以支持新的驱动类型

2. **实现驱动客户端**
在 `internal/client/` 中创建新驱动的客户端实现

3. **更新负载均衡策略**
如需要，在 `internal/pkg/balancer/` 中实现或更新负载均衡算法

4. **更新前端**
在前端LLM资源管理界面中添加新驱动选项

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

### v1.5.0 (2026-02-08)
- **自动重试功能**: 新增可配置的自动重试机制，支持指数退避策略，自动重试失败的请求
- **Provider配置**: 新增Provider设置（超时时间、最大重试次数、重试延迟），可通过管理界面配置
- **错误分类**: 智能错误分类，区分可重试和不可重试的错误
- **流式重试**: 重试逻辑现在也适用于流式请求（在流建立之前）
- **API Key管理**: 将"Token管理"更名为"API Key管理"，避免与LLM消耗的token混淆
- **文档更新**: 全面更新所有文档（README、配置指南、API参考、架构文档）

### v1.4.0 (2026-02-05)
- **国际化功能**: 前端完整支持中英文切换的i18n功能
- **流式响应支持**: 新增 Server-Sent Events (SSE) 流式响应支持，适用于聊天补全接口
- **策略增强**: 随机选择策略支持LLM资源池配置
- **客户端库**: 新增Python、JavaScript、Go标准客户端实现
- **代码清理**: 删除冗余的 `backend/examples` 目录，统一客户端示例到 `clients/` 目录
- **文档更新**: 全面更新所有语言的文档

### v1.3.0 (2026-02-03)
- **驱动架构**: 将"服务提供商"改为"驱动"概念，目前仅支持OpenAI驱动
- **批量导入导出**: 新增Excel模板下载和批量导入LLM资源功能
- **搜索增强**: 新增资源名称、基础URL和模型标识的模糊搜索支持
- **前端改进**: 修复批量导入后数据不显示问题，优化搜索体验
- **模板管理**: Excel导入模板包含核心字段（名称、类别、驱动、模型标识、基础URL、API密钥、状态）

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