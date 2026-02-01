import axios from 'axios'

// 创建axios实例
const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
apiClient.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    // 处理401错误
    if (error.response && error.response.status === 401) {
      // 清除token并跳转到登录页
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// API方法
const api = {
  // 健康检查
  healthCheck() {
    return apiClient.get('/health')
  },

  // 用户管理
  getUsers(params) {
    return apiClient.get('/users', { params })
  },
  getUser(id) {
    return apiClient.get(`/users/${id}`)
  },
  createUser(user) {
    return apiClient.post('/users', user)
  },
  updateUser(id, user) {
    return apiClient.put(`/users/${id}`, user)
  },
  deleteUser(id) {
    return apiClient.delete(`/users/${id}`)
  },

  // LLM资源管理
  getLLMResources(params) {
    return apiClient.get('/llm-resources', { params })
  },
  getLLMResource(id) {
    return apiClient.get(`/llm-resources/${id}`)
  },
  createLLMResource(resource) {
    return apiClient.post('/llm-resources', resource)
  },
  updateLLMResource(id, resource) {
    return apiClient.put(`/llm-resources/${id}`, resource)
  },
  deleteLLMResource(id) {
    return apiClient.delete(`/llm-resources/${id}`)
  },
  getLLMResourceModels(resourceId) {
    return apiClient.get(`/llm-resources/${resourceId}/models`)
  },

  // 模型管理
  getModels(params) {
    return apiClient.get('/models', { params })
  },
  getModel(id) {
    return apiClient.get(`/models/${id}`)
  },
  createModel(model) {
    return apiClient.post('/models', model)
  },
  updateModel(id, model) {
    return apiClient.put(`/models/${id}`, model)
  },
  deleteModel(id) {
    return apiClient.delete(`/models/${id}`)
  },
  getModelPricing(id) {
    return apiClient.get(`/models/${id}/pricing`)
  },
  getModelTypes() {
    return apiClient.get('/models/types')
  },

  // 端点管理
  getEndpoints(params) {
    return apiClient.get('/endpoints', { params })
  },
  getEndpoint(id) {
    return apiClient.get(`/endpoints/${id}`)
  },
  createEndpoint(endpoint) {
    return apiClient.post('/endpoints', endpoint)
  },
  updateEndpoint(id, endpoint) {
    return apiClient.put(`/endpoints/${id}`, endpoint)
  },
  deleteEndpoint(id) {
    return apiClient.delete(`/endpoints/${id}`)
  },

  // 请求管理
  getRequests(params) {
    return apiClient.get('/requests', { params })
  },
  getRequestDetail(id) {
    return apiClient.get(`/requests/${id}`)
  },
  exportRequests(params) {
    return apiClient.get('/requests/export', { params })
  },

  // 统计信息
  getSystemStats() {
    return apiClient.get('/stats/system')
  },
  getLLMResourceStats(id) {
    return apiClient.get(`/stats/llm-resources/${id}`)
  },
  getUserStats(id) {
    return apiClient.get(`/stats/users/${id}`)
  },

  // 系统设置
  getSettings() {
    return apiClient.get('/settings')
  },
  updateSettings(settings) {
    return apiClient.put('/settings', settings)
  },
  getSystemInfo() {
    return apiClient.get('/system/info')
  }
}

export default api