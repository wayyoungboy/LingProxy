<template>
  <div class="requests-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="page-title">{{ $t('requests.title') }}</span>
          <el-button type="primary" @click="exportRequests">
            <el-icon><Download /></el-icon>
            {{ $t('requests.export') }}
          </el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="searchForm" class="search-form mb-4">
        <el-form-item :label="$t('requests.path')">
          <el-input v-model="searchForm.path" :placeholder="$t('requests.pathPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('requests.selectStatus')">
            <el-option :label="$t('requests.success')" value="success" />
            <el-option :label="$t('requests.failed')" value="error" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('llmResourceUsage.timeRange')">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            :range-separator="$t('llmResourceUsage.to')"
            :start-placeholder="$t('llmResourceUsage.startDate')"
            :end-placeholder="$t('llmResourceUsage.endDate')"
            style="width: 250px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getRequests">{{ $t('common.search') }}</el-button>
          <el-button @click="resetForm">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
      
      <el-table 
        :data="requestsList" 
        style="width: 100%" 
        v-loading="loading"
        border
        stripe
      >
        <el-table-column prop="id" label="ID" width="180" />
        <el-table-column prop="endpoint" :label="$t('requests.path')" />
        <el-table-column prop="method" :label="$t('requests.method')" width="100" />
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'success' ? 'success' : 'danger'">
              {{ scope.row.status === 'success' ? $t('requests.success') : $t('requests.failed') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" :label="$t('requests.duration')" width="120" />
        <el-table-column prop="tokens" :label="$t('requests.tokens')" width="120" />
        <el-table-column prop="created_at" :label="$t('requests.requestTime')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewRequestDetail(scope.row)">
              {{ $t('requests.viewDetails') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container mt-4">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 请求详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="$t('requests.detailTitle')"
      width="800px"
    >
      <div v-if="currentRequest" class="request-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item :label="$t('requests.requestId')">{{ currentRequest.id }}</el-descriptions-item>
          <el-descriptions-item :label="$t('dashboard.userId')">{{ currentRequest.user_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('requests.path')">{{ currentRequest.endpoint }}</el-descriptions-item>
          <el-descriptions-item :label="$t('requests.method')">{{ currentRequest.method }}</el-descriptions-item>
          <el-descriptions-item :label="$t('common.status')">
            <el-tag :type="currentRequest.status === 'success' ? 'success' : 'danger'">
              {{ currentRequest.status === 'success' ? $t('requests.success') : $t('requests.failed') }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('requests.duration')">{{ currentRequest.duration }}ms</el-descriptions-item>
          <el-descriptions-item :label="$t('requests.tokens')">{{ currentRequest.tokens || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('requests.requestTime')">{{ formatDate(currentRequest.created_at) }}</el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="detailDialogVisible = false">{{ $t('common.close') }}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download } from '@element-plus/icons-vue'
import api from '../api'

const { t } = useI18n()

const requestsList = ref([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const currentRequest = ref(null)

const searchForm = reactive({
  path: '',
  status: '',
  dateRange: []
})

const pagination = reactive({
  current: 1,
  size: 10,
  total: 0
})

const getRequests = async () => {
  loading.value = true
  try {
    const params = {
      limit: pagination.size * pagination.current  // 后端使用limit参数
    }
    
    const response = await api.getRequests(params)
    // 后端返回格式: { "data": [...] }
    if (response && response.data) {
      // 如果是数组，直接使用
      if (Array.isArray(response.data)) {
        // 客户端分页和过滤
        let filtered = response.data
        
        // 按路径过滤
        if (searchForm.path) {
          filtered = filtered.filter(r => r.endpoint && r.endpoint.includes(searchForm.path))
        }
        
        // 按状态过滤
        if (searchForm.status) {
          filtered = filtered.filter(r => r.status === searchForm.status)
        }
        
        // 按时间范围过滤
        if (searchForm.dateRange && searchForm.dateRange.length === 2) {
          const start = new Date(searchForm.dateRange[0])
          const end = new Date(searchForm.dateRange[1])
          end.setHours(23, 59, 59, 999) // 包含结束日期
          filtered = filtered.filter(r => {
            const createdAt = new Date(r.created_at)
            return createdAt >= start && createdAt <= end
          })
        }
        
        // 客户端分页
        const start = (pagination.current - 1) * pagination.size
        const end = start + pagination.size
        requestsList.value = filtered.slice(start, end)
        pagination.total = filtered.length
      } else {
        // 如果后端返回分页格式
        requestsList.value = response.data.items || []
        pagination.total = response.data.total || 0
      }
    } else {
      requestsList.value = []
      pagination.total = 0
    }
  } catch (error) {
    console.error('获取请求列表失败:', error)
    ElMessage.error(t('requests.getListFailed'))
    requestsList.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

const viewRequestDetail = async (request) => {
  try {
    const response = await api.getRequestDetail(request.id)
    currentRequest.value = response.data
    detailDialogVisible.value = true
  } catch (error) {
    ElMessage.error(t('requests.getDetailFailed'))
  }
}

const exportRequests = async () => {
  try {
    await api.exportRequests(searchForm)
    ElMessage.success(t('requests.exportSuccess'))
  } catch (error) {
    ElMessage.error(t('requests.exportFailed'))
  }
}

const resetForm = () => {
  Object.assign(searchForm, {
    path: '',
    status: '',
    dateRange: []
  })
  getRequests()
}

const handleSizeChange = (size) => {
  pagination.size = size
  getRequests()
}

const handleCurrentChange = (current) => {
  pagination.current = current
  getRequests()
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

const formatJson = (jsonString) => {
  try {
    if (typeof jsonString === 'string') {
      return JSON.stringify(JSON.parse(jsonString), null, 2)
    }
    return JSON.stringify(jsonString, null, 2)
  } catch (error) {
    return jsonString
  }
}

onMounted(() => {
  getRequests()
})
</script>

<style scoped>
.requests-container {
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

.search-form {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}

.request-detail {
  max-height: 600px;
  overflow-y: auto;
}

pre {
  margin: 0;
  font-family: monospace;
  font-size: 14px;
  line-height: 1.5;
}
</style>