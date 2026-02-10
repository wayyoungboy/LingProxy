# 快速开始指南

## 前置要求

- **后端**：Go 1.21 或更高版本
- **前端**：Node.js 18+ 和 npm/yarn
- **数据库**：SQLite（已包含）或 MySQL/PostgreSQL（可选）

## 安装

### 1. 克隆仓库

```bash
git clone https://github.com/your-org/lingproxy.git
cd lingproxy
```

### 2. 后端设置

#### 安装依赖

```bash
cd backend
go mod tidy
```

#### 配置应用

```bash
# 复制示例配置文件
cp configs/config.yaml.example configs/config.yaml

# 编辑配置文件
# 默认管理员凭据：
# 用户名: admin
# 密码: admin123
```

#### 运行后端

```bash
# 开发模式
go run cmd/main.go

# 或构建后运行
go build -o lingproxy cmd/main.go
./lingproxy
```

后端将在 `http://localhost:8080` 启动

### 3. 前端设置

#### 安装依赖

```bash
cd frontend
npm install
```

#### 运行前端

```bash
npm run dev
```

前端将在 `http://localhost:3000` 可用

### 4. 访问管理后台

1. 打开浏览器访问 `http://localhost:3000`
2. 使用默认凭据登录：
   - 用户名：`admin`
   - 密码：`admin123`
3. **重要**：首次登录后请修改默认密码！

## Docker 部署

### 前置要求

- **Docker**: Docker 20.10+ 和 Docker Compose 2.0+
- **配置文件**: 确保 `backend/configs/config.yaml` 存在且配置正确

### 使用 Makefile 快速启动（推荐）

最简单的启动方式：

```bash
# 在项目根目录执行 - 一条命令启动所有服务
make docker-compose-up
```

此命令会自动完成：
1. ✅ 检查配置文件是否存在（不存在则从示例创建）
2. ✅ 创建必要的目录（logs, run）
3. ✅ 启动 SeekDB 和 LingProxy 服务（需要时自动构建）
4. ✅ 等待 SeekDB 就绪
5. ✅ 自动创建数据库
6. ✅ 显示访问地址

**访问地址**（启动后）：
- **后端 API**: http://localhost:8080/api/v1
- **健康检查**: http://localhost:8080/api/v1/health

**注意**：Docker 部署仅包含后端服务。前端需要单独运行：
```bash
cd frontend
npm run dev
# 前端将在 http://localhost:3000 可用
```

### 手动使用 Docker Compose

如果更喜欢直接使用 Docker Compose：

```bash
# 1. 准备配置文件
cp backend/configs/config.yaml.example backend/configs/config.yaml
# 编辑 backend/configs/config.yaml，配置 SeekDB 连接：
# storage:
#   type: "gorm"
#   gorm:
#     driver: "mysql"
#     dsn: "root:@tcp(seekdb:2881)/lingproxy?charset=utf8mb4&parseTime=True&loc=Local"

# 2. 启动服务
docker-compose -f docker/docker-compose.yml up -d --build

# 3. 数据库会在后端启动时自动创建
# 无需手动创建数据库

# 4. 查看日志
docker-compose -f docker/docker-compose.yml logs -f

# 5. 停止服务
docker-compose -f docker/docker-compose.yml down
```

### 其他有用的 Makefile 命令

```bash
# 查看服务状态
make docker-compose-ps

# 查看日志
make docker-compose-logs

# 停止服务
make docker-compose-down

# 重启服务
make docker-compose-restart

# 仅初始化数据库
make docker-compose-init-db
```

### 服务架构

Docker 部署采用**前后端分离**架构：

- **SeekDB**: MySQL 兼容数据库（端口 2881, 2886）
  - 数据存储在 Docker volume `seekdb-data` 中
  - 健康检查确保服务就绪后再启动后端

- **LingProxy 后端**: 后端 API 服务（端口 8080）
  - 纯后端服务，不包含前端
  - 启动时如果数据库不存在会自动创建
  - 使用 Docker 专用配置（`config.yaml.docker`）

- **前端**: 开发时需单独运行
  - 在 `frontend` 目录下使用 `npm run dev`
  - 运行在 3000 端口（Vite 开发服务器）
  - API 请求通过 Vite 代理转发到 `http://localhost:8080`

### 故障排查

**服务无法启动：**
```bash
# 检查服务状态
make docker-compose-ps

# 查看日志
make docker-compose-logs

# 检查端口是否被占用
lsof -i :8080
lsof -i :2881
```

**数据库连接问题：**
```bash
# 验证 SeekDB 是否运行
docker exec seekdb mysql -h127.0.0.1 -uroot -P2881 -e "SHOW DATABASES;"

# 重新创建数据库
make docker-compose-init-db
```

**配置问题：**
```bash
# 检查配置文件
cat backend/configs/config.yaml

# 验证 SeekDB 连接字符串
grep -A 3 "storage:" backend/configs/config.yaml
```

更多详情请参考[配置指南](03-configuration.md)和[开发指南](06-development.md)。

## 第一步

### 1. 配置 LLM 资源

1. 在管理后台导航到 **LLM 资源管理**
2. 点击 **添加资源**
3. 填写资源信息：
   - 名称：描述性名称
   - 类型：模型类别（对话、生图、嵌入等）
   - 驱动：目前仅支持 "openai"
   - 模型：模型标识（如 gpt-4, gpt-3.5-turbo）
   - 基础 URL：API 端点 URL
   - API 密钥：您的 API 密钥

### 2. 创建 API Key

1. 在管理后台导航到 **API Key 管理**
2. 点击 **创建 API Key**
3. 填写 API Key 信息：
   - 名称：API Key 名称/描述
   - 策略：选择路由策略（可选）

### 3. 测试 API

```bash
# 使用 curl（后端运行在 8080 端口）
curl -X POST http://localhost:8080/llm/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "你好！"}
    ]
  }'
```

**注意**：如果您在本地运行前端（`npm run dev`），前端的 API 请求会自动代理到后端。

## 下一步

- 阅读 [配置指南](03-configuration.md) 了解详细配置选项
- 查看 [API 文档](04-api-reference.md) 了解 API 使用方法
- 查看 [架构指南](05-architecture.md) 了解系统设计详情
- 查看 [开发指南](06-development.md) 了解如何贡献代码
