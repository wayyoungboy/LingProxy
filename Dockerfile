# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git

# 设置 Go 代理（可选，加速依赖下载）
ENV GOPROXY=https://goproxy.cn,direct

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
# 使用 CGO_ENABLED=0 构建静态二进制文件，便于在 alpine 中运行
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o lingproxy \
    ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
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
COPY --from=builder /app/lingproxy .

# 复制配置文件
COPY --from=builder /app/configs ./configs

# 创建必要的目录（数据目录、日志目录）
RUN mkdir -p /app/data /app/logs && \
    chown -R lingproxy:lingproxy /app

# 切换到非root用户
USER lingproxy

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/health || exit 1

# 启动应用
CMD ["./lingproxy"]
