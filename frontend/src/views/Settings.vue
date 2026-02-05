<template>
  <div class="settings-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="page-title">{{ $t('settings.title') }}</span>
          <el-button type="primary" @click="saveSettings">
            <el-icon><Check /></el-icon>
            {{ $t('settings.saveSettings') }}
          </el-button>
        </div>
      </template>
      
      <el-tabs v-model="activeTab" class="mt-4">
        <!-- 基本设置 -->
        <el-tab-pane :label="$t('settings.basicSettings')" name="basic">
          <el-form :model="settingsForm.basic" label-width="150px">
            <el-form-item :label="$t('settings.systemName')">
              <el-input v-model="settingsForm.basic.system_name" :placeholder="$t('settings.systemNamePlaceholder')" />
            </el-form-item>
            <el-form-item :label="$t('settings.systemVersion')">
              <el-input v-model="settingsForm.basic.system_version" disabled />
            </el-form-item>
            <el-form-item :label="$t('settings.environment')">
              <el-input v-model="settingsForm.basic.environment" disabled />
            </el-form-item>
            <el-form-item :label="$t('settings.apiUrl')">
              <el-input v-model="settingsForm.basic.api_url" :placeholder="$t('settings.apiUrlPlaceholder')" />
            </el-form-item>
            <el-form-item :label="$t('settings.port')">
              <el-input-number v-model="settingsForm.basic.port" :min="1" :max="65535" />
              <el-alert type="warning" :closable="false" style="margin-top: 10px">
                {{ $t('settings.portChangeWarning') }}
              </el-alert>
            </el-form-item>
            <el-form-item :label="$t('settings.host')">
              <el-input v-model="settingsForm.basic.host" :placeholder="$t('settings.hostPlaceholder')" />
              <el-alert type="warning" :closable="false" style="margin-top: 10px">
                {{ $t('settings.hostChangeWarning') }}
              </el-alert>
            </el-form-item>
            <el-form-item :label="$t('settings.debugMode')">
              <el-switch v-model="settingsForm.basic.debug_mode" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 缓存设置 -->
        <el-tab-pane :label="$t('settings.cacheSettings')" name="cache">
          <el-form :model="settingsForm.cache" label-width="150px">
            <el-form-item :label="$t('settings.enableCache')">
              <el-switch v-model="settingsForm.cache.enabled" />
            </el-form-item>
            <el-form-item :label="$t('settings.cacheType')">
              <el-select v-model="settingsForm.cache.type" :placeholder="$t('settings.selectCacheType')" disabled>
                <el-option :label="$t('settings.memoryCache')" value="memory" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('settings.cacheTTL')">
              <el-input-number v-model="settingsForm.cache.ttl" :min="1" :max="86400" />
            </el-form-item>
            <el-form-item :label="$t('settings.cacheSizeLimit')">
              <el-input-number v-model="settingsForm.cache.size_limit" :min="1" :max="1024" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 限流设置 -->
        <el-tab-pane :label="$t('settings.rateLimitSettings')" name="rateLimit">
          <el-form :model="settingsForm.rate_limit" label-width="150px">
            <el-form-item :label="$t('settings.enableRateLimit')">
              <el-switch v-model="settingsForm.rate_limit.enabled" />
            </el-form-item>
            <el-form-item :label="$t('settings.requestsPerMinute')">
              <el-input-number v-model="settingsForm.rate_limit.requests_per_minute" :min="1" :max="10000" />
            </el-form-item>
            <el-form-item :label="$t('settings.concurrencyLimit')">
              <el-input-number v-model="settingsForm.rate_limit.concurrency" :min="1" :max="1000" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 安全设置 -->
        <el-tab-pane :label="$t('settings.securitySettings')" name="security">
          <el-form :model="settingsForm.security" label-width="150px">
            <el-form-item :label="$t('settings.enableAuth')">
              <el-switch v-model="settingsForm.security.auth_enabled" />
              <el-alert type="warning" :closable="false" style="margin-top: 10px">
                {{ $t('settings.authDisabledWarning') }}
              </el-alert>
            </el-form-item>
            <el-divider />
            <el-form-item :label="$t('settings.jwtSecret')">
              <el-input 
                v-model="settingsForm.security.jwt_secret" 
                type="password" 
                :placeholder="$t('settings.jwtSecretPlaceholder')"
                show-password
              />
              <el-alert type="warning" :closable="false" style="margin-top: 10px">
                {{ $t('settings.jwtSecretChangeWarning') }}
              </el-alert>
            </el-form-item>
            <el-form-item :label="$t('settings.tokenExpiration')">
              <el-input-number v-model="settingsForm.security.token_expiration" :min="1" :max="720" />
            </el-form-item>
            <el-form-item :label="$t('settings.enableHTTPS')">
              <el-switch v-model="settingsForm.security.https_enabled" disabled />
              <span style="color: #909399; margin-left: 10px">{{ $t('settings.notImplemented') }}</span>
            </el-form-item>
            <el-divider>{{ $t('settings.adminPassword') }}</el-divider>
            <el-form-item :label="$t('settings.currentPassword')">
              <el-input
                v-model="passwordForm.oldPassword"
                type="password"
                :placeholder="$t('settings.currentPasswordPlaceholder')"
                show-password
                style="width: 300px"
              />
            </el-form-item>
            <el-form-item :label="$t('settings.newPassword')">
              <el-input
                v-model="passwordForm.newPassword"
                type="password"
                :placeholder="$t('settings.newPasswordPlaceholder')"
                show-password
                style="width: 300px"
              />
            </el-form-item>
            <el-form-item :label="$t('settings.confirmPassword')">
              <el-input
                v-model="passwordForm.confirmPassword"
                type="password"
                :placeholder="$t('settings.confirmPasswordPlaceholder')"
                show-password
                style="width: 300px"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="updateAdminPassword" :loading="passwordUpdating">
                {{ $t('settings.updatePassword') }}
              </el-button>
              <el-button @click="resetPasswordForm">{{ $t('settings.reset') }}</el-button>
            </el-form-item>
            <el-divider>{{ $t('settings.corsSettings') }}</el-divider>
            <el-form-item :label="$t('settings.enableCORS')">
              <el-switch v-model="settingsForm.security.cors.enabled" />
            </el-form-item>
            <el-form-item :label="$t('settings.allowedOrigins')">
              <el-select 
                v-model="settingsForm.security.cors.allow_origins" 
                multiple 
                :placeholder="$t('settings.selectAllowedOrigins')"
                style="width: 100%"
              >
                <el-option :label="$t('settings.allOrigins')" value="*" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('settings.allowedMethods')">
              <el-select 
                v-model="settingsForm.security.cors.allow_methods" 
                multiple 
                :placeholder="$t('settings.selectAllowedMethods')"
                style="width: 100%"
              >
                <el-option label="GET" value="GET" />
                <el-option label="POST" value="POST" />
                <el-option label="PUT" value="PUT" />
                <el-option label="DELETE" value="DELETE" />
                <el-option label="OPTIONS" value="OPTIONS" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('settings.allowedHeaders')">
              <el-select 
                v-model="settingsForm.security.cors.allow_headers" 
                multiple 
                :placeholder="$t('settings.selectAllowedHeaders')"
                style="width: 100%"
              >
                <el-option :label="$t('settings.allOrigins')" value="*" />
              </el-select>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 日志设置 -->
        <el-tab-pane :label="$t('settings.logSettings')" name="log">
          <el-form :model="settingsForm.log" label-width="150px">
            <el-form-item :label="$t('settings.logLevel')">
              <el-select v-model="settingsForm.log.level" :placeholder="$t('settings.selectLogLevel')">
                <el-option label="DEBUG" value="debug" />
                <el-option label="INFO" value="info" />
                <el-option label="WARN" value="warn" />
                <el-option label="ERROR" value="error" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('settings.logFormat')">
              <el-select v-model="settingsForm.log.format" :placeholder="$t('settings.selectLogFormat')">
                <el-option :label="$t('settings.json')" value="json" />
                <el-option :label="$t('settings.text')" value="text" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('settings.logOutput')">
              <el-select v-model="settingsForm.log.output" :placeholder="$t('settings.selectLogOutput')">
                <el-option :label="$t('settings.stdout')" value="stdout" />
                <el-option :label="$t('settings.file')" value="file" />
                <el-option :label="$t('settings.both')" value="both" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('settings.logFilePath')" v-if="settingsForm.log.output === 'file' || settingsForm.log.output === 'both'">
              <el-input v-model="settingsForm.log.file_path" :placeholder="$t('settings.logFilePathPlaceholder')" />
            </el-form-item>
            <el-form-item :label="$t('settings.maxFileSize')">
              <el-input-number v-model="settingsForm.log.max_size" :min="1" :max="1024" />
            </el-form-item>
            <el-form-item :label="$t('settings.maxBackups')">
              <el-input-number v-model="settingsForm.log.max_backups" :min="1" :max="100" />
            </el-form-item>
            <el-form-item :label="$t('settings.maxAge')">
              <el-input-number v-model="settingsForm.log.max_age" :min="1" :max="365" />
            </el-form-item>
            <el-form-item :label="$t('settings.compressOldLogs')">
              <el-switch v-model="settingsForm.log.compress" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 负载均衡设置 -->
        <el-tab-pane :label="$t('settings.loadBalancerSettings')" name="loadBalancer">
          <el-form :model="settingsForm.load_balancer" label-width="150px">
            <el-form-item :label="$t('settings.defaultStrategy')">
              <el-select v-model="settingsForm.load_balancer.default_strategy" :placeholder="$t('settings.selectDefaultStrategy')">
                <el-option :label="$t('settings.roundRobin')" value="round_robin" />
                <el-option :label="$t('settings.leastConnections')" value="least_connections" />
                <el-option :label="$t('settings.weighted')" value="weighted" />
              </el-select>
            </el-form-item>
            <el-divider>{{ $t('settings.healthCheck') }}</el-divider>
            <el-form-item :label="$t('settings.enableHealthCheck')">
              <el-switch v-model="settingsForm.load_balancer.health_check.enabled" />
            </el-form-item>
            <el-form-item :label="$t('settings.checkInterval')">
              <el-input-number v-model="settingsForm.load_balancer.health_check.interval" :min="5" :max="300" />
            </el-form-item>
            <el-form-item :label="$t('settings.timeout')">
              <el-input-number v-model="settingsForm.load_balancer.health_check.timeout" :min="1" :max="60" />
            </el-form-item>
            <el-form-item :label="$t('settings.maxFailures')">
              <el-input-number v-model="settingsForm.load_balancer.health_check.max_failures" :min="1" :max="10" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 系统信息 -->
        <el-tab-pane :label="$t('settings.systemInfo')" name="info">
          <el-descriptions :column="2" border class="mt-4">
            <el-descriptions-item :label="$t('settings.systemInfoName')">{{ systemInfo.system_name }}</el-descriptions-item>
            <el-descriptions-item :label="$t('settings.systemInfoVersion')">{{ systemInfo.system_version }}</el-descriptions-item>
            <el-descriptions-item :label="$t('settings.systemInfoEnvironment')">{{ systemInfo.environment }}</el-descriptions-item>
            <el-descriptions-item :label="$t('settings.goVersion')">{{ systemInfo.go_version }}</el-descriptions-item>
            <el-descriptions-item :label="$t('settings.os')">{{ systemInfo.os }}</el-descriptions-item>
            <el-descriptions-item :label="$t('settings.arch')">{{ systemInfo.arch }}</el-descriptions-item>
            <el-descriptions-item :label="$t('settings.startTime')">{{ formatTime(systemInfo.start_time) }}</el-descriptions-item>
            <el-descriptions-item :label="$t('settings.currentTime')">{{ formatTime(systemInfo.current_time) }}</el-descriptions-item>
            <el-descriptions-item :label="$t('settings.uptime')">{{ systemInfo.uptime }}</el-descriptions-item>
            <el-descriptions-item :label="$t('settings.cpuCores')">{{ systemInfo.cpu_cores }}</el-descriptions-item>
            <el-descriptions-item :label="$t('settings.memoryUsage')">
              {{ formatBytes(systemInfo.memory_used) }} / {{ formatBytes(systemInfo.memory_total) }} 
              ({{ formatPercentage(systemInfo.memory_usage) }})
            </el-descriptions-item>
            <el-descriptions-item :label="$t('settings.diskUsage')" v-if="systemInfo.disk_total > 0">
              {{ formatBytes(systemInfo.disk_used) }} / {{ formatBytes(systemInfo.disk_total) }} 
              ({{ formatPercentage(systemInfo.disk_usage) }})
            </el-descriptions-item>
          </el-descriptions>
          <el-button type="primary" class="mt-4" @click="refreshSystemInfo">{{ $t('settings.refreshSystemInfo') }}</el-button>
        </el-tab-pane>
      </el-tabs>
      
      <!-- 重启提示对话框 -->
      <el-dialog
        v-model="restartDialogVisible"
        :title="$t('settings.restartRequired')"
        width="500px"
      >
        <p>{{ $t('settings.restartRequiredMessage') }}</p>
        <ul>
          <li v-for="field in restartRequiredFields" :key="field">{{ field }}</li>
        </ul>
        <template #footer>
          <el-button @click="restartDialogVisible = false">{{ $t('settings.gotIt') }}</el-button>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { Check } from '@element-plus/icons-vue'
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import api from '../api'

const { t } = useI18n()

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
    ElMessage.error(t('settings.getSettingsFailed'))
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
    
    ElMessage.success(t('settings.saveSettingsSuccess'))
  } catch (error) {
    console.error('保存设置失败:', error)
    ElMessage.error(t('settings.saveSettingsFailed'))
  }
}

// 刷新系统信息
const refreshSystemInfo = async () => {
  try {
    const response = await api.getSystemInfo()
    if (response && response.data) {
      Object.assign(systemInfo, response.data)
    }
    ElMessage.success(t('settings.refreshSystemInfoSuccess'))
  } catch (error) {
    console.error('刷新系统信息失败:', error)
    ElMessage.error(t('settings.refreshSystemInfoFailed'))
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
    ElMessage.warning(t('settings.currentPasswordRequired'))
    return
  }
  
  if (!passwordForm.newPassword) {
    ElMessage.warning(t('settings.newPasswordRequired'))
    return
  }
  
  if (passwordForm.newPassword.length < 6) {
    ElMessage.warning(t('settings.newPasswordMinLength'))
    return
  }
  
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.warning(t('settings.passwordMismatch'))
    return
  }
  
  if (passwordForm.oldPassword === passwordForm.newPassword) {
    ElMessage.warning(t('settings.passwordSameAsOld'))
    return
  }
  
  try {
    passwordUpdating.value = true
    
    await api.updateAdminPassword({
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })
    
    ElMessage.success(t('settings.passwordUpdateSuccess'))
    resetPasswordForm()
  } catch (error) {
    console.error('更新密码失败:', error)
    const errorMsg = error.response?.data?.error || t('settings.passwordUpdateFailed')
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
