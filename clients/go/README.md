# LingProxy Go Client

标准 Go 客户端库，用于与 LingProxy AI API 网关交互。

## 安装

```bash
go get github.com/lingproxy/lingproxy/clients/go
```

## 使用方法

### 基本用法

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
        APIKey:  "your-api-key",
        BaseURL: "http://localhost:8080/llm/v1",
    })
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    // 创建聊天补全
    temp := 0.7
    maxTokens := int64(100)
    req := client.ChatCompletionRequest{
        Model: "gpt-3.5-turbo",
        Messages: []client.ChatMessage{
            {Role: "user", Content: "Hello!"},
        },
        Temperature: &temp,
        MaxTokens:   &maxTokens,
    }
    
    resp, err := c.CreateChatCompletion(ctx, req)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(resp.Choices[0].Message.Content)
}
```

### 使用环境变量

```bash
export LINGPROXY_API_KEY=your-api-key
```

```go
c, err := client.NewClient(&client.ClientOptions{})
```

### 使用 OpenAI SDK 风格

客户端提供了对底层 OpenAI SDK 的直接访问：

```go
openaiClient := c.GetClient()

response, err := openaiClient.Chat.Completions.New(ctx, params)
```

### 列出可用模型

```go
models, err := c.ListModels(ctx)
if err != nil {
    panic(err)
}

for _, model := range models {
    fmt.Println(model.ID)
}
```

### 文本补全

```go
req := client.CompletionRequest{
    Model:  "gpt-3.5-turbo",
    Prompt: "Write a haiku about programming:",
}

resp, err := c.CreateCompletion(ctx, req)
if err != nil {
    panic(err)
}

fmt.Println(resp.Choices[0].Text)
```

## 运行示例

```bash
# 设置 API key
export LINGPROXY_API_KEY=your-api-key

# 运行示例
go run example.go
```

## API 参考

### NewClient(options *ClientOptions) (*LingProxyClient, error)

创建新的 LingProxy 客户端。

- `options.APIKey`: API 密钥或 Token。如果不提供，将从 `LINGPROXY_API_KEY` 环境变量读取
- `options.BaseURL`: LingProxy 服务器的基础 URL（默认: http://localhost:8080/llm/v1）
- `options.Timeout`: 请求超时时间（默认: 30秒）

### ListModels(ctx context.Context) ([]openai.Model, error)

列出所有可用的模型。

### CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)

创建聊天补全。

- `req.Model`: 模型名称
- `req.Messages`: 消息列表，每个消息包含 `Role` 和 `Content`
- `req.Temperature`: 采样温度（0-2）
- `req.MaxTokens`: 最大生成 token 数
- `req.TopP`: Nucleus 采样参数

### CreateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

创建文本补全。

- `req.Model`: 模型名称
- `req.Prompt`: 文本提示
- `req.Temperature`: 采样温度（0-2）
- `req.MaxTokens`: 最大生成 token 数

### GetClient() *openai.Client

返回底层 OpenAI 客户端，用于直接访问 OpenAI SDK 的所有功能。

## 兼容性

本客户端基于 OpenAI Go SDK，完全兼容 OpenAI API 规范。你可以通过 `GetClient()` 方法直接使用 OpenAI SDK 的所有功能。
