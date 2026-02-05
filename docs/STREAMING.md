# LingProxy 流式响应支持

LingProxy 现在支持流式响应（Streaming Responses），允许客户端实时接收 AI 生成的文本，提供更好的用户体验。

## 功能特性

- ✅ **Server-Sent Events (SSE)**: 使用标准的 SSE 格式传输流式数据
- ✅ **OpenAI 兼容**: 完全兼容 OpenAI API 的流式响应格式
- ✅ **实时输出**: 支持实时流式输出，减少首字延迟
- ✅ **自动统计**: 自动记录 token 使用量和请求统计

## 使用方法

### Python 客户端

```python
from openai import OpenAI

client = OpenAI(
    api_key="your-api-key",
    base_url="http://localhost:8080/llm/v1"
)

# 启用流式响应
stream = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[
        {"role": "user", "content": "Hello!"}
    ],
    stream=True  # 启用流式响应
)

# 处理流式响应
for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

### JavaScript/TypeScript 客户端

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: 'your-api-key',
  baseURL: 'http://localhost:8080/llm/v1'
});

const stream = await client.chat.completions.create({
  model: 'gpt-3.5-turbo',
  messages: [{ role: 'user', content: 'Hello!' }],
  stream: true
});

for await (const chunk of stream) {
  if (chunk.choices[0].delta.content) {
    process.stdout.write(chunk.choices[0].delta.content);
  }
}
```

### Go 客户端

```go
import (
    "context"
    "github.com/openai/openai-go/v3"
    "github.com/openai/openai-go/v3/option"
)

client := openai.NewClient(
    option.WithAPIKey("your-api-key"),
    option.WithBaseURL("http://localhost:8080/llm/v1"),
)

stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
    Model: "gpt-3.5-turbo",
    Messages: []openai.ChatCompletionMessageParamUnion{
        openai.UserMessage("Hello!"),
    },
})

for stream.Next() {
    chunk := stream.Current()
    // 处理 chunk
}
```

## API 格式

### 请求

```json
{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ],
  "stream": true
}
```

### 响应格式（SSE）

流式响应使用 Server-Sent Events (SSE) 格式：

```
data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":"stop"}]}

data: [DONE]
```

## 实现细节

### 后端实现

1. **Handler 层** (`backend/internal/handler/openai.go`):
   - 检查请求中的 `stream` 参数
   - 如果为 `true`，调用流式处理方法
   - 设置 SSE 响应头
   - 将流式数据转换为 SSE 格式发送

2. **Service 层** (`backend/internal/service/openai_service.go`):
   - `CreateChatCompletionStream`: 创建流式聊天补全请求
   - 调用 OpenAI SDK 的 `NewStreaming` 方法

3. **客户端层** (`backend/internal/client/openai/client.go`):
   - 封装 OpenAI SDK 的流式响应接口

### 响应头

流式响应会设置以下 HTTP 响应头：

```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
X-Accel-Buffering: no
```

## 注意事项

1. **性能**: 流式响应可以减少首字延迟，提供更好的用户体验
2. **网络**: 确保客户端和服务器之间的网络连接稳定
3. **超时**: 长时间运行的流式响应可能会遇到超时问题，需要适当配置超时时间
4. **错误处理**: 流式响应中的错误会通过 SSE 事件发送

## 故障排查

### 流式响应不工作

1. 检查请求中的 `stream` 参数是否为 `true`
2. 确认客户端支持 SSE（Server-Sent Events）
3. 检查网络连接是否稳定
4. 查看服务器日志以获取详细错误信息

### 响应中断

1. 检查网络连接稳定性
2. 确认服务器没有超时
3. 检查 LLM 资源是否正常工作

## 更新日志

- **v1.4.0** (2026-02-05): 添加流式响应支持
  - 实现 SSE 格式的流式响应
  - 支持实时流式输出
  - 自动记录流式请求统计
