# API Reference

## Base URL

- Development: `http://localhost:8080`
- Production: `https://your-domain.com`

## Authentication

Most APIs require authentication. Include the API key in the Authorization header:

```
Authorization: Bearer YOUR_API_KEY
```

## OpenAI-Compatible API

### Chat Completions

Create a chat completion. Supports streaming responses (Server-Sent Events).

**Endpoint:** `POST /llm/v1/chat/completions`

**Request Parameters:**
- `model` (string, required): Model name
- `messages` (array, required): Array of messages
- `temperature` (number, optional): Sampling temperature, between 0 and 2
- `max_tokens` (integer, optional): Maximum tokens to generate
- `stream` (boolean, optional): Enable streaming response, defaults to `false`

**Non-streaming Request Example:**
```json
{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello!"}
  ],
  "temperature": 0.7,
  "max_tokens": 100
}
```

**Non-streaming Response:**
```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "gpt-3.5-turbo",
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "Hello! How can I help you today?"
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 10,
    "total_tokens": 20
  }
}
```

**Streaming Request Example:**
```json
{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ],
  "stream": true
}
```

**Streaming Response (SSE Format):**
```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":"stop"}]}

data: [DONE]
```

**Streaming Response Notes:**
- Response headers are set to `Content-Type: text/event-stream`
- Each data chunk starts with `data: ` followed by a JSON object
- Streaming response ends with `data: [DONE]`
- Each chunk contains a `delta` field with incremental content

**Client Examples:**

Python:
```python
from openai import OpenAI

client = OpenAI(
    api_key="your-api-key",
    base_url="http://localhost:8080/llm/v1"
)

stream = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": "Hello!"}],
    stream=True
)

for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

JavaScript:
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

**Notes:**
- Streaming reduces first token latency and provides better user experience
- Ensure stable network connection between client and server
- Configure appropriate timeout for long-running streams
- Errors in streaming responses are sent via SSE events

### List Models

List available models.

**Endpoint:** `GET /llm/v1/models`

**Response:**
```json
{
  "data": [
    {
      "id": "gpt-4",
      "object": "model",
      "created": 1677610602,
      "owned_by": "openai"
    }
  ]
}
```

## Management API

### Admin APIs

#### Get Admin Info

**Endpoint:** `GET /api/v1/admin/info`

**Response:**
```json
{
  "data": {
    "id": "admin-id",
    "username": "admin",
    "api_key": "sk-...",
    "role": "admin",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Update Admin Info

**Endpoint:** `PUT /api/v1/admin/info`

**Request:**
```json
{
  "password": "current_password",
  "new_username": "new_admin",
  "new_password": "new_password123"
}
```

#### Update Admin Username

**Endpoint:** `PUT /api/v1/admin/username`

**Request:**
```json
{
  "password": "current_password",
  "username": "new_username"
}
```

#### Update Admin Password

**Endpoint:** `PUT /api/v1/admin/password`

**Request:**
```json
{
  "old_password": "old_password",
  "new_password": "new_password123"
}
```

#### Reset API Key

**Endpoint:** `PUT /api/v1/admin/api-key`

**Response:**
```json
{
  "message": "API key reset successfully",
  "data": {
    "api_key": "sk-new-api-key"
  }
}
```

### LLM Resources

#### List LLM Resources

**Endpoint:** `GET /api/v1/llm-resources`

**Query Parameters:**
- `search`: Search keyword (fuzzy search on name, base_url, model)

**Response:**
```json
{
  "data": [
    {
      "id": "resource-id",
      "name": "OpenAI GPT-4",
      "type": "chat",
      "driver": "openai",
      "model": "gpt-4",
      "base_url": "https://api.openai.com/v1",
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### Create LLM Resource

**Endpoint:** `POST /api/v1/llm-resources`

**Request:**
```json
{
  "name": "OpenAI GPT-4",
  "type": "chat",
  "driver": "openai",
  "model": "gpt-4",
  "base_url": "https://api.openai.com/v1",
  "api_key": "sk-...",
  "status": "active"
}
```

#### Update LLM Resource

**Endpoint:** `PUT /api/v1/llm-resources/:id`

#### Delete LLM Resource

**Endpoint:** `DELETE /api/v1/llm-resources/:id`

#### Import LLM Resources

**Endpoint:** `POST /api/v1/llm-resources/import`

**Request:** 
- Excel: Multipart form data with Excel file
- JSON: `Content-Type: application/json` with JSON array in request body

**JSON Request Example:**
```json
[
  {
    "name": "OpenAI GPT-4",
    "type": "chat",
    "driver": "openai",
    "model": "gpt-4",
    "base_url": "https://api.openai.com/v1",
    "api_key": "sk-...",
    "status": "active"
  }
]
```

**Response:**
```json
{
  "message": "导入完成",
  "success": 3,
  "fail": 0,
  "duplicate": 2,
  "total": 5,
  "errors": [],
  "duplicates": [
    {
      "row": 1,
      "name": "OpenAI GPT-4",
      "type": "chat",
      "model": "gpt-4",
      "base_url": "https://api.openai.com/v1"
    }
  ]
}
```

**Notes:**
- Duplicate detection: Resources with same `type`, `model`, `base_url`, and `api_key` are considered duplicates
- Automatic trimming: Leading and trailing whitespace are removed from all fields during import
- Duplicate resources are not imported and are reported in the response

#### Download Import Template

**Endpoint:** `GET /api/v1/llm-resources/import/template`

**Response:** Excel file download

#### Test LLM Resource

**Endpoint:** `POST /api/v1/llm-resources/:id/test`

**Description:** Test if an LLM resource can be called successfully. Only resources with `active` status can be tested.

**Response (Success):**
```json
{
  "success": true,
  "message": "测试成功",
  "model": "THUDM/GLM-4-9B-0414",
  "response": "\nHello 👋! I'm ChatGLM",
  "usage": {
    "prompt_tokens": 6,
    "completion_tokens": 10,
    "total_tokens": 16
  },
  "duration_ms": 529
}
```

**Response (Failure):**
```json
{
  "success": false,
  "error": "context deadline exceeded",
  "message": "测试失败: context deadline exceeded",
  "duration_ms": 30000
}
```

**Supported Resource Types:**
- `chat`: Tests by sending a simple "Hello" message with MaxTokens=10
- `embedding`: Tests by embedding the text "test"
- `rerank`: Currently not implemented

**Notes:**
- Test timeout is 30 seconds
- For chat resources, MaxTokens is limited to 10 to minimize costs
- Returns detailed information including model, response content, token usage, and duration

### API Keys

#### List API Keys

**Endpoint:** `GET /api/v1/api-keys`

**Note:** The old endpoint `GET /api/v1/tokens` is still available for backward compatibility but is deprecated.

#### Create API Key

**Endpoint:** `POST /api/v1/api-keys`

**Note:** The old endpoint `POST /api/v1/tokens` is still available for backward compatibility but is deprecated.

**Request:**
```json
{
  "name": "My API Key",
  "policy_id": "policy-id"
}
```

**Response:**
```json
{
  "data": {
    "id": "api-key-id",
    "name": "My API Key",
    "api_key": "ling-...",
    "prefix": "ling-...",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Update API Key

**Endpoint:** `PUT /api/v1/api-keys/:id`

**Note:** The old endpoint `PUT /api/v1/tokens/:id` is still available for backward compatibility but is deprecated.

#### Delete API Key

**Endpoint:** `DELETE /api/v1/api-keys/:id`

**Note:** The old endpoint `DELETE /api/v1/tokens/:id` is still available for backward compatibility but is deprecated.

#### Reset API Key

**Endpoint:** `POST /api/v1/api-keys/:id/reset`

**Note:** The old endpoint `POST /api/v1/tokens/:id/reset` is still available for backward compatibility but is deprecated.

### Policies

#### List Policies

**Endpoint:** `GET /api/v1/policies`

#### Create Policy

**Endpoint:** `POST /api/v1/policies`

**Request:**
```json
{
  "name": "My Policy",
  "template_id": "template-id",
  "type": "round_robin",
  "parameters": "{}",
  "enabled": true
}
```

### Statistics

#### Get System Statistics

**Endpoint:** `GET /api/v1/stats/system`

**Response:**
```json
{
  "data": {
    "total_requests": 1000,
    "total_users": 10,
    "total_llm_resources": 5,
    "success_rate": 98.5,
    "avg_response_time": 120.5
  }
}
```

#### Get LLM Resource Usage Statistics

**Endpoint:** `GET /api/v1/stats/llm-resources/usage`

**Description:** Get usage statistics grouped by LLM resources, including token usage, request count, success rate, etc.

**Response:**
```json
{
  "data": [
    {
      "resource_id": "1770294403900247000",
      "resource_name": "SiliconFlow-Chat-Qwen-Qwen2.5-7B-Instruct",
      "resource_type": "chat",
      "model": "Qwen/Qwen2.5-7B-Instruct",
      "total_tokens": 50000,
      "total_requests": 100,
      "success_requests": 98,
      "failed_requests": 2,
      "success_rate": 98.0,
      "avg_tokens_per_request": 500,
      "last_request_time": "2026-02-06T00:29:27+08:00"
    }
  ]
}
```

**Field Descriptions:**
- `resource_id`: LLM resource ID
- `resource_name`: Resource name
- `resource_type`: Resource type (chat, image, embedding, etc.)
- `model`: Model identifier
- `total_tokens`: Total token usage
- `total_requests`: Total request count
- `success_requests`: Successful request count
- `failed_requests`: Failed request count
- `success_rate`: Success rate (percentage)
- `avg_tokens_per_request`: Average tokens per request
- `last_request_time`: Last request time

#### Get Single LLM Resource Statistics

**Endpoint:** `GET /api/v1/stats/llm-resources/:id`

**Parameters:**
- `id`: LLM resource ID

**Response:**
```json
{
  "data": {
    "llm_resource_id": "1770294403900247000",
    "total_requests": 100,
    "success_rate": 98.0,
    "avg_response_time": 110
  }
}
```

#### Get User Statistics

**Endpoint:** `GET /api/v1/stats/users/:id`

**Parameters:**
- `id`: User ID

**Response:**
```json
{
  "data": {
    "user_id": "1770272193110467000",
    "total_requests": 100,
    "total_tokens": 50000,
    "success_rate": 97.8,
    "avg_response_time": 130
  }
}
```

### System

#### Get System Info

**Endpoint:** `GET /api/v1/system/info`

**Response:**
```json
{
  "data": {
    "cpu_usage": 25.5,
    "memory_usage": 60.2,
    "uptime": 3600,
    "version": "1.0.0"
  }
}
```

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message description"
}
```

**HTTP Status Codes:**
- `200`: Success
- `400`: Bad Request
- `401`: Unauthorized
- `404`: Not Found
- `500`: Internal Server Error

## Automatic Retry

LingProxy automatically retries failed requests based on configurable settings:

**Retry Behavior:**
- Failed requests are automatically retried for network errors, timeouts, and 5xx server errors
- Uses exponential backoff: delay increases with each retry attempt
- Maximum retry count and delay are configurable via admin interface (Settings → Provider Settings)
- Default: 3 retries with 1 second base delay

**Retryable Errors:**
- Network connection failures
- Request timeouts
- 5xx server errors (500, 502, 503, 504)
- 429 rate limit errors

**Non-Retryable Errors:**
- 4xx client errors (except 429)
- Authentication failures (401, 403)
- Invalid request parameters

**Note:** Retry configuration changes take effect immediately without requiring a service restart.

## Swagger Documentation

Interactive API documentation is available at:
- Swagger UI: `http://localhost:8080/swagger/index.html`
