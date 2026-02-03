# API Reference

## Base URL

- Development: `http://localhost:8080`
- Production: `https://your-domain.com`

## Authentication

Most APIs require authentication. Include the API key or token in the Authorization header:

```
Authorization: Bearer YOUR_API_KEY_OR_TOKEN
```

## OpenAI-Compatible API

### Chat Completions

Create a chat completion.

**Endpoint:** `POST /llm/v1/chat/completions`

**Request:**
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

**Response:**
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

**Request:** Multipart form data with Excel file

**Response:**
```json
{
  "message": "Import completed",
  "success": 3,
  "failed": 0
}
```

#### Download Import Template

**Endpoint:** `GET /api/v1/llm-resources/import/template`

**Response:** Excel file download

### Tokens

#### List Tokens

**Endpoint:** `GET /api/v1/tokens`

#### Create Token

**Endpoint:** `POST /api/v1/tokens`

**Request:**
```json
{
  "name": "My Token",
  "policy_id": "policy-id"
}
```

#### Update Token

**Endpoint:** `PUT /api/v1/tokens/:id`

#### Delete Token

**Endpoint:** `DELETE /api/v1/tokens/:id`

#### Reset Token

**Endpoint:** `POST /api/v1/tokens/:id/reset`

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

## Swagger Documentation

Interactive API documentation is available at:
- Swagger UI: `http://localhost:8080/swagger/index.html`
