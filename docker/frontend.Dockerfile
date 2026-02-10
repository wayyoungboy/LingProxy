# ============================================================================
# 前端专用 Dockerfile
# 构建前端静态文件并使用 Nginx 提供服务
# ============================================================================

# 构建阶段
FROM node:20-alpine AS builder

# 配置 Alpine 镜像源（使用阿里云镜像，加速国内访问）
RUN if [ -f /etc/apk/repositories ]; then \
        sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
    else \
        echo "https://mirrors.aliyun.com/alpine/v3.19/main" > /etc/apk/repositories && \
        echo "https://mirrors.aliyun.com/alpine/v3.19/community" >> /etc/apk/repositories; \
    fi

WORKDIR /app

# 安装构建依赖
RUN apk update && apk add --no-cache git

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

# 运行阶段
FROM nginx:alpine

# 配置 Alpine 镜像源（使用阿里云镜像，加速国内访问）
RUN if [ -f /etc/apk/repositories ]; then \
        sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
    else \
        echo "https://mirrors.aliyun.com/alpine/v3.19/main" > /etc/apk/repositories && \
        echo "https://mirrors.aliyun.com/alpine/v3.19/community" >> /etc/apk/repositories; \
    fi

# 安装 curl 用于健康检查
RUN apk update && apk add --no-cache curl

# 从构建阶段复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 Nginx 配置文件
COPY docker/nginx-frontend.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]
