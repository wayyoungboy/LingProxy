<template>
  <div class="settings-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>系统设置</span>
          <el-button type="primary" @click="saveSettings">保存设置</el-button>
        </div>
      </template>
      
      <el-tabs v-model="activeTab" class="mt-4">
        <el-tab-pane label="基本设置" name="basic">
          <el-form :model="settingsForm" label-width="120px">
            <el-form-item label="系统名称">
              <el-input v-model="settingsForm.system_name" placeholder="请输入系统名称" />
            </el-form-item>
            <el-form-item label="系统版本">
              <el-input v-model="settingsForm.system_version" placeholder="请输入系统版本" disabled />
            </el-form-item>
            <el-form-item label="API地址">
              <el-input v-model="settingsForm.api_url" placeholder="请输入API地址" />
            </el-form-item>
            <el-form-item label="是否启用调试模式">
              <el-switch v-model="settingsForm.debug_mode" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="缓存设置" name="cache">
          <el-form :model="settingsForm.cache" label-width="120px">
            <el-form-item label="缓存类型">
              <el-select v-model="settingsForm.cache.type" placeholder="请选择缓存类型">
                <el-option label="内存缓存" value="memory" />
                <el-option label="Redis" value="redis" />
              </el-select>
            </el-form-item>
            <el-form-item label="缓存过期时间(秒)">
              <el-input-number v-model="settingsForm.cache.expiration" :min="1" :max="86400" />
            </el-form-item>
            <el-form-item label="缓存大小限制(MB)">
              <el-input-number v-model="settingsForm.cache.size_limit" :min="1" :max="1024" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="限流设置" name="rateLimit">
          <el-form :model="settingsForm.rate_limit" label-width="120px">
            <el-form-item label="是否启用限流">
              <el-switch v-model="settingsForm.rate_limit.enabled" />
            </el-form-item>
            <el-form-item label="请求频率限制(QPS)">
              <el-input-number v-model="settingsForm.rate_limit.qps" :min="1" :max="1000" />
            </el-form-item>
            <el-form-item label="并发请求限制">
              <el-input-number v-model="settingsForm.rate_limit.concurrency" :min="1" :max="1000" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="安全设置" name="security">
          <el-form :model="settingsForm.security" label-width="120px">
            <el-form-item label="JWT密钥">
              <el-input v-model="settingsForm.security.jwt_secret" type="password" placeholder="请输入JWT密钥" />
            </el-form-item>
            <el-form-item label="Token过期时间(小时)">
              <el-input-number v-model="settingsForm.security.token_expiration" :min="1" :max="720" />
            </el-form-item>
            <el-form-item label="是否启用HTTPS">
              <el-switch v-model="settingsForm.security.https_enabled" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="系统信息" name="info">
          <el-descriptions :column="1" border class="mt-4">
            <el-descriptions-item label="系统名称">{{ systemInfo.system_name }}</el-descriptions-item>
            <el-descriptions-item label="系统版本">{{ systemInfo.system_version }}</el-descriptions-item>
            <el-descriptions-item label="运行环境">{{ systemInfo.environment }}</el-descriptions-item>
            <el-descriptions-item label="启动时间">{{ systemInfo.start_time }}</el-descriptions-item>
            <el-descriptions-item label="当前时间">{{ systemInfo.current_time }}</el-descriptions-item>
            <el-descriptions-item label="CPU核心数">{{ systemInfo.cpu_cores }}</el-descriptions-item>
            <el-descriptions-item label="内存使用">{{ systemInfo.memory_usage }}</el-descriptions-item>
            <el-descriptions-item label="磁盘使用">{{ systemInfo.disk_usage }}</el-descriptions-item>
          </el-descriptions>
          <el-button type="primary" class="mt-4" @click="refreshSystemInfo">刷新系统信息</el-button>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api'

const activeTab = ref('basic')
const settingsForm = reactive({
  system_name: 'LingProxy',
  system_version: '1.0.0',
  api_url: 'http://localhost:8080',
  debug_mode: false,
  cache: {
    type: 'memory',
    expiration: 3600,
    size_limit: 100
  },
  rate_limit: {
    enabled: true,
    qps: 100,
    concurrency: 50
  },
  security: {
    jwt_secret: 'your-secret-key',
    token_expiration: 24,
    https_enabled: false
  }
})

const systemInfo = reactive({
  system_name: 'LingProxy',
  system_version: '1.0.0',
  environment: 'development',
  start_time: '2024-01-01 00:00:00',
  current_time: new Date().toLocaleString(),
  cpu_cores: 4,
  memory_usage: '50%',
  disk_usage: '30%'
})

const loadSettings = async () => {
  try {
    const response = await api.getSettings()
    if (response.data) {
      Object.assign(settingsForm, response.data)
    }
  } catch (error) {
    ElMessage.error('获取设置失败，使用默认设置')
  }
}

const saveSettings = async () => {
  try {
    await api.updateSettings(settingsForm)
    ElMessage.success('保存设置成功')
  } catch (error) {
    ElMessage.error('保存设置失败')
  }
}

const refreshSystemInfo = async () => {
  try {
    const response = await api.getSystemInfo()
    if (response.data) {
      Object.assign(systemInfo, response.data)
      systemInfo.current_time = new Date().toLocaleString()
    }
    ElMessage.success('刷新系统信息成功')
  } catch (error) {
    ElMessage.error('刷新系统信息失败')
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

.el-tabs {
  margin-top: 20px;
}

.el-form {
  margin-top: 20px;
}
</style>