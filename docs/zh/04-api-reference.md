# API 参考文档

## 基础 URL

- 开发环境：`http://localhost:8080`
- 生产环境：`https://your-domain.com`

## 认证

大多数 API 需要认证。在 Authorization 头中包含 API 密钥或 Token：

```
Authorization: Bearer YOUR_API_KEY_OR_TOKEN
```

## OpenAI 兼容 API

### 聊天补全

创建聊天补全。

**端点：** `POST /llm/v1/chat/completions`

**请求：**
```json
{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "system", "content": "你是一个有用的助手。"},
    {"role": "user", "content": "你好！"}
  ],
  "temperature": 0.7,
  "max_tokens": 100
}
```

**响应：**
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
      "content": "你好！今天我能为你做些什么？"
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

### 列出模型

列出可用模型。

**端点：** `GET /llm/v1/models`

**响应：**
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

## 管理 API

### 管理员 API

#### 获取管理员信息

**端点：** `GET /api/v1/admin/info`

**响应：**
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

#### 更新管理员信息

**端点：** `PUT /api/v1/admin/info`

**请求：**
```json
{
  "password": "当前密码",
  "new_username": "新管理员名",
  "new_password": "新密码123"
}
```

#### 更新管理员用户名

**端点：** `PUT /api/v1/admin/username`

**请求：**
```json
{
  "password": "当前密码",
  "username": "新用户名"
}
```

#### 更新管理员密码

**端点：** `PUT /api/v1/admin/password`

**请求：**
```json
{
  "old_password": "旧密码",
  "new_password": "新密码123"
}
```

#### 重置 API 密钥

**端点：** `PUT /api/v1/admin/api-key`

**响应：**
```json
{
  "message": "API key reset successfully",
  "data": {
    "api_key": "sk-new-api-key"
  }
}
```

### LLM 资源

#### 列出 LLM 资源

**端点：** `GET /api/v1/llm-resources`

**查询参数：**
- `search`: 搜索关键词（对名称、base_url、模型进行模糊搜索）

**响应：**
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

#### 创建 LLM 资源

**端点：** `POST /api/v1/llm-resources`

**请求：**
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

#### 更新 LLM 资源

**端点：** `PUT /api/v1/llm-resources/:id`

#### 删除 LLM 资源

**端点：** `DELETE /api/v1/llm-resources/:id`

#### 导入 LLM 资源

**端点：** `POST /api/v1/llm-resources/import`

**请求：** 包含 Excel 文件的多部分表单数据

**响应：**
```json
{
  "message": "Import completed",
  "success": 3,
  "failed": 0
}
```

#### 下载导入模板

**端点：** `GET /api/v1/llm-resources/import/template`

**响应：** Excel 文件下载

### Token

#### 列出 Token

**端点：** `GET /api/v1/tokens`

#### 创建 Token

**端点：** `POST /api/v1/tokens`

**请求：**
```json
{
  "name": "我的 Token",
  "policy_id": "policy-id"
}
```

#### 更新 Token

**端点：** `PUT /api/v1/tokens/:id`

#### 删除 Token

**端点：** `DELETE /api/v1/tokens/:id`

#### 重置 Token

**端点：** `POST /api/v1/tokens/:id/reset`

### 策略

#### 列出策略

**端点：** `GET /api/v1/policies`

#### 创建策略

**端点：** `POST /api/v1/policies`

**请求：**
```json
{
  "name": "我的策略",
  "template_id": "template-id",
  "type": "round_robin",
  "parameters": "{}",
  "enabled": true
}
```

### 系统

#### 获取系统信息

**端点：** `GET /api/v1/system/info`

**响应：**
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

## 错误响应

所有错误遵循以下格式：

```json
{
  "error": "错误消息描述"
}
```

**HTTP 状态码：**
- `200`: 成功
- `400`: 错误请求
- `401`: 未授权
- `404`: 未找到
- `500`: 内部服务器错误

## Swagger 文档

交互式 API 文档可在以下位置访问：
- Swagger UI: `http://localhost:8080/swagger/index.html`
