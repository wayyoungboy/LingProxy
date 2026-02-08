# 配置指南

## 配置文件

主配置文件位于 `backend/configs/config.yaml`。复制 `config.yaml.example` 来创建您的配置文件。

## 配置结构

### 应用配置

```yaml
app:
  environment: "development"  # development, staging, production
  port: 8080                  # 服务器端口
  # name 和 version 有默认值，通常不需要配置
  # host 有默认值 "0.0.0.0"，通常不需要配置
```

### 存储配置

```yaml
storage:
  type: "gorm"  # memory 或 gorm
  gorm:
    driver: "sqlite"  # sqlite, mysql
    dsn: "lingproxy.db"
    # MySQL 示例：
    # driver: "mysql"
    # dsn: "username:password@tcp(localhost:3306)/lingproxy?charset=utf8mb4&parseTime=True&loc=Local"
```

**存储类型：**
- `memory`: 内存存储（用于开发/测试）
- `gorm`: 数据库存储（SQLite、MySQL、PostgreSQL）

### 日志配置

```yaml
log:
  level: "info"      # debug, info, warn, error, fatal
  format: "json"     # text, json
  output: "both"     # stdout, file, both（推荐：both）
  file_path: "./logs/lingproxy.log"
  # max_size, max_backups, max_age, compress 有默认值
```

**日志级别：**
- `debug`: 详细的调试信息
- `info`: 一般信息消息
- `warn`: 警告消息
- `error`: 错误消息
- `fatal`: 致命错误

**输出模式：**
- `stdout`: 仅控制台输出
- `file`: 仅文件输出
- `both`: 控制台和文件（推荐）

### 安全配置

```yaml
security:
  auth:
    enabled: true  # 全局启用/禁用认证
  # cors, rate_limit, jwt 可通过管理界面配置
```

**认证：**
- 当 `enabled: true` 时：所有 API（除登录外）都需要认证
- 当 `enabled: false` 时：所有 API（除登录外）都无需认证即可访问

### Provider配置

```yaml
provider:
  timeout: "30s"        # 请求超时时间
  max_retries: 3        # 最大重试次数（0 = 禁用重试）
  retry_delay: "1s"     # 基础重试延迟（实际延迟会指数增长）
  max_idle_conns: 100   # 最大空闲连接数
  max_conns_per_host: 100  # 每个主机最大连接数
  idle_conn_timeout: "90s"  # 空闲连接超时时间
```

**重试配置说明：**
- `max_retries`: 失败请求的最大重试次数。设置为 `0` 可禁用重试。
- `retry_delay`: 重试之间的基础延迟时间。实际延迟会指数增长：`延迟 = retry_delay × 重试次数`。
- 重试适用于网络错误、超时和5xx服务器错误。
- 不会重试4xx客户端错误（429限流除外）、认证错误或上下文取消错误。
- 可通过管理界面配置：设置 → Provider设置。

## 环境变量

您可以使用 `LINGPROXY_` 前缀的环境变量覆盖配置值：

```bash
# 示例：覆盖端口
export LINGPROXY_APP_PORT=9000

# 示例：覆盖数据库 DSN
export LINGPROXY_STORAGE_GORM_DSN="mysql://user:pass@localhost/db"
```

## 默认值

所有配置选项都有合理的默认值。完整的默认值请参阅 `backend/internal/config/config.go`：

- `app.name`: "LingProxy"
- `app.version`: "1.0.0"
- `app.host`: "0.0.0.0"
- `app.port`: 8080
- `app.environment`: "development"
- `storage.type`: "memory"
- `log.level`: "info"
- `log.format`: "json"
- `log.output`: "both"
- `security.auth.enabled`: true
- `provider.timeout`: "30s"
- `provider.max_retries`: 3
- `provider.retry_delay`: "1s"

## 管理员凭据

默认管理员凭据（首次启动时设置）：
- 用户名：`admin`
- 密码：`admin123`

**重要**：首次登录后请修改默认密码！

## 通过管理界面配置

许多设置可以通过管理后台配置：
- **系统设置**：基础设置、缓存、限流、安全、日志、负载均衡、Provider重试
- **LLM 资源**：添加、编辑和管理 AI 服务资源
- **模型**：配置模型详情、定价和使用限制
- **策略**：创建和管理路由策略
- **API Key**：创建和管理请求端 API Key，支持策略绑定

通过管理界面所做的更改会保存到配置文件。

### Provider重试设置

您可以通过管理界面配置资源请求的重试行为：

1. 导航到 **设置** → **Provider设置**
2. 配置：
   - **请求超时时间**：等待响应的最大时间（秒）
   - **最大重试次数**：最大重试尝试次数（0-10，0 = 禁用）
   - **重试延迟**：重试之间的基础延迟（秒，实际延迟会指数增长）

这些设置会立即应用到所有新请求，无需重启服务。

## 生产环境配置

生产环境部署时：

1. 将 `app.environment` 设置为 `"production"`
2. 使用 `gorm` 存储类型和生产数据库（MySQL/PostgreSQL）
3. 启用认证（`security.auth.enabled: true`）
4. 配置适当的日志轮转
5. 使用强密码和 API 密钥
6. 为您的域名适当配置 CORS

## 故障排除

### 配置文件未找到

如果找不到配置文件，应用程序将使用默认值。检查日志输出中的配置加载消息。

### 无效配置

如果配置验证失败，请检查错误消息并验证：
- YAML 语法正确
- 必填字段存在
- 值在有效范围内

### 配置未生效

修改配置文件后：
1. 重启应用程序
2. 检查日志中的配置加载消息
3. 通过管理界面验证更改
