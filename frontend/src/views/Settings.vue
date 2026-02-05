<template>
  <div class="settings-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="page-title">系统设置</span>
          <el-button type="primary" @click="saveSettings">
            <el-icon><Check /></el-icon>
            保存设置
          </el-button>
        </div>
      </template>
      
      <el-tabs v-model="activeTab" class="mt-4">
        <!-- 基本设置 -->
        <el-tab-pane label="基本设置" name="basic">
          <el-form :model="settingsForm.basic" label-width="150px">
            <el-form-item label="系统名称">
              <el-input v-model="settingsForm.basic.system_name" placeholder="请输入系统名称" />
            </el-form-item>
            <el-form-item label="系统版本">
              <el-input v-model="settingsForm.basic.system_version" disabled />
            </el-form-item>
            <el-form-item label="运行环境">
              <el-input v-model="settingsForm.basic.environment" disabled />
            </el-form-item>
            <el-form-item label="API地址">
              <el-input v-model="settingsForm.basic.api_url" placeholder="请输入API地址" />
            </el-form-item>
            <el-form-item label="服务端口">
              <el-input-number v-model="settingsForm.basic.port" :min="1" :max="65535" />
              <el-alert type="warning" :closable="false" style="margin-top: 10px">
                修改端口需要重启服务才能生效
              </el-alert>
            </el-form-item>
            <el-form-item label="监听地址">
              <el-input v-model="settingsForm.basic.host" placeholder="请输入监听地址" />
              <el-alert type="warning" :closable="false" style="margin-top: 10px">
                修改监听地址需要重启服务才能生效
              </el-alert>
            </el-form-item>
            <el-form-item label="调试模式">
              <el-switch v-model="settingsForm.basic.debug_mode" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 缓存设置 -->
        <el-tab-pane label="缓存设置" name="cache">
          <el-form :model="settingsForm.cache" label-width="150px">
            <el-form-item label="启用缓存">
              <el-switch v-model="settingsForm.cache.enabled" />
            </el-form-item>
            <el-form-item label="缓存类型">
              <el-select v-model="settingsForm.cache.type" placeholder="请选择缓存类型" disabled>
                <el-option label="内存缓存" value="memory" />
              </el-select>
            </el-form-item>
            <el-form-item label="缓存过期时间(秒)">
              <el-input-number v-model="settingsForm.cache.ttl" :min="1" :max="86400" />
            </el-form-item>
            <el-form-item label="缓存大小限制(MB)">
              <el-input-number v-model="settingsForm.cache.size_limit" :min="1" :max="1024" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 限流设置 -->
        <el-tab-pane label="限流设置" name="rateLimit">
          <el-form :model="settingsForm.rate_limit" label-width="150px">
            <el-form-item label="启用限流">
              <el-switch v-model="settingsForm.rate_limit.enabled" />
            </el-form-item>
            <el-form-item label="每分钟请求数限制">
              <el-input-number v-model="settingsForm.rate_limit.requests_per_minute" :min="1" :max="10000" />
            </el-form-item>
            <el-form-item label="并发请求限制">
              <el-input-number v-model="settingsForm.rate_limit.concurrency" :min="1" :max="1000" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 安全设置 -->
        <el-tab-pane label="安全设置" name="security">
          <el-form :model="settingsForm.security" label-width="150px">
            <el-form-item label="启用认证">
              <el-switch v-model="settingsForm.security.auth_enabled" />
              <el-alert type="warning" :closable="false" style="margin-top: 10px">
                关闭认证后，所有API（除登录外）都不需要认证即可访问。修改此设置需要重启服务。
              </el-alert>
            </el-form-item>
            <el-divider />
            <el-form-item label="JWT密钥">
              <el-input 
                v-model="settingsForm.security.jwt_secret" 
                type="password" 
                placeholder="留空则不修改（显示为******）"
                show-password
              />
              <el-alert type="warning" :closable="false" style="margin-top: 10px">
                修改JWT密钥需要重启服务，且会导致所有现有Token失效
              </el-alert>
            </el-form-item>
            <el-form-item label="Token过期时间(小时)">
              <el-input-number v-model="settingsForm.security.token_expiration" :min="1" :max="720" />
            </el-form-item>
            <el-form-item label="启用HTTPS">
              <el-switch v-model="settingsForm.security.https_enabled" disabled />
              <span style="color: #909399; margin-left: 10px">暂未实现</span>
            </el-form-item>
            <el-divider>管理员密码</el-divider>
            <el-form-item label="当前密码">
              <el-input
                v-model="passwordForm.oldPassword"
                type="password"
                placeholder="请输入当前密码"
                show-password
                style="width: 300px"
              />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input
                v-model="passwordForm.newPassword"
                type="password"
                placeholder="请输入新密码（至少6位）"
                show-password
                style="width: 300px"
              />
            </el-form-item>
            <el-form-item label="确认新密码">
              <el-input
                v-model="passwordForm.confirmPassword"
                type="password"
                placeholder="请再次输入新密码"
                show-password
                style="width: 300px"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="updateAdminPassword" :loading="passwordUpdating">
                更新密码
              </el-button>
              <el-button @click="resetPasswordForm">重置</el-button>
            </el-form-item>
            <el-divider>CORS设置</el-divider>
            <el-form-item label="启用CORS">
              <el-switch v-model="settingsForm.security.cors.enabled" />
            </el-form-item>
            <el-form-item label="允许的来源">
              <el-select 
                v-model="settingsForm.security.cors.allow_origins" 
                multiple 
                placeholder="选择允许的来源"
                style="width: 100%"
              >
                <el-option label="全部 (*)" value="*" />
              </el-select>
            </el-form-item>
            <el-form-item label="允许的方法">
              <el-select 
                v-model="settingsForm.security.cors.allow_methods" 
                multiple 
                placeholder="选择允许的HTTP方法"
                style="width: 100%"
              >
                <el-option label="GET" value="GET" />
                <el-option label="POST" value="POST" />
                <el-option label="PUT" value="PUT" />
                <el-option label="DELETE" value="DELETE" />
                <el-option label="OPTIONS" value="OPTIONS" />
              </el-select>
            </el-form-item>
            <el-form-item label="允许的请求头">
              <el-select 
                v-model="settingsForm.security.cors.allow_headers" 
                multiple 
                placeholder="选择允许的请求头"
                style="width: 100%"
              >
                <el-option label="全部 (*)" value="*" />
              </el-select>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 日志设置 -->
        <el-tab-pane label="日志设置" name="log">
          <el-form :model="settingsForm.log" label-width="150px">
            <el-form-item label="日志级别">
              <el-select v-model="settingsForm.log.level" placeholder="请选择日志级别">
                <el-option label="DEBUG" value="debug" />
                <el-option label="INFO" value="info" />
                <el-option label="WARN" value="warn" />
                <el-option label="ERROR" value="error" />
              </el-select>
            </el-form-item>
            <el-form-item label="日志格式">
              <el-select v-model="settingsForm.log.format" placeholder="请选择日志格式">
                <el-option label="JSON" value="json" />
                <el-option label="文本" value="text" />
              </el-select>
            </el-form-item>
            <el-form-item label="输出方式">
              <el-select v-model="settingsForm.log.output" placeholder="请选择输出方式">
                <el-option label="标准输出" value="stdout" />
                <el-option label="文件" value="file" />
                <el-option label="同时输出（推荐）" value="both" />
              </el-select>
            </el-form-item>
            <el-form-item label="日志文件路径" v-if="settingsForm.log.output === 'file' || settingsForm.log.output === 'both'">
              <el-input v-model="settingsForm.log.file_path" placeholder="请输入日志文件路径" />
            </el-form-item>
            <el-form-item label="单个文件最大大小(MB)">
              <el-input-number v-model="settingsForm.log.max_size" :min="1" :max="1024" />
            </el-form-item>
            <el-form-item label="保留备份数量">
              <el-input-number v-model="settingsForm.log.max_backups" :min="1" :max="100" />
            </el-form-item>
            <el-form-item label="保留天数">
              <el-input-number v-model="settingsForm.log.max_age" :min="1" :max="365" />
            </el-form-item>
            <el-form-item label="压缩旧日志">
              <el-switch v-model="settingsForm.log.compress" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 负载均衡设置 -->
        <el-tab-pane label="负载均衡" name="loadBalancer">
          <el-form :model="settingsForm.load_balancer" label-width="150px">
            <el-form-item label="默认策略">
              <el-select v-model="settingsForm.load_balancer.default_strategy" placeholder="请选择默认策略">
                <el-option label="轮询" value="round_robin" />
                <el-option label="最少连接" value="least_connections" />
                <el-option label="加权" value="weighted" />
              </el-select>
            </el-form-item>
            <el-divider>健康检查</el-divider>
            <el-form-item label="启用健康检查">
              <el-switch v-model="settingsForm.load_balancer.health_check.enabled" />
            </el-form-item>
            <el-form-item label="检查间隔(秒)">
              <el-input-number v-model="settingsForm.load_balancer.health_check.interval" :min="5" :max="300" />
            </el-form-item>
            <el-form-item label="超时时间(秒)">
              <el-input-number v-model="settingsForm.load_balancer.health_check.timeout" :min="1" :max="60" />
            </el-form-item>
            <el-form-item label="最大失败次数">
              <el-input-number v-model="settingsForm.load_balancer.health_check.max_failures" :min="1" :max="10" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 系统信息 -->
        <el-tab-pane label="系统信息" name="info">
          <el-descriptions :column="2" border class="mt-4">
            <el-descriptions-item label="系统名称">{{ systemInfo.system_name }}</el-descriptions-item>
            <el-descriptions-item label="系统版本">{{ systemInfo.system_version }}</el-descriptions-item>
            <el-descriptions-item label="运行环境">{{ systemInfo.environment }}</el-descriptions-item>
            <el-descriptions-item label="Go版本">{{ systemInfo.go_version }}</el-descriptions-item>
            <el-descriptions-item label="操作系统">{{ systemInfo.os }}</el-descriptions-item>
            <el-descriptions-item label="系统架构">{{ systemInfo.arch }}</el-descriptions-item>
            <el-descriptions-item label="启动时间">{{ formatTime(systemInfo.start_time) }}</el-descriptions-item>
            <el-descriptions-item label="当前时间">{{ formatTime(systemInfo.current_time) }}</el-descriptions-item>
            <el-descriptions-item label="运行时长">{{ systemInfo.uptime }}</el-descriptions-item>
            <el-descriptions-item label="CPU核心数">{{ systemInfo.cpu_cores }}</el-descriptions-item>
            <el-descriptions-item label="内存使用">
              {{ formatBytes(systemInfo.memory_used) }} / {{ formatBytes(systemInfo.memory_total) }} 
              ({{ formatPercentage(systemInfo.memory_usage) }})
            </el-descriptions-item>
            <el-descriptions-item label="磁盘使用" v-if="systemInfo.disk_total > 0">
              {{ formatBytes(systemInfo.disk_used) }} / {{ formatBytes(systemInfo.disk_total) }} 
              ({{ formatPercentage(systemInfo.disk_usage) }})
            </el-descriptions-item>
          </el-descriptions>
          <el-button type="primary" class="mt-4" @click="refreshSystemInfo">刷新系统信息</el-button>
        </el-tab-pane>
      </el-tabs>
      
      <!-- 重启提示对话框 -->
      <el-dialog
        v-model="restartDialogVisible"
        title="需要重启服务"
        width="500px"
      >
        <p>以下配置项已修改，需要重启服务才能生效：</p>
        <ul>
          <li v-for="field in restartRequiredFields" :key="field">{{ field }}</li>
        </ul>
        <template #footer>
          <el-button @click="restartDialogVisible = false">知道了</el-button>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { Check } from '@element-plus/icons-vue'
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

const activeTab = ref('basic')
const restartDialogVisible = ref(false)
const restartRequiredFields = ref([])
const passwordUpdating = ref(false)

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const settingsForm = reactive({
  basic: {
    system_name: 'LingProxy',
    system_version: '1.0.0',
    environment: 'development',
    api_url: 'http://localhost:8080',
    port: 8080,
    host: '0.0.0.0',
    debug_mode: false
  },
  cache: {
    enabled: true,
    type: 'memory',
    ttl: 3600,
    size_limit: 100
  },
  rate_limit: {
    enabled: true,
    requests_per_minute: 1000,
    concurrency: 50
  },
  security: {
    auth_enabled: true,
    jwt_secret: '******',
    token_expiration: 24,
    https_enabled: false,
    cors: {
      enabled: true,
      allow_origins: ['*'],
      allow_methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
      allow_headers: ['*']
    }
  },
  log: {
    level: 'info',
    format: 'json',
    output: 'stdout',
    file_path: './logs/lingproxy.log',
    max_size: 100,
    max_backups: 3,
    max_age: 28,
    compress: true
  },
  load_balancer: {
    default_strategy: 'round_robin',
    health_check: {
      enabled: true,
      interval: 30,
      timeout: 5,
      max_failures: 3
    }
  },
})

const systemInfo = reactive({
  system_name: 'LingProxy',
  system_version: '1.0.0',
  environment: 'development',
  start_time: '',
  current_time: '',
  uptime: '',
  cpu_cores: 0,
  cpu_usage: 0,
  memory_total: 0,
  memory_used: 0,
  memory_usage: 0,
  disk_total: 0,
  disk_used: 0,
  disk_usage: 0,
  go_version: '',
  os: '',
  arch: ''
})

// 加载设置
const loadSettings = async () => {
  try {
    const response = await api.getSettings()
    if (response && response.data) {
      // 更新基本设置
      if (response.data.basic) {
        Object.assign(settingsForm.basic, response.data.basic)
      }
      // 更新缓存设置
      if (response.data.cache) {
        Object.assign(settingsForm.cache, response.data.cache)
      }
      // 更新限流设置
      if (response.data.rate_limit) {
        Object.assign(settingsForm.rate_limit, response.data.rate_limit)
      }
      // 更新安全设置
      if (response.data.security) {
        if (response.data.security.auth_enabled !== undefined) {
          settingsForm.security.auth_enabled = response.data.security.auth_enabled
        }
        if (response.data.security.token_expiration !== undefined) {
          settingsForm.security.token_expiration = response.data.security.token_expiration
        }
        if (response.data.security.cors) {
          settingsForm.security.cors = response.data.security.cors
        }
      }
      // 更新日志设置
      if (response.data.log) {
        Object.assign(settingsForm.log, response.data.log)
      }
      // 更新负载均衡设置
      if (response.data.load_balancer) {
        Object.assign(settingsForm.load_balancer, response.data.load_balancer)
        if (response.data.load_balancer.health_check) {
          settingsForm.load_balancer.health_check = response.data.load_balancer.health_check
        }
      }
    }
  } catch (error) {
    console.error('获取设置失败:', error)
    ElMessage.error('获取设置失败，使用默认设置')
  }
}

// 保存设置
const saveSettings = async () => {
  try {
    // 准备更新请求（只发送需要更新的字段）
    const updateData = {}
    
    if (settingsForm.basic.system_name) {
      updateData.basic = {
        system_name: settingsForm.basic.system_name,
        port: settingsForm.basic.port,
        host: settingsForm.basic.host,
        debug_mode: settingsForm.basic.debug_mode
      }
    }
    
    if (settingsForm.cache) {
      updateData.cache = {
        enabled: settingsForm.cache.enabled,
        ttl: settingsForm.cache.ttl,
        size_limit: settingsForm.cache.size_limit
      }
    }
    
    if (settingsForm.rate_limit) {
      updateData.rate_limit = {
        enabled: settingsForm.rate_limit.enabled,
        requests_per_minute: settingsForm.rate_limit.requests_per_minute,
        concurrency: settingsForm.rate_limit.concurrency
      }
    }
    
    if (settingsForm.security) {
      updateData.security = {
        auth_enabled: settingsForm.security.auth_enabled,
        token_expiration: settingsForm.security.token_expiration,
        cors: settingsForm.security.cors
      }
      // 只有当JWT密钥不是脱敏标记时才发送
      if (settingsForm.security.jwt_secret && settingsForm.security.jwt_secret !== '******') {
        updateData.security.jwt_secret = settingsForm.security.jwt_secret
      }
    }
    
    if (settingsForm.log) {
      updateData.log = settingsForm.log
    }
    
    if (settingsForm.load_balancer) {
      updateData.load_balancer = settingsForm.load_balancer
    }
    
    const response = await api.updateSettings(updateData)
    
    if (response && response.requires_restart) {
      restartRequiredFields.value = response.restart_required_fields || []
      restartDialogVisible.value = true
    }
    
    ElMessage.success('保存设置成功')
  } catch (error) {
    console.error('保存设置失败:', error)
    ElMessage.error('保存设置失败')
  }
}

// 刷新系统信息
const refreshSystemInfo = async () => {
  try {
    const response = await api.getSystemInfo()
    if (response && response.data) {
      Object.assign(systemInfo, response.data)
    }
    ElMessage.success('刷新系统信息成功')
  } catch (error) {
    console.error('刷新系统信息失败:', error)
    ElMessage.error('刷新系统信息失败')
  }
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  try {
    const date = new Date(timeStr)
    return date.toLocaleString('zh-CN')
  } catch {
    return timeStr
  }
}

// 格式化字节
const formatBytes = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

// 格式化百分比
const formatPercentage = (value) => {
  if (typeof value !== 'number' || isNaN(value)) return '0.00%'
  return value.toFixed(2) + '%'
}

// 重置密码表单
const resetPasswordForm = () => {
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}

// 更新管理员密码
const updateAdminPassword = async () => {
  // 验证表单
  if (!passwordForm.oldPassword) {
    ElMessage.warning('请输入当前密码')
    return
  }
  
  if (!passwordForm.newPassword) {
    ElMessage.warning('请输入新密码')
    return
  }
  
  if (passwordForm.newPassword.length < 6) {
    ElMessage.warning('新密码长度至少为6位')
    return
  }
  
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.warning('两次输入的新密码不一致')
    return
  }
  
  if (passwordForm.oldPassword === passwordForm.newPassword) {
    ElMessage.warning('新密码不能与当前密码相同')
    return
  }
  
  try {
    passwordUpdating.value = true
    
    await api.updateAdminPassword({
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })
    
    ElMessage.success('密码更新成功')
    resetPasswordForm()
  } catch (error) {
    console.error('更新密码失败:', error)
    const errorMsg = error.response?.data?.error || '更新密码失败'
    ElMessage.error(errorMsg)
  } finally {
    passwordUpdating.value = false
  }
}

onMounted(() => {
  loadSettings()
  refreshSystemInfo()
})
</script>

<style scoped>
.settings-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.el-tabs {
  margin-top: 20px;
}

.el-form {
  margin-top: 20px;
}

.mt-4 {
  margin-top: 20px;
}
</style>
