# LingProxy Python Client

标准 Python 客户端，用于与 LingProxy AI API 网关交互。

## 安装

```bash
pip install -r requirements.txt
```

或者直接安装 OpenAI SDK：

```bash
pip install openai>=1.0.0
```

## 使用方法

### 基本用法

```python
from lingproxy_client import LingProxyClient

# 初始化客户端
client = LingProxyClient(
    api_key="your-api-key",
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

### 使用环境变量

```bash
export LINGPROXY_API_KEY=your-api-key
```

```python
from lingproxy_client import LingProxyClient

# 自动从环境变量读取 API key
client = LingProxyClient()
```

### 使用 OpenAI SDK 风格

客户端完全兼容 OpenAI Python SDK：

```python
from lingproxy_client import LingProxyClient

client = LingProxyClient(api_key="your-api-key")

# 直接使用 OpenAI SDK 的方法
response = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)

print(response.choices[0].message.content)
```

### 列出可用模型

```python
models = client.list_models()
for model in models:
    print(model['id'])
```

### 文本补全

```python
response = client.create_completion(
    model="gpt-3.5-turbo",
    prompt="Write a haiku about programming:",
    max_tokens=50
)

print(response['choices'][0]['text'])
```

### 流式响应

```python
stream = client.create_chat_completion(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": "Tell me a story"}],
    stream=True
)

for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

## 运行示例

```bash
# 设置 API key
export LINGPROXY_API_KEY=your-api-key

# 运行示例
python lingproxy_client.py
```

## API 参考

### LingProxyClient

#### `__init__(api_key=None, base_url="http://localhost:8080/llm/v1", timeout=30.0)`

初始化客户端。

- `api_key`: API 密钥或 Token。如果不提供，将从 `LINGPROXY_API_KEY` 环境变量读取
- `base_url`: LingProxy 服务器的基础 URL
- `timeout`: 请求超时时间（秒）

#### `list_models() -> List[Dict]`

列出所有可用的模型。

#### `create_chat_completion(model, messages, **kwargs) -> Dict`

创建聊天补全。

- `model`: 模型名称
- `messages`: 消息列表，每个消息包含 `role` 和 `content`
- `temperature`: 采样温度（0-2）
- `max_tokens`: 最大生成 token 数
- `top_p`: Nucleus 采样参数
- `stream`: 是否流式返回

#### `create_completion(model, prompt, **kwargs) -> Dict`

创建文本补全。

- `model`: 模型名称
- `prompt`: 文本提示
- `temperature`: 采样温度（0-2）
- `max_tokens`: 最大生成 token 数

## 兼容性

本客户端基于 OpenAI Python SDK，完全兼容 OpenAI API 规范。你可以直接使用 OpenAI SDK 的所有功能。
