<template>
  <div class="usage-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="page-title">{{ $t('llmResourceUsage.title') }}</span>
          <el-button type="primary" @click="refreshData">
            <el-icon><Refresh /></el-icon>
            {{ $t('dashboard.refreshData') }}
          </el-button>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="20" class="mb-4">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-label">{{ $t('llmResourceUsage.totalTokens') }}</div>
              <div class="stat-value">{{ formatNumber(totalTokens) }}</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-label">{{ $t('llmResourceUsage.totalRequests') }}</div>
              <div class="stat-value">{{ formatNumber(totalRequests) }}</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-label">{{ $t('llmResourceUsage.successRequests') }}</div>
              <div class="stat-value">{{ formatNumber(successRequests) }}</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-label">{{ $t('llmResourceUsage.averageTokens') }}</div>
              <div class="stat-value">{{ formatNumber(avgTokensPerRequest) }}</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 筛选表单 -->
      <el-form :inline="true" :model="searchForm" class="mb-4">
        <el-form-item :label="$t('llmResourceUsage.timeRange')">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            :range-separator="$t('llmResourceUsage.to')"
            :start-placeholder="$t('llmResourceUsage.startDate')"
            :end-placeholder="$t('llmResourceUsage.endDate')"
            style="width: 250px"
            @change="handleDateRangeChange"
          />
        </el-form-item>
        <el-form-item :label="$t('llmResourceUsage.resourceName')">
          <el-input
            v-model="searchForm.resourceName"
            :placeholder="$t('llmResourceUsage.resourceNamePlaceholder')"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadUsageData">{{ $t('common.search') }}</el-button>
          <el-button @click="resetForm">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- Token使用统计表格 -->
      <el-table
        :data="usageList"
        style="width: 100%"
        v-loading="loading"
        :default-sort="{ prop: 'total_tokens', order: 'descending' }"
        stripe
      >
        <el-table-column
          prop="resource_name"
          :label="$t('llmResourceUsage.resourceName')"
          width="200"
          show-overflow-tooltip
        />
        <el-table-column
          prop="resource_type"
          :label="$t('llmResourceUsage.resourceType')"
          width="120"
        >
          <template #default="scope">
            <el-tag type="info">{{ getTypeLabel(scope.row.resource_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="model"
          :label="$t('llmResourceUsage.modelId')"
          width="180"
          show-overflow-tooltip
        />
        <el-table-column
          prop="total_tokens"
          :label="$t('llmResourceUsage.tokenUsage')"
          width="150"
          sortable
        >
          <template #default="scope">
            <span class="token-value">{{ formatNumber(scope.row.total_tokens || 0) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="total_requests"
          :label="$t('llmResourceUsage.requestCount')"
          width="120"
          sortable
        >
          <template #default="scope">
            {{ formatNumber(scope.row.total_requests || 0) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="success_requests"
          :label="$t('llmResourceUsage.successCount')"
          width="120"
        >
          <template #default="scope">
            <el-tag type="success" size="small">
              {{ formatNumber(scope.row.success_requests || 0) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="failed_requests"
          :label="$t('llmResourceUsage.failedCount')"
          width="120"
        >
          <template #default="scope">
            <el-tag type="danger" size="small" v-if="scope.row.failed_requests > 0">
              {{ formatNumber(scope.row.failed_requests || 0) }}
            </el-tag>
            <span v-else>0</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="success_rate"
          :label="$t('llmResourceUsage.successRate')"
          width="120"
          sortable
        >
          <template #default="scope">
            <el-tag :type="getSuccessRateType(scope.row.success_rate)" size="small">
              {{ (scope.row.success_rate || 0).toFixed(2) }}%
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="avg_tokens_per_request"
          :label="$t('llmResourceUsage.avgTokensPerRequest')"
          width="150"
          sortable
        >
          <template #default="scope">
            {{ formatNumber(Math.round(scope.row.avg_tokens_per_request || 0)) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="last_request_time"
          :label="$t('llmResourceUsage.lastRequestTime')"
          width="180"
          sortable
        >
          <template #default="scope">
            {{ scope.row.last_request_time ? formatDate(scope.row.last_request_time) : '-' }}
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import api from '../api'
import { MODEL_TYPE_LABELS } from '../utils/constants'

const { t } = useI18n()

const loading = ref(false)
const usageList = ref([])
const totalTokens = ref(0)
const totalRequests = ref(0)
const successRequests = ref(0)

const searchForm = ref({
  dateRange: null,
  resourceName: ''
})

const pagination = ref({
  current: 1,
  size: 20,
  total: 0
})

// 计算平均Token/请求
const avgTokensPerRequest = computed(() => {
  if (totalRequests.value === 0) return 0
  return Math.round(totalTokens.value / totalRequests.value)
})

// 加载使用量数据
const loadUsageData = async () => {
  try {
    loading.value = true

    // 调用后端API获取按资源分组的统计
    // 注意：apiClient的响应拦截器已经返回了response.data，所以这里直接使用response.data
    const response = await api.getLLMResourceUsageStats()
    const usageArray = response?.data || []

    // 应用筛选
    let filtered = usageArray
    if (searchForm.value.resourceName) {
      filtered = filtered.filter(item =>
        item.resource_name.toLowerCase().includes(searchForm.value.resourceName.toLowerCase())
      )
    }

    // 应用时间范围筛选
    if (searchForm.value.dateRange && searchForm.value.dateRange.length === 2) {
      const startDate = new Date(searchForm.value.dateRange[0])
      const endDate = new Date(searchForm.value.dateRange[1])
      endDate.setHours(23, 59, 59, 999) // 设置为当天的最后一刻

      filtered = filtered.filter(item => {
        if (!item.last_request_time) return false
        const requestTime = new Date(item.last_request_time)
        return requestTime >= startDate && requestTime <= endDate
      })
    }

    // 计算总计
    totalTokens.value = filtered.reduce((sum, item) => sum + (item.total_tokens || 0), 0)
    totalRequests.value = filtered.reduce((sum, item) => sum + (item.total_requests || 0), 0)
    successRequests.value = filtered.reduce((sum, item) => sum + (item.success_requests || 0), 0)

    // 分页
    const startIndex = (pagination.value.current - 1) * pagination.value.size
    const endIndex = startIndex + pagination.value.size
    usageList.value = filtered.slice(startIndex, endIndex)
    pagination.value.total = filtered.length
  } catch (error) {
    console.error('加载使用量数据失败:', error)
    ElMessage.error(
      t('llmResourceUsage.loadFailed') + ': ' + (error.response?.data?.error || error.message)
    )
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  loadUsageData()
}

// 处理日期范围变化
const handleDateRangeChange = () => {
  loadUsageData()
}

// 重置表单
const resetForm = () => {
  searchForm.value = {
    dateRange: null,
    resourceName: ''
  }
  pagination.value.current = 1
  loadUsageData()
}

// 分页处理
const handleSizeChange = size => {
  pagination.value.size = size
  pagination.value.current = 1
  loadUsageData()
}

const handleCurrentChange = current => {
  pagination.value.current = current
  loadUsageData()
}

// 格式化数字
const formatNumber = num => {
  if (num === null || num === undefined) return '0'
  return num.toLocaleString('zh-CN')
}

// 格式化日期
const formatDate = dateString => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 获取类型标签
const getTypeLabel = type => {
  const labels = {
    chat: t('llmResources.typeChat'),
    image: t('llmResources.typeImage'),
    embedding: t('llmResources.typeEmbedding'),
    rerank: t('llmResources.typeRerank'),
    audio: t('llmResources.typeAudio'),
    video: t('llmResources.typeVideo')
  }
  return labels[type] || MODEL_TYPE_LABELS[type] || type || '-'
}

// 获取成功率标签类型
const getSuccessRateType = rate => {
  if (rate >= 95) return 'success'
  if (rate >= 80) return 'warning'
  return 'danger'
}

// 组件挂载时加载数据
onMounted(() => {
  loadUsageData()
})
</script>

<style scoped>
.usage-container {
  padding: 0;
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

.stat-card {
  text-align: center;
}

.stat-content {
  padding: 10px 0;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #2563eb;
}

.token-value {
  font-weight: 600;
  color: #2563eb;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
}

.mb-4 {
  margin-bottom: 16px;
}

.mt-4 {
  margin-top: 16px;
}
</style>
