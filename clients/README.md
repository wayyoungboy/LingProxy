# LingProxy 客户端库

本目录包含多种语言的 LingProxy 客户端实现，方便不同技术栈的开发者使用。

## 支持的客户端

### 🌐 Web 测试页面（推荐用于快速测试）

独立的静态 HTML 页面，可以直接在浏览器中打开，无需安装任何依赖。

- 📁 位置: `test-page.html`
- ✨ 特性:
  - ✅ 对话接口测试
  - ✅ 文本补全测试
  - ✅ 模型列表查询
  - ✅ 流式响应支持
  - ✅ 配置自动保存
- 🚀 快速开始:
  ```bash
  # 直接在浏览器中打开
  open clients/test-page.html  # macOS
  # 或使用简单 HTTP 服务器
  cd clients
  python3 -m http.server 8000
  # 访问 http://localhost:8000/test-page.html
  ```

### Python 客户端

标准 Python 客户端，基于 OpenAI Python SDK。

- 📁 位置: `python/`
- 📖 文档: [Python 客户端文档](python/README.md)
- 🚀 快速开始:
  ```bash
  cd python
  pip install -r requirements.txt
  python lingproxy_client.py
  ```

### JavaScript/TypeScript 客户端

标准 JavaScript/TypeScript 客户端，基于 OpenAI JavaScript SDK。

- 📁 位置: `javascript/`
- 📖 文档: [JavaScript 客户端文档](javascript/README.md)
- 🚀 快速开始:
  ```bash
  cd javascript
  npm install
  npm run demo
  ```

### Go 客户端

标准 Go 客户端库，基于 OpenAI Go SDK。

- 📁 位置: `go/`
- 📖 文档: [Go 客户端文档](go/README.md)
- 🚀 快速开始:
  ```bash
  cd go
  go run example.go
  ```

## 通用特性

所有客户端都提供以下特性：

- ✅ OpenAI API 兼容
- ✅ 环境变量支持（`LINGPROXY_API_KEY`）
- ✅ 聊天补全（Chat Completions）
- ✅ 文本补全（Text Completions）
- ✅ 模型列表（List Models）
- ✅ 流式响应支持（Streaming）
- ✅ 完整的类型定义

## 使用示例

### Python

```python
from lingproxy_client import LingProxyClient

client = LingProxyClient(api_key="your-api-key")
response = client.create_chat_completion(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": "Hello!"}]
)
```

### JavaScript

```javascript
import { LingProxyClient } from './lingproxy-client.js';

const client = new LingProxyClient({ apiKey: 'your-api-key' });
const response = await client.createChatCompletion({
  model: 'gpt-3.5-turbo',
  messages: [{ role: 'user', content: 'Hello!' }]
});
```

### Go

```go
import "github.com/lingproxy/lingproxy/clients/go/client"

c, _ := client.NewClient(&client.ClientOptions{APIKey: "your-api-key"})
resp, _ := c.CreateChatCompletion(ctx, client.ChatCompletionRequest{
    Model: "gpt-3.5-turbo",
    Messages: []client.ChatMessage{
        {Role: "user", Content: "Hello!"},
    },
})
```

## 环境变量

所有客户端都支持通过环境变量设置 API Key：

```bash
export LINGPROXY_API_KEY=your-api-key
```

## 测试工具

### 🌐 Web 测试页面（推荐用于快速测试）

独立的静态 HTML 页面，可以直接在浏览器中打开，用于测试 LingProxy 的 API 转发接口。

#### 功能特性

- ✅ **对话接口测试** - 测试 `/llm/v1/chat/completions` 接口
- ✅ **文本补全测试** - 测试 `/llm/v1/completions` 接口
- ✅ **模型列表查询** - 查询可用的模型列表
- ✅ **流式响应支持** - 支持流式对话响应
- ✅ **配置持久化** - 自动保存 API URL 和 Key 到本地存储
- ✅ **美观界面** - 现代化的 UI 设计，易于使用

#### 使用方法

**1. 打开测试页面**

直接在浏览器中打开 `test-page.html` 文件：

```bash
# 方法1: 直接双击文件
# 方法2: 在浏览器中打开
open clients/test-page.html  # macOS
xdg-open clients/test-page.html  # Linux
start clients/test-page.html  # Windows
```

或者使用简单的 HTTP 服务器：

```bash
# Python 3
cd clients
python3 -m http.server 8000
# 然后访问 http://localhost:8000/test-page.html

# Node.js (需要安装 http-server)
npx http-server -p 8000
# 然后访问 http://localhost:8000/test-page.html
```

**2. 配置 API 信息**

在页面顶部填写：

- **API Base URL**: LingProxy 服务的地址，默认 `http://localhost:8080/llm/v1`
- **API Key**: 你的 LingProxy API Key（格式：`ling-xxxxx`）

> 💡 提示：配置会自动保存到浏览器的本地存储，下次打开页面时会自动填充。

**3. 测试对话接口**

1. 切换到 **"💬 对话接口"** 标签页
2. 填写模型名称（例如：`gpt-3.5-turbo`）
3. 添加消息（支持 user、assistant、system 角色）
4. 可选：启用流式响应、调整 Temperature、Max Tokens
5. 点击 **"发送请求"** 按钮

**响应结果**：
- 显示 AI 的回复内容
- 显示 Token 使用统计（Prompt Tokens、Completion Tokens、Total Tokens）

**4. 测试文本补全接口**

1. 切换到 **"📝 文本补全"** 标签页
2. 填写模型名称（例如：`text-davinci-003`）
3. 输入提示词（Prompt）
4. 调整参数（Temperature、Max Tokens）
5. 点击 **"发送请求"** 按钮

**5. 查询模型列表**

1. 切换到 **"📋 模型列表"** 标签页
2. 点击 **"加载模型列表"** 按钮
3. 查看所有可用的模型

#### 功能说明

**对话接口**：
- **多消息支持**：可以添加多条消息，模拟多轮对话
- **角色选择**：支持 user、assistant、system 三种角色
- **流式响应**：勾选"启用流式响应"可以实时查看 AI 生成的内容
- **参数调整**：
  - Temperature: 控制输出的随机性（0-2）
  - Max Tokens: 限制生成的最大 Token 数

**文本补全接口**：
- **单次补全**：输入提示词，AI 会补全后续内容
- **参数调整**：支持 Temperature 和 Max Tokens 参数

**模型列表**：
- **快速查询**：一键查询所有可用的模型
- **模型识别**：帮助确认配置的模型是否正确

#### 注意事项

1. **CORS 问题**：如果遇到跨域问题，请确保：
   - LingProxy 服务已配置允许跨域请求
   - 或者使用浏览器扩展禁用 CORS（仅用于测试）

2. **API Key 格式**：确保 API Key 以 `ling-` 开头

3. **服务地址**：确保 LingProxy 服务正在运行，并且地址正确

4. **网络连接**：确保浏览器可以访问 LingProxy 服务地址

#### 故障排除

**问题：无法连接到服务**

**解决方案**：
- 检查 LingProxy 服务是否正在运行
- 确认 API Base URL 是否正确
- 检查防火墙设置

**问题：401 未授权错误**

**解决方案**：
- 检查 API Key 是否正确
- 确认 Token 状态是否为 `active`
- 检查 LingProxy 配置中的认证是否已启用

**问题：404 未找到接口**

**解决方案**：
- 确认 API Base URL 路径正确（应该是 `/llm/v1`）
- 检查 LingProxy 路由配置

**问题：模型不存在错误**

**解决方案**：
- 在 LingProxy 管理界面添加 LLM 资源
- 确认模型名称正确
- 使用"模型列表"功能查询可用模型

#### 技术说明

- **纯静态页面**：无需服务器，可直接在浏览器中打开
- **使用 Fetch API**：使用现代浏览器 API 发送请求
- **本地存储**：使用 localStorage 保存配置
- **流式处理**：支持 Server-Sent Events (SSE) 流式响应

#### 浏览器兼容性

- ✅ Chrome/Edge (推荐)
- ✅ Firefox
- ✅ Safari
- ✅ Opera

需要支持：
- Fetch API
- async/await
- localStorage

#### 安全提示

⚠️ **重要**：此测试页面仅用于开发和测试环境，不要在生产环境中使用。

- API Key 会保存在浏览器的本地存储中
- 不要在公共计算机上使用此页面
- 测试完成后建议清除浏览器数据

## 更多信息

- 完整的 API 文档请参考 [API 参考文档](../../docs/zh/04-api-reference.md)
- 项目主页: [LingProxy README](../../README_zh.md)
