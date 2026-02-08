# LingProxy 客户端快速开始指南

本指南帮助你快速开始使用 LingProxy 的各个客户端。

## 前置要求

1. 确保 LingProxy 服务正在运行（默认地址：`http://localhost:8080`）
2. 获取你的 API Key（从 LingProxy 管理界面创建 API Key）

## Python 客户端

### 安装

```bash
cd clients/python
pip install -r requirements.txt
```

### 基本使用

```python
from lingproxy_client import LingProxyClient

# 初始化客户端（从环境变量读取 API Key）
client = LingProxyClient()

# 或者直接指定 API Key
client = LingProxyClient(
    api_key="ling-your-api-key-here",
    base_url="http://localhost:8080/llm/v1"
)

# 创建聊天补全
response = client.create_chat_completion(
    model="gpt-3.5-turbo",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)

print(response['choices'][0]['message']['content'])
```

### 运行示例

```bash
export LINGPROXY_API_KEY=your-api-key
python lingproxy_client.py
```

## JavaScript/TypeScript 客户端

### 安装

```bash
cd clients/javascript
npm install
```

### 基本使用

```javascript
import { LingProxyClient } from './lingproxy-client.js';

// 初始化客户端
const client = new LingProxyClient({
  apiKey: 'ling-your-api-key-here',
  baseURL: 'http://localhost:8080/llm/v1'
});

// 创建聊天补全
const response = await client.createChatCompletion({
  model: 'gpt-3.5-turbo',
  messages: [
    { role: 'user', content: 'Hello!' }
  ]
});

console.log(response.choices[0].message.content);
```

### 运行示例

```bash
export LINGPROXY_API_KEY=your-api-key
npm run demo
```

## Go 客户端

### 安装

```bash
cd clients/go
go mod init lingproxy-client-example
go get github.com/openai/openai-go/v3
```

### 基本使用

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/lingproxy/lingproxy/clients/go/client"
)

func main() {
    // 初始化客户端
    c, err := client.NewClient(&client.ClientOptions{
        APIKey:  "ling-your-api-key-here",
        BaseURL: "http://localhost:8080/llm/v1",
    })
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    // 创建聊天补全
    temp := 0.7
    req := client.ChatCompletionRequest{
        Model: "gpt-3.5-turbo",
        Messages: []client.ChatMessage{
            {Role: "user", Content: "Hello!"},
        },
        Temperature: &temp,
    }
    
    resp, err := c.CreateChatCompletion(ctx, req)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(resp.Choices[0].Message.Content)
}
```

### 运行示例

```bash
export LINGPROXY_API_KEY=your-api-key
go run example.go
```

## 环境变量配置

所有客户端都支持通过环境变量设置 API Key：

```bash
# Linux/macOS
export LINGPROXY_API_KEY=ling-your-api-key-here

# Windows (PowerShell)
$env:LINGPROXY_API_KEY="ling-your-api-key-here"

# Windows (CMD)
set LINGPROXY_API_KEY=ling-your-api-key-here
```

## 常见问题

### 1. 如何获取 API Key？

1. 启动 LingProxy 服务
2. 访问管理界面（通常是 `http://localhost:3000`）
3. 登录后，在 API Key 管理页面创建新的 API Key
4. 复制生成的 API Key（格式：`ling-xxxxx`）

### 2. 连接失败怎么办？

- 检查 LingProxy 服务是否正在运行
- 确认 `base_url` 配置正确（默认：`http://localhost:8080/llm/v1`）
- 检查防火墙设置

### 3. 认证失败怎么办？

- 确认 API Key 正确（以 `ling-` 开头）
- 检查 API Key 状态是否为 `active`
- 确认认证功能已启用（检查配置文件中的 `security.auth.enabled`）

### 4. 模型不存在错误？

- 在 LingProxy 管理界面添加 LLM 资源
- 确认模型名称正确
- 检查资源状态是否为 `active`

## 更多资源

- [完整 API 文档](../../docs/zh/04-api-reference.md)
- [Python 客户端文档](python/README.md)
- [JavaScript 客户端文档](javascript/README.md)
- [Go 客户端文档](go/README.md)
- [项目主页](../../README_zh.md)
