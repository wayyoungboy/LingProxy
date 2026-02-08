# 构建 LingProxy 官网

本文档说明如何构建和部署 LingProxy 的 Docusaurus 官网。

## 前置要求

- Node.js 18.0 或更高版本
- npm 或 yarn

## 快速开始

### 1. 安装依赖

```bash
cd website
npm install
```

### 2. 本地开发

启动开发服务器（支持热重载）：

```bash
npm start
```

网站将在 `http://localhost:5000` 启动，修改文档后会自动刷新。

### 3. 构建生产版本

构建静态网站：

```bash
npm run build
```

构建完成后，静态文件会生成在 `build/` 目录中。

### 4. 预览构建结果

预览构建后的网站：

```bash
npm run serve
```

这会在本地启动一个服务器来预览构建结果。

## 部署方式

### 方式一：GitHub Pages

1. **配置 GitHub Pages**

   编辑 `docusaurus.config.js`：
   ```js
   url: 'https://lingproxy.github.io',
   baseUrl: '/lingproxy/',  // 如果部署在子目录
   ```

2. **部署**

   ```bash
   GIT_USER=<Your GitHub username> USE_SSH=true npm run deploy
   ```

   这会自动构建并推送到 `gh-pages` 分支。

### 方式二：Docker

1. **构建镜像**

   ```bash
   cd website
   docker build -t lingproxy-website .
   ```

2. **运行容器**

   ```bash
   docker run -d -p 5000:80 --name lingproxy-website lingproxy-website
   ```

   网站将在 `http://localhost:5000` 访问。

### 方式三：静态文件服务器

构建后，将 `build/` 目录的内容部署到任何静态文件服务器：

- **Vercel**: 连接 GitHub 仓库，自动部署
- **Netlify**: 拖拽 `build/` 目录或连接 GitHub
- **Nginx**: 将 `build/` 目录内容复制到 web 根目录
- **AWS S3 + CloudFront**: 上传到 S3，配置 CloudFront

### 方式四：使用 Makefile

如果项目根目录有 Makefile，可以添加构建命令：

```makefile
.PHONY: website-build website-serve website-deploy

website-build:
	cd website && npm install && npm run build

website-serve:
	cd website && npm run serve

website-deploy:
	cd website && GIT_USER=$(GIT_USER) USE_SSH=true npm run deploy
```

然后使用：

```bash
make website-build    # 构建
make website-serve    # 预览
make website-deploy   # 部署到 GitHub Pages
```

## 文档结构

网站从项目根目录的 `docs/` 读取文档：

```
docs/
├── en/              # 英文文档
│   ├── README.md
│   ├── 01-introduction.md
│   └── ...
└── zh/              # 中文文档
    ├── README.md
    ├── 01-introduction.md
    └── ...
```

## 配置说明

### 修改网站信息

编辑 `website/docusaurus.config.js`：

```js
title: 'LingProxy',                    // 网站标题
tagline: 'High-performance AI API Gateway',  // 副标题
url: 'https://lingproxy.github.io',     // 网站 URL
baseUrl: '/',                           // 基础路径
```

### 修改侧边栏

编辑 `website/sidebars.js` 来调整文档导航结构。

### 添加 Logo

将 Logo 文件放到 `website/static/img/` 目录：
- `logo.svg` - 浅色模式 Logo
- `logo-dark.svg` - 深色模式 Logo（可选）
- `favicon.ico` - 网站图标

## 常见问题

### 构建失败

1. **清除缓存后重新构建**：
   ```bash
   npm run clear
   npm run build
   ```

2. **重新安装依赖**：
   ```bash
   rm -rf node_modules package-lock.json
   npm install
   ```

### 端口被占用

使用不同端口启动：

```bash
npm start -- --port 6000
```

### 文档未更新

确保 `docusaurus.config.js` 中配置的文档路径正确：

```js
docs: {
  path: '../docs',  // 从项目根目录的 docs/ 读取
  // ...
}
```

## 更多信息

- [Docusaurus 官方文档](https://docusaurus.io/docs)
- [项目 README](../README.md)
