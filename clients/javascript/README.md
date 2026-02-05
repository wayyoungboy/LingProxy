# LingProxy JavaScript/TypeScript Client

标准 JavaScript/TypeScript 客户端，用于与 LingProxy AI API 网关交互。

## 安装

```bash
npm install
```

或者直接安装 OpenAI SDK：

```bash
npm install openai
```

## 使用方法

### ES Modules

```javascript
import { LingProxyClient } from './lingproxy-client.js';

const client = new LingProxyClient({
  apiKey: 'your-api-key',
  baseURL: 'http://localhost:8080/llm/v1'
});

const response = await client.createChatCompletion({
  model: 'gpt-3.5-turbo',
  messages: [
    { role: 'user', content: 'Hello!' }
  ]
});

console.log(response.choices[0].message.content);
```

### CommonJS

```javascript
const { LingProxyClient } = require('./lingproxy-client.js');

const client = new LingProxyClient({
  apiKey: 'your-api-key',
  baseURL: 'http://localhost:8080/llm/v1'
});
```

### 使用环境变量

```bash
export LINGPROXY_API_KEY=your-api-key
```

```javascript
import { LingProxyClient } from './lingproxy-client.js';

// 自动从环境变量读取 API key
const client = new LingProxyClient();
```

### 使用 OpenAI SDK 风格

客户端完全兼容 OpenAI JavaScript SDK：

```javascript
import { LingProxyClient } from './lingproxy-client.js';

const client = new LingProxyClient({ apiKey: 'your-api-key' });

// 直接使用 OpenAI SDK 的方法
const response = await client.chat.completions.create({
  model: 'gpt-3.5-turbo',
  messages: [
    { role: 'user', content: 'Hello!' }
  ]
});

console.log(response.choices[0].message.content);
```

### 列出可用模型

```javascript
const models = await client.listModels();
models.forEach(model => {
  console.log(model.id);
});
```

### 文本补全

```javascript
const response = await client.createCompletion({
  model: 'gpt-3.5-turbo',
  prompt: 'Write a haiku about programming:',
  max_tokens: 50
});

console.log(response.choices[0].text);
```

### 流式响应

```javascript
const stream = await client.createChatCompletion({
  model: 'gpt-3.5-turbo',
  messages: [{ role: 'user', content: 'Tell me a story' }],
  stream: true
});

for await (const chunk of stream) {
  if (chunk.choices[0].delta.content) {
    process.stdout.write(chunk.choices[0].delta.content);
  }
}
```

## 运行示例

```bash
# 设置 API key
export LINGPROXY_API_KEY=your-api-key

# 运行示例
npm run demo
# 或
node lingproxy-client.js
```

## TypeScript 支持

客户端完全兼容 TypeScript，你可以直接使用类型定义：

```typescript
import { LingProxyClient } from './lingproxy-client.js';

const client = new LingProxyClient({
  apiKey: 'your-api-key',
  baseURL: 'http://localhost:8080/llm/v1'
});

const response = await client.createChatCompletion({
  model: 'gpt-3.5-turbo',
  messages: [
    { role: 'user', content: 'Hello!' }
  ]
});
```

## API 参考

### LingProxyClient

#### `constructor(options?)`

初始化客户端。

- `options.apiKey`: API 密钥或 Token。如果不提供，将从 `LINGPROXY_API_KEY` 环境变量读取
- `options.baseURL`: LingProxy 服务器的基础 URL（默认: http://localhost:8080/llm/v1）
- `options.timeout`: 请求超时时间（毫秒，默认: 30000）

#### `listModels(): Promise<Array>`

列出所有可用的模型。

#### `createChatCompletion(params): Promise<Object>`

创建聊天补全。

- `params.model`: 模型名称
- `params.messages`: 消息列表，每个消息包含 `role` 和 `content`
- `params.temperature`: 采样温度（0-2）
- `params.max_tokens`: 最大生成 token 数
- `params.top_p`: Nucleus 采样参数
- `params.stream`: 是否流式返回

#### `createCompletion(params): Promise<Object>`

创建文本补全。

- `params.model`: 模型名称
- `params.prompt`: 文本提示
- `params.temperature`: 采样温度（0-2）
- `params.max_tokens`: 最大生成 token 数

## 兼容性

本客户端基于 OpenAI JavaScript SDK，完全兼容 OpenAI API 规范。你可以直接使用 OpenAI SDK 的所有功能。
