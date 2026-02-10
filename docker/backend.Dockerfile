# ============================================================================
# 后端专用 Dockerfile
# 只构建和运行后端服务
# ============================================================================

# 构建阶段
FROM golang:1.24-alpine AS builder

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

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add \
    ca-certificates \
    curl \
    tzdata \
    && update-ca-certificates

# 创建非root用户
RUN addgroup -g 1000 -S lingproxy && \
    adduser -u 1000 -S lingproxy -G lingproxy

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/lingproxy /app/lingproxy

# 从构建阶段复制配置文件
COPY --from=builder /app/configs /app/configs
# 在 Docker 环境中使用 Docker 专用配置文件
RUN cp /app/configs/config.yaml.docker /app/configs/config.yaml

# 创建必要的目录
RUN mkdir -p /app/logs /app/run && \
    chown -R lingproxy:lingproxy /app

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口（后端服务端口）
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/health || exit 1

# 切换到非root用户
USER lingproxy

# 启动后端服务
CMD ["/app/lingproxy"]
