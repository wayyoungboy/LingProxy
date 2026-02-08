# Docker 部署指南

本目录包含 LingProxy 的 Docker 部署相关文件。

## 文件说明

- `Dockerfile` - 统一的前后端 Dockerfile（只暴露一个端口）
- `docker-compose.yml` - Docker Compose 配置文件
- `docker-compose.example.yml` - Docker Compose 配置示例文件
- `nginx.conf` - Nginx 配置文件（代理到本地后端服务）
- `supervisord.conf` - Supervisor 配置（管理 Nginx 和后端服务）
- `.dockerignore` - Docker 构建忽略文件列表

## 架构说明

统一部署架构：
- **只暴露一个端口（80）**
- **前端和后端在同一个容器中**
- **使用 Supervisor 管理多个进程**
- **Nginx 服务前端静态文件并代理 API 请求到本地后端服务**

```
容器内部：
├─ Supervisor (root)
│   ├─ Nginx (nginx用户) → 监听 80 端口
│   │   ├─ 服务前端静态文件
│   │   └─ 代理 /api, /llm, /swagger → 127.0.0.1:8080
│   └─ LingProxy (lingproxy用户) → 监听 127.0.0.1:8080

对外暴露：
└─ 只暴露 80 端口
```

## 快速开始

### 方式一：使用 Docker Compose（推荐）

1. **准备配置文件**
```bash
# 在项目根目录执行
cp backend/configs/config.yaml.example backend/configs/config.yaml
# 编辑 backend/configs/config.yaml 根据需要修改配置
```

2. **启动服务**
```bash
# 在项目根目录执行
docker-compose -f docker/docker-compose.yml up -d

# 查看日志
docker-compose -f docker/docker-compose.yml logs -f

# 停止服务
docker-compose -f docker/docker-compose.yml down
```

3. **访问服务**
- **前端界面**: http://localhost
- **后端 API**: http://localhost/api/v1
- **API 文档**: http://localhost/swagger/index.html

### 方式二：直接使用 Docker

1. **构建镜像**
```bash
# 在项目根目录执行
docker build -f docker/Dockerfile -t lingproxy:latest .
```

2. **运行容器**
```bash
docker run -d \
  --name lingproxy \
  -p 80:80 \
  -v $(pwd)/backend/configs:/app/configs:ro \
  -v $(pwd)/logs:/app/logs \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/run:/app/run \
  -e GIN_MODE=release \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  lingproxy:latest
```

3. **访问服务**
- **前端界面**: http://localhost
- **后端 API**: http://localhost/api/v1
- **API 文档**: http://localhost/swagger/index.html

## 配置说明

### 卷挂载

- `backend/configs:/app/configs:ro` - 配置文件目录（只读）
- `logs:/app/logs` - 日志目录
- `data:/app/data` - 数据目录（SQLite数据库等）
- `run:/app/run` - 运行目录（PID文件等）

### 环境变量

- `GIN_MODE=release` - Gin框架运行模式（release/production）
- `TZ=Asia/Shanghai` - 时区设置

### 端口

- `80:80` - 统一端口（前端 + 后端）

## 服务说明

- **容器名**：`lingproxy`
- **端口**：80（统一端口）
- **内部架构**：
  - Nginx（端口80）：服务前端静态文件 + 代理 API 请求
  - 后端服务（端口8080）：在容器内运行，只监听 localhost
  - Supervisor：管理 Nginx 和后端服务进程
- **健康检查**：每30秒检查一次 `/api/v1/health` 端点

## 健康检查

每30秒检查一次 `http://localhost/api/v1/health` 端点。

## 访问地址

启动服务后，可以通过以下地址访问：

- **前端界面**: http://localhost
- **后端 API**: http://localhost/api/v1
- **API 文档**: http://localhost/swagger/index.html

前端会自动代理 `/api`、`/llm`、`/swagger` 路径的请求到后端服务。

## 注意事项

1. **配置文件**: 运行前必须创建 `backend/configs/config.yaml` 文件
2. **数据持久化**: 确保数据目录（`data/`）和日志目录（`logs/`）有适当的权限
3. **端口冲突**: 确保 80 端口未被占用
4. **前端构建**: 前端镜像构建时会自动执行 `npm ci` 和 `npm run build`，构建时间可能较长
5. **统一部署说明**：
   - 使用 Supervisor 同时管理 Nginx 和后端服务
   - 后端服务在容器内监听 `127.0.0.1:8080`，不对外暴露
   - Nginx 代理 `/api`、`/llm`、`/swagger` 到本地后端服务
   - 只暴露 80 端口，简化部署和防火墙配置

## 优势

- ✅ **只暴露一个端口**：简化防火墙配置
- ✅ **前后端统一**：前端和后端在同一个容器中，部署简单
- ✅ **进程管理**：使用 Supervisor 管理进程，自动重启
- ✅ **适合生产环境**：简化部署流程
