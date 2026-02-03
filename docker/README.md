# Docker 部署指南

本目录包含 LingProxy 的 Docker 部署相关文件。

## 文件说明

- `Dockerfile` - Docker 镜像构建文件
- `docker-compose.yml` - Docker Compose 配置文件（从项目根目录运行）
- `docker-compose.example.yml` - Docker Compose 配置示例文件
- `.dockerignore` - Docker 构建忽略文件列表

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
  -p 8080:8080 \
  -v $(pwd)/backend/configs:/app/configs:ro \
  -v $(pwd)/logs:/app/logs \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/run:/app/run \
  -e GIN_MODE=release \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  lingproxy:latest
```

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

- `8080:8080` - HTTP API 端口

## 健康检查

容器包含健康检查配置，每30秒检查一次 `/api/v1/health` 端点。

## 注意事项

1. **配置文件**: 运行前必须创建 `backend/configs/config.yaml` 文件
2. **数据持久化**: 确保数据目录（`data/`）和日志目录（`logs/`）有适当的权限
3. **网络**: 默认使用 `lingproxy-network` 网络，如需与其他服务通信，请配置网络
