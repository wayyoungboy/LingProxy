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
    // 对于blob响应，直接返回response对象
    if (response.config.responseType === 'blob') {
      return response.data
    }
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

  // 认证相关
  register(userData) {
    return apiClient.post('/auth/register', userData)
  },
  login(credentials) {
    return apiClient.post('/auth/login', credentials)
  },
  getCurrentUser() {
    return apiClient.get('/users/me')
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
  resetAPIKey(id) {
    return apiClient.post(`/users/${id}/reset-api-key`)
  },
  updatePassword(id, passwordData) {
    return apiClient.put(`/users/${id}/password`, passwordData)
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
  importLLMResources(file) {
    const formData = new FormData()
    formData.append('file', file)
    return apiClient.post('/llm-resources/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },
  downloadLLMResourcesTemplate() {
    return apiClient.get('/llm-resources/import/template', {
      responseType: 'blob'
    })
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

  // Token管理
  getTokens() {
    return apiClient.get('/tokens')
  },
  getToken(id) {
    return apiClient.get(`/tokens/${id}`)
  },
  createToken(token) {
    return apiClient.post('/tokens', token)
  },
  updateToken(id, token) {
    return apiClient.put(`/tokens/${id}`, token)
  },
  deleteToken(id) {
    return apiClient.delete(`/tokens/${id}`)
  },
  resetToken(id) {
    return apiClient.post(`/tokens/${id}/reset`)
  },
  setTokenPolicy(id, data) {
    return apiClient.put(`/tokens/${id}/policy`, data)
  },
  removeTokenPolicy(id) {
    return apiClient.delete(`/tokens/${id}/policy`)
  },

  // 管理员
  getAdminInfo() {
    return apiClient.get('/admin/info')
  },
  resetAdminAPIKey() {
    return apiClient.put('/admin/api-key')
  },
  updateAdminPassword(passwordData) {
    return apiClient.put('/admin/password', passwordData)
  },

  // 策略管理
  getPolicyTemplates() {
    return apiClient.get('/policy-templates')
  },
  getPolicyTemplate(id) {
    return apiClient.get(`/policy-templates/${id}`)
  },
  getPolicies() {
    return apiClient.get('/policies')
  },
  getPolicy(id) {
    return apiClient.get(`/policies/${id}`)
  },
  createPolicy(policy) {
    return apiClient.post('/policies', policy)
  },
  updatePolicy(id, policy) {
    return apiClient.put(`/policies/${id}`, policy)
  },
  deletePolicy(id) {
    return apiClient.delete(`/policies/${id}`)
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
  },

  // 日志管理
  getLogFiles() {
    return apiClient.get('/logs/files')
  },
  getLogs(params) {
    return apiClient.get('/logs', { params })
  },
  downloadLogFile(fileName) {
    return apiClient.get(`/logs/files/${fileName}/download`, {
      responseType: 'blob'
    })
  },
  clearLogs(fileName) {
    return apiClient.post('/logs/clear', null, { params: { file: fileName } })
  }
}

export default api