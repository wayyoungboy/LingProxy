# LingProxy 客户端库

本目录包含多种语言的 LingProxy 客户端实现，方便不同技术栈的开发者使用。

## 支持的客户端

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

## 更多信息

- 完整的 API 文档请参考 [API 参考文档](../../docs/zh/04-api-reference.md)
- 项目主页: [LingProxy README](../../README_zh.md)
