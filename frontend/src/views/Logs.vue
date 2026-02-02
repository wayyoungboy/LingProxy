<template>
  <div class="logs-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>日志管理</span>
          <div>
            <el-button type="success" @click="refreshLogs" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button type="danger" @click="handleClearLogs" :disabled="!selectedFile">
              <el-icon><Delete /></el-icon>
              清空日志
            </el-button>
            <el-button type="primary" @click="handleDownloadLog" :disabled="!selectedFile">
              <el-icon><Download /></el-icon>
              下载日志
            </el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20">
        <!-- 左侧：日志文件列表 -->
        <el-col :span="6">
          <el-card shadow="hover">
            <template #header>
              <span>日志文件</span>
            </template>
            <el-scrollbar height="600px">
              <el-menu
                :default-active="selectedFile"
                @select="handleFileSelect"
                class="log-files-menu"
              >
                <el-menu-item
                  v-for="file in logFiles"
                  :key="file.name"
                  :index="file.name"
                >
                  <div class="file-item">
                    <div class="file-name">{{ file.name }}</div>
                    <div class="file-info">
                      <span class="file-size">{{ formatFileSize(file.size) }}</span>
                      <span class="file-time">{{ formatTime(file.mod_time) }}</span>
                    </div>
                  </div>
                </el-menu-item>
              </el-menu>
              <el-empty v-if="logFiles.length === 0" description="暂无日志文件" />
            </el-scrollbar>
          </el-card>
        </el-col>

        <!-- 右侧：日志内容 -->
        <el-col :span="18">
          <el-card shadow="hover">
            <template #header>
              <div class="log-header">
                <span>日志内容</span>
                <div class="log-controls">
                  <el-select
                    v-model="logLevel"
                    placeholder="日志级别"
                    style="width: 120px; margin-right: 10px"
                    @change="handleFilterChange"
                  >
                    <el-option label="全部" value="" />
                    <el-option label="DEBUG" value="DEBUG" />
                    <el-option label="INFO" value="INFO" />
                    <el-option label="WARN" value="WARN" />
                    <el-option label="ERROR" value="ERROR" />
                    <el-option label="FATAL" value="FATAL" />
                  </el-select>
                  <el-input
                    v-model="keyword"
                    placeholder="搜索关键词"
                    style="width: 200px; margin-right: 10px"
                    clearable
                    @clear="handleFilterChange"
                    @keyup.enter="handleFilterChange"
                  >
                    <template #prefix>
                      <el-icon><Search /></el-icon>
                    </template>
                  </el-input>
                  <el-input-number
                    v-model="lines"
                    :min="10"
                    :max="1000"
                    :step="10"
                    style="width: 120px; margin-right: 10px"
                    @change="handleFilterChange"
                  />
                  <el-checkbox v-model="tail" @change="handleFilterChange">
                    从尾部读取
                  </el-checkbox>
                </div>
              </div>
            </template>

            <div class="log-content" ref="logContentRef">
              <div
                v-for="entry in logEntries"
                :key="entry.line"
                :class="['log-entry', `log-level-${entry.level?.toLowerCase()}`]"
              >
                <span class="log-line">{{ entry.line }}</span>
                <span class="log-timestamp" v-if="entry.timestamp">
                  {{ formatTimestamp(entry.timestamp) }}
                </span>
                <span class="log-level" v-if="entry.level">
                  [{{ entry.level }}]
                </span>
                <span class="log-message">{{ entry.message || entry.raw }}</span>
              </div>
              <el-empty v-if="logEntries.length === 0" description="暂无日志内容" />
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Delete, Download, Search } from '@element-plus/icons-vue'
import api from '../api'

const loading = ref(false)
const logFiles = ref([])
const logEntries = ref([])
const selectedFile = ref('')
const logLevel = ref('')
const keyword = ref('')
const lines = ref(100)
const tail = ref(true)
const logContentRef = ref(null)

// 获取日志文件列表
const getLogFiles = async () => {
  try {
    loading.value = true
    const response = await api.getLogFiles()
    if (response && response.data) {
      logFiles.value = response.data
      // 如果没有选中文件，选择第一个
      if (!selectedFile.value && logFiles.value.length > 0) {
        selectedFile.value = logFiles.value[0].name
        getLogs()
      }
    }
  } catch (error) {
    console.error('获取日志文件列表失败:', error)
    ElMessage.error('获取日志文件列表失败')
  } finally {
    loading.value = false
  }
}

// 获取日志内容
const getLogs = async () => {
  if (!selectedFile.value) {
    return
  }

  try {
    loading.value = true
    const params = {
      file: selectedFile.value,
      lines: lines.value,
      tail: tail.value
    }
    if (logLevel.value) {
      params.level = logLevel.value
    }
    if (keyword.value) {
      params.keyword = keyword.value
    }

    const response = await api.getLogs(params)
    if (response && response.data) {
      logEntries.value = response.data
      // 滚动到底部
      await nextTick()
      if (logContentRef.value) {
        logContentRef.value.scrollTop = logContentRef.value.scrollHeight
      }
    }
  } catch (error) {
    console.error('获取日志内容失败:', error)
    ElMessage.error('获取日志内容失败')
  } finally {
    loading.value = false
  }
}

// 选择文件
const handleFileSelect = (fileName) => {
  selectedFile.value = fileName
  getLogs()
}

// 过滤条件改变
const handleFilterChange = () => {
  getLogs()
}

// 刷新日志
const refreshLogs = () => {
  getLogFiles()
  if (selectedFile.value) {
    getLogs()
  }
}

// 下载日志文件
const handleDownloadLog = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择要下载的日志文件')
    return
  }

  try {
    const response = await api.downloadLogFile(selectedFile.value)
    // 创建blob对象
    const blob = new Blob([response], {
      type: 'application/octet-stream'
    })
    // 创建下载链接
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = selectedFile.value
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('日志文件下载成功')
  } catch (error) {
    console.error('下载日志文件失败:', error)
    ElMessage.error('下载日志文件失败')
  }
}

// 清空日志
const handleClearLogs = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择要清空的日志文件')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要清空日志文件 "${selectedFile.value}" 吗？此操作不可恢复。`,
      '清空确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await api.clearLogs(selectedFile.value)
    ElMessage.success('日志文件已清空')
    // 刷新日志内容
    getLogs()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清空日志文件失败:', error)
      ElMessage.error('清空日志文件失败')
    }
  }
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
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

// 格式化时间戳
const formatTimestamp = (timestamp) => {
  if (!timestamp) return ''
  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      fractionalSecondDigits: 3
    })
  } catch {
    return timestamp
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getLogFiles()
})
</script>

<style scoped>
.logs-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.log-files-menu {
  border: none;
}

.file-item {
  width: 100%;
}

.file-name {
  font-weight: 500;
  margin-bottom: 4px;
  word-break: break-all;
}

.file-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
}

.file-size {
  margin-right: 10px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.log-controls {
  display: flex;
  align-items: center;
}

.log-content {
  height: 600px;
  overflow-y: auto;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  background-color: #1e1e1e;
  color: #d4d4d4;
  padding: 10px;
  border-radius: 4px;
}

.log-entry {
  display: flex;
  margin-bottom: 2px;
  padding: 2px 0;
  line-height: 1.5;
}

.log-entry:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.log-line {
  color: #858585;
  margin-right: 10px;
  min-width: 50px;
  text-align: right;
}

.log-timestamp {
  color: #608b4e;
  margin-right: 10px;
  min-width: 180px;
}

.log-level {
  margin-right: 10px;
  min-width: 60px;
  font-weight: bold;
}

.log-level-debug .log-level {
  color: #569cd6;
}

.log-level-info .log-level {
  color: #4ec9b0;
}

.log-level-warn .log-level {
  color: #dcdcaa;
}

.log-level-error .log-level {
  color: #f48771;
}

.log-level-fatal .log-level {
  color: #f48771;
}

.log-message {
  flex: 1;
  word-break: break-all;
}

.log-level-debug {
  color: #569cd6;
}

.log-level-info {
  color: #4ec9b0;
}

.log-level-warn {
  color: #dcdcaa;
}

.log-level-error {
  color: #f48771;
}

.log-level-fatal {
  color: #f48771;
}
</style>
