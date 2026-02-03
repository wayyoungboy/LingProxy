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

### 使用 Docker Compose（推荐）

```bash
# 在项目根目录执行
docker-compose -f docker/docker-compose.yml up -d

# 查看日志
docker-compose -f docker/docker-compose.yml logs -f

# 停止服务
docker-compose -f docker/docker-compose.yml down
```

更多详情请参阅 [Docker 部署指南](../docker/README.md)。

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

### 2. 创建 Token

1. 在管理后台导航到 **Token 管理**
2. 点击 **创建 Token**
3. 填写 Token 信息：
   - 名称：Token 名称/描述
   - 策略：选择路由策略（可选）

### 3. 测试 API

```bash
# 使用 curl
curl -X POST http://localhost:8080/llm/v1/chat/completions \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "你好！"}
    ]
  }'
```

## 下一步

- 阅读 [配置指南](03-configuration.md) 了解详细配置选项
- 查看 [API 文档](04-api-reference.md) 了解 API 使用方法
- 查看 [架构指南](05-architecture.md) 了解系统设计详情
- 查看 [开发指南](06-development.md) 了解如何贡献代码
