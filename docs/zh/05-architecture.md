# 架构指南

## 系统架构

LingProxy 采用现代微服务架构，关注点分离清晰：

```
┌─────────────┐
│   前端      │  Vue 3 + Element Plus
│  (端口 3000)│
└──────┬──────┘
       │ HTTP/REST API
┌──────▼──────────────────┐
│      后端 API           │  Go + Gin
│     (端口 8080)         │
├──────────────────────────┤
│  ┌────────────────────┐ │
│  │   HTTP 处理器      │ │  请求处理
│  └──────────┬─────────┘ │
│  ┌──────────▼─────────┐ │
│  │   中间件           │ │  认证、CORS、日志
│  └──────────┬─────────┘ │
│  ┌──────────▼─────────┐ │
│  │   服务层           │ │  业务逻辑
│  └──────────┬─────────┘ │
│  ┌──────────▼─────────┐ │
│  │   存储层           │ │  数据持久化
│  └────────────────────┘ │
└──────────────────────────┘
       │
┌──────▼──────┐
│   数据库    │  SQLite/MySQL/PostgreSQL
└─────────────┘
```

## 后端架构

### 目录结构

```
backend/
├── cmd/
│   └── main.go              # 应用入口
├── configs/
│   └── config.yaml.example  # 配置模板
├── examples/
│   └── llm_demo.go          # 示例代码
├── internal/
│   ├── cache/               # 缓存实现
│   ├── client/              # AI 服务客户端
│   │   ├── embedding/       # 嵌入客户端
│   │   └── openai/          # OpenAI 客户端
│   ├── config/              # 配置管理
│   ├── handler/             # HTTP 处理器
│   ├── middleware/          # HTTP 中间件
│   ├── pkg/                 # 内部工具包
│   │   ├── balancer/        # 负载均衡
│   │   ├── logger/          # 日志工具
│   │   ├── monitor/         # 监控工具
│   │   └── password/        # 密码工具
│   ├── router/              # 路由配置
│   ├── service/             # 业务逻辑服务
│   └── storage/             # 存储层
│       ├── models.go        # 数据模型
│       ├── storage.go       # 存储接口
│       ├── storage_facade.go # 存储门面
│       ├── memory_storage.go # 内存实现
│       └── gorm_storage.go   # GORM 实现
└── swagger/                 # API 文档
```

### 分层架构

#### 1. Handler 层
- **目的**：HTTP 请求/响应处理
- **职责**：
  - 解析 HTTP 请求
  - 验证输入数据
  - 调用服务层
  - 格式化 HTTP 响应
- **文件**：`internal/handler/*.go`

#### 2. Service 层
- **目的**：业务逻辑实现
- **职责**：
  - 实现业务规则
  - 协调处理器和存储层
  - 处理复杂操作
- **文件**：`internal/service/*.go`

#### 3. Storage 层
- **目的**：数据持久化抽象
- **职责**：
  - 定义存储接口
  - 实现存储后端（内存、GORM）
  - 处理数据操作
- **文件**：`internal/storage/*.go`

#### 4. Middleware 层
- **目的**：横切关注点
- **职责**：
  - 认证
  - CORS 处理
  - 请求日志
  - 限流
- **文件**：`internal/middleware/*.go`

## 数据模型

### 核心模型

#### User
```go
type User struct {
    ID           string
    Username     string
    PasswordHash string
    APIKey       string
    Role         string    // admin
    Status       string    // active, inactive, suspended
    LastLoginAt  *time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

#### LLMResource
```go
type LLMResource struct {
    ID        string
    Name      string
    Type      string    // chat, image, embedding 等
    Driver    string    // openai（目前仅支持 openai）
    Model     string
    BaseURL   string
    APIKey    string
    Status    string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### Token
```go
type Token struct {
    ID         string
    Name       string
    Token      string
    Prefix     string
    Status     string
    PolicyID   string
    LastUsedAt *time.Time
    ExpiresAt  *time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

#### Policy
```go
type Policy struct {
    ID         string
    Name       string
    TemplateID string
    Type       string    // round_robin, random 等
    Parameters string    // JSON
    Enabled    bool
    Builtin    bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

## 请求流程

### OpenAI 兼容 API 请求

```
1. 客户端请求
   ↓
2. 认证中间件
   - 验证 API 密钥/Token
   ↓
3. 请求日志中间件
   - 记录请求详情
   ↓
4. OpenAI 处理器
   - 解析请求
   - 选择资源（通过策略）
   ↓
5. 策略执行器
   - 执行路由策略
   - 选择 LLM 资源
   ↓
6. 客户端管理器
   - 创建/获取客户端
   - 转发请求到 AI 服务
   ↓
7. 响应处理
   - 格式化响应
   - 记录响应
   ↓
8. 返回客户端
```

### 管理 API 请求

```
1. 客户端请求
   ↓
2. 认证中间件
   - 验证管理员凭据
   ↓
3. 处理器
   - 解析请求
   - 验证输入
   ↓
4. 服务层
   - 执行业务逻辑
   - 更新存储
   ↓
5. 存储层
   - 持久化更改
   ↓
6. 响应
   - 返回结果
```

## 存储后端

### 内存存储
- **使用场景**：开发和测试
- **特点**：快速、临时、无持久化
- **实现**：`memory_storage.go`

### GORM 存储
- **使用场景**：生产环境
- **特点**：持久化，支持 SQLite/MySQL/PostgreSQL
- **实现**：`gorm_storage.go`

## 安全架构

### 认证流程

```
1. 客户端发送带 API 密钥/Token 的请求
   ↓
2. 认证中间件提取凭据
   ↓
3. 验证凭据：
   - 在 TokenService 中检查 Token
   - 或在 User 存储中检查 API 密钥
   ↓
4. 设置用户上下文
   ↓
5. 继续到处理器
```

### 密码安全
- 使用 bcrypt 哈希密码
- 绝不存储明文密码
- 密码验证使用恒定时间比较

## 负载均衡

### 支持的策略
- **轮询（Round Robin）**：顺序分配请求
- **随机（Random）**：随机选择
- **加权（Weighted）**：加权分配
- **模型匹配（Model Match）**：按模型名称匹配
- **正则匹配（Regex Match）**：按模式匹配
- **优先级（Priority）**：基于优先级选择
- **故障转移（Failover）**：自动故障转移

## 配置管理

### 配置源（优先级顺序）
1. 环境变量（`LINGPROXY_*`）
2. 配置文件（`config.yaml`）
3. 默认值（代码中）

### 配置结构
- 应用设置
- 存储配置
- 日志配置
- 安全设置

## 错误处理

### 错误类型
- **验证错误**：400 Bad Request
- **认证错误**：401 Unauthorized
- **未找到错误**：404 Not Found
- **服务器错误**：500 Internal Server Error

### 错误响应格式
```json
{
  "error": "错误消息"
}
```

## 日志

### 日志级别
- **Debug**：详细调试信息
- **Info**：一般信息消息
- **Warn**：警告消息
- **Error**：错误消息
- **Fatal**：致命错误

### 日志输出
- 控制台（stdout）
- 文件（轮转日志）
- 两者（推荐）

## 性能考虑

### 缓存
- 内存缓存用于频繁访问的数据
- 可配置 TTL

### 连接池
- HTTP 客户端连接池
- 可配置池大小

### 数据库优化
- 索引查询
- 高效数据模型
- 连接池
