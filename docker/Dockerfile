# ============================================================================
# 多阶段构建：统一的前后端 Dockerfile
# 只暴露一个端口（80），同时提供前端和后端服务
# ============================================================================

# 阶段1: 构建后端
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 设置 Go 代理（可选，加速依赖下载）
ENV GOPROXY=https://goproxy.cn,direct

# 复制go mod文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY backend/ .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o lingproxy \
    ./cmd/main.go

# 阶段2: 构建前端
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git

# 设置 npm 镜像（可选，加速依赖下载）
RUN npm config set registry https://registry.npmmirror.com

# 复制 package 文件
COPY frontend/package.json frontend/package-lock.json ./

# 安装依赖
RUN npm ci

# 复制源代码
COPY frontend/ .

# 构建应用
RUN npm run build

# 阶段3: 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add \
    ca-certificates \
    curl \
    nginx \
    supervisor \
    tzdata \
    && update-ca-certificates

# 创建非root用户
RUN addgroup -g 1000 -S lingproxy && \
    adduser -u 1000 -S lingproxy -G lingproxy

# 设置工作目录
WORKDIR /app

# 从后端构建阶段复制二进制文件
COPY --from=backend-builder /app/lingproxy /app/lingproxy

# 从后端构建阶段复制配置文件
COPY --from=backend-builder /app/configs /app/configs

# 从前端构建阶段复制构建产物到 nginx 目录
COPY --from=frontend-builder /app/dist /usr/share/nginx/html

# 复制 nginx 配置文件（修改为代理到 localhost:8080）
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

# 复制 supervisor 配置文件
COPY docker/supervisord.conf /etc/supervisord.conf

# 创建必要的目录
RUN mkdir -p /app/data /app/logs /app/run /var/log/supervisor && \
    chown -R lingproxy:lingproxy /app && \
    chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/log/nginx

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口（只暴露一个端口）
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD curl -f http://localhost/api/v1/health || exit 1

# 使用 supervisor 启动所有服务
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
