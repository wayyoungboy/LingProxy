<template>
  <div class="requests-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('requests.title') }}</h1>
    </div>

    <el-card class="filter-card">

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

      <el-table :data="requestsList" style="width: 100%" v-loading="loading" border stripe>
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
    <el-dialog v-model="detailDialogVisible" :title="$t('requests.detailTitle')" width="800px">
      <div v-if="currentRequest" class="request-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item :label="$t('requests.requestId')">
            {{ currentRequest.id }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('dashboard.userId')">
            {{ currentRequest.user_id || '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('requests.path')">
            {{ currentRequest.endpoint }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('requests.method')">
            {{ currentRequest.method }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('common.status')">
            <el-tag :type="currentRequest.status === 'success' ? 'success' : 'danger'">
              {{
                currentRequest.status === 'success' ? $t('requests.success') : $t('requests.failed')
              }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('requests.duration')">
            {{ currentRequest.duration }}ms
          </el-descriptions-item>
          <el-descriptions-item :label="$t('requests.tokens')">
            {{ currentRequest.tokens || 0 }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('requests.requestTime')">
            {{ formatDate(currentRequest.created_at) }}
          </el-descriptions-item>
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
import { ElMessage } from 'element-plus'
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
    // 获取足够的数据用于客户端分页（最多1000条）
    const params = {
      limit: 1000
    }

    // 添加搜索参数
    if (searchForm.path) {
      params.endpoint = searchForm.path
    }

    if (searchForm.status) {
      params.status = searchForm.status
    }

    // 时间范围参数
    if (searchForm.dateRange && searchForm.dateRange.length === 2) {
      const start = new Date(searchForm.dateRange[0])
      const end = new Date(searchForm.dateRange[1])
      end.setHours(23, 59, 59, 999) // 包含结束日期
      params.start_time = start.toISOString()
      params.end_time = end.toISOString()
    }

    const response = await api.getRequests(params)
    if (response && response.data) {
      const allRequests = response.data
      // 客户端分页
      const start = (pagination.current - 1) * pagination.size
      const end = start + pagination.size
      requestsList.value = allRequests.slice(start, end)
      pagination.total = allRequests.length
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

const viewRequestDetail = async request => {
  try {
    const response = await api.getRequestDetail(request.id)
    currentRequest.value = response.data
    detailDialogVisible.value = true
  } catch (error) {
    ElMessage.error(t('requests.getDetailFailed'))
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

const handleSizeChange = size => {
  pagination.size = size
  getRequests()
}

const handleCurrentChange = current => {
  pagination.current = current
  getRequests()
}

const formatDate = dateString => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  getRequests()
})
</script>

<style scoped>
.requests-page {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-family: var(--font-serif);
  font-size: 32px;
  font-weight: 500;
  line-height: 1.1;
  color: var(--claude-text-primary);
}

.filter-card {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 0;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.request-detail {
  max-height: 600px;
  overflow-y: auto;
}

pre {
  margin: 0;
  font-family: var(--font-mono);
  font-size: 14px;
  line-height: 1.5;
}

.text-caption {
  font-size: 14px;
  color: var(--claude-text-secondary);
}
</style>
