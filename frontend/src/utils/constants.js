/**
 * 常量定义
 */

// API相关常量
export const API_TIMEOUT = 10000
export const API_BASE_URL = '/api/v1'

// 存储键名
export const STORAGE_KEYS = {
  TOKEN: 'token',
  USER_INFO: 'userInfo',
  THEME: 'theme',
  LANGUAGE: 'language'
}

// 路由路径
export const ROUTE_PATHS = {
  LOGIN: '/login',
  DASHBOARD: '/dashboard',
  TOKENS: '/tokens',
  LLM_RESOURCES: '/llm-resources',
  LLM_RESOURCE_USAGE: '/llm-resources/usage',
  MODELS: '/models',
  REQUESTS: '/requests',
  POLICIES: '/policies',
  SETTINGS: '/settings',
  LOGS: '/logs',
  USERS: '/users',
  ENDPOINTS: '/endpoints'
}

// 模型类型
export const MODEL_TYPES = {
  CHAT: 'chat',
  IMAGE: 'image',
  EMBEDDING: 'embedding',
  RERANK: 'rerank',
  AUDIO: 'audio',
  VIDEO: 'video'
}

// 模型类型标签映射
export const MODEL_TYPE_LABELS = {
  [MODEL_TYPES.CHAT]: '对话',
  [MODEL_TYPES.IMAGE]: '生图',
  [MODEL_TYPES.EMBEDDING]: '嵌入',
  [MODEL_TYPES.RERANK]: '重排序',
  [MODEL_TYPES.AUDIO]: '语音',
  [MODEL_TYPES.VIDEO]: '视频'
}

// 服务提供商
export const PROVIDERS = {
  OPENAI: 'openai',
  ANTHROPIC: 'anthropic',
  GOOGLE: 'google',
  COHERE: 'cohere',
  JINA: 'jina',
  QWEN: 'qwen',
  OLLAMA: 'ollama',
  CUSTOM: 'custom'
}

// 服务提供商标签映射
export const PROVIDER_LABELS = {
  [PROVIDERS.OPENAI]: 'OpenAI',
  [PROVIDERS.ANTHROPIC]: 'Anthropic',
  [PROVIDERS.GOOGLE]: 'Google',
  [PROVIDERS.COHERE]: 'Cohere',
  [PROVIDERS.JINA]: 'Jina',
  [PROVIDERS.QWEN]: 'Qwen',
  [PROVIDERS.OLLAMA]: 'Ollama',
  [PROVIDERS.CUSTOM]: '自定义'
}

// 状态
export const STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive',
  PENDING: 'pending'
}

// 状态标签映射
export const STATUS_LABELS = {
  [STATUS.ACTIVE]: '活跃',
  [STATUS.INACTIVE]: '禁用',
  [STATUS.PENDING]: '待处理'
}

// 请求状态
export const REQUEST_STATUS = {
  SUCCESS: 'success',
  FAILED: 'failed',
  PENDING: 'pending'
}

// 请求状态标签映射
export const REQUEST_STATUS_LABELS = {
  [REQUEST_STATUS.SUCCESS]: '成功',
  [REQUEST_STATUS.FAILED]: '失败',
  [REQUEST_STATUS.PENDING]: '处理中'
}

// 分页默认值
export const PAGINATION = {
  DEFAULT_PAGE_SIZE: 10,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100]
}

// 表格列宽
export const TABLE_COLUMN_WIDTHS = {
  ID: 180,
  STATUS: 100,
  ACTIONS: 200,
  DATE: 180,
  DURATION: 120
}

// 文件上传限制
export const FILE_UPLOAD = {
  MAX_SIZE: 10 * 1024 * 1024, // 10MB
  ACCEPTED_TYPES: ['.xlsx', '.xls', '.csv']
}

// 消息提示持续时间
export const MESSAGE_DURATION = {
  SUCCESS: 2000,
  ERROR: 3000,
  WARNING: 2500,
  INFO: 2000
}

// 菜单相关工具函数（从menu.js导入）
// 注意：getMenuTitle 现在需要传入 i18n 的 t 函数，请从 menu.js 直接导入使用
