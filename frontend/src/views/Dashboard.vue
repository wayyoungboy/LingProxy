<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1 class="page-title">{{ $t('dashboard.systemDashboard') }}</h1>
      <el-button type="primary" @click="handleRefresh" :loading="refreshing">
        <el-icon><Refresh /></el-icon>
        {{ $t('dashboard.refreshData') }}
      </el-button>
    </div>

    <!-- 统计卡片 - Claude style -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-card-content">
          <div class="stat-icon stat-icon-terracotta">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ formatNumber(stats.total_users) }}</div>
            <div class="stat-label">{{ $t('dashboard.totalUsers') }}</div>
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-card-content">
          <div class="stat-icon stat-icon-green">
            <el-icon><Cpu /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ formatNumber(stats.total_llm_resources) }}</div>
            <div class="stat-label">{{ $t('dashboard.totalLLMResources') }}</div>
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-card-content">
          <div class="stat-icon stat-icon-coral">
            <el-icon><Message /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ formatNumber(stats.total_requests) }}</div>
            <div class="stat-label">{{ $t('dashboard.totalRequests') }}</div>
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-card-content">
          <div class="stat-icon stat-icon-success">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ (stats.success_rate || 0).toFixed(2) }}%</div>
            <div class="stat-label">{{ $t('dashboard.successRate') }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 性能指标 - Claude style -->
    <el-card class="performance-card">
      <template #header>
        <span class="section-title">{{ $t('dashboard.systemPerformance') }}</span>
      </template>
      <div class="performance-metrics">
        <div class="metric-card">
          <div class="metric-label">{{ $t('dashboard.avgResponseTime') }}</div>
          <div class="metric-value">{{ (stats.avg_response_time || 0).toFixed(2) }}ms</div>
        </div>
        <div class="metric-card">
          <div class="metric-label">{{ $t('dashboard.todayRequests') }}</div>
          <div class="metric-value metric-value-serif">{{ formatNumber(todayRequests) }}</div>
        </div>
        <div class="metric-card">
          <div class="metric-label">{{ $t('dashboard.systemStatus') }}</div>
          <div class="metric-value metric-status">
            <span class="status-dot"></span>
            {{ $t('dashboard.running') }}
          </div>
        </div>
        <div class="metric-card">
          <div class="metric-label">{{ $t('dashboard.runningTime') }}</div>
          <div class="metric-value">{{ uptime || '--' }}</div>
        </div>
      </div>
    </el-card>

    <!-- 最近请求 - Claude style -->
    <el-card class="requests-card">
      <template #header>
        <div class="card-header-row">
          <span class="section-title">{{ $t('dashboard.recentRequests') }}</span>
          <el-button type="primary" size="small" @click="navigateToRequests">
            {{ $t('dashboard.viewAll') }}
          </el-button>
        </div>
      </template>

      <el-table v-if="recentRequests.length > 0" :data="recentRequests" style="width: 100%">
        <el-table-column prop="id" :label="$t('requests.requestId')" width="180" />
        <el-table-column prop="user_id" :label="$t('dashboard.userId')" width="180" />
        <el-table-column prop="endpoint" :label="$t('llmResources.endpoint')" />
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ scope.row.status === 'success' ? $t('requests.success') : $t('requests.failed') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" :label="$t('requests.duration')" width="120">
          <template #default="scope">{{ scope.row.duration }}ms</template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180">
          <template #default="scope">{{ formatDate(scope.row.created_at) }}</template>
        </el-table-column>
      </el-table>

      <el-empty v-else :description="$t('dashboard.noRequests')" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh, User, Cpu, Message, Check } from '@element-plus/icons-vue'
import api from '../api'
import { formatDate, formatNumber } from '../utils/index'

const router = useRouter()
const { t } = useI18n()
const refreshing = ref(false)
const stats = reactive({
  total_users: 0,
  total_llm_resources: 0,
  total_requests: 0,
  success_rate: 0,
  avg_response_time: 0
})
const recentRequests = ref([])
const systemInfo = ref(null)

// 计算属性
const todayRequests = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return recentRequests.value.filter(req => {
    const reqDate = new Date(req.created_at)
    return reqDate >= today
  }).length
})

const uptime = computed(() => {
  return systemInfo.value?.uptime || '--'
})

// 获取系统统计信息
const getSystemStats = async () => {
  try {
    refreshing.value = true

    const [statsResponse, systemInfoResponse, requestsResponse] = await Promise.allSettled([
      api.getSystemStats(),
      api.getSystemInfo(),
      api.getRequests({ limit: 10 })
    ])

    if (statsResponse.status === 'fulfilled' && statsResponse.value) {
      const data = statsResponse.value.data || statsResponse.value
      Object.assign(stats, {
        total_users: data.total_users || 0,
        total_llm_resources: data.total_llm_resources || 0,
        total_requests: data.total_requests || 0,
        success_rate: data.success_rate || 0,
        avg_response_time: data.avg_response_time || 0
      })
    }

    if (systemInfoResponse.status === 'fulfilled' && systemInfoResponse.value) {
      systemInfo.value = systemInfoResponse.value.data || systemInfoResponse.value
    }

    if (requestsResponse.status === 'fulfilled' && requestsResponse.value) {
      const requestsData = requestsResponse.value.data || requestsResponse.value
      recentRequests.value = Array.isArray(requestsData)
        ? requestsData
        : Array.isArray(requestsData?.items)
          ? requestsData.items
          : []
    }
  } catch (error) {
    console.error('获取系统统计信息失败:', error)
    Object.assign(stats, {
      total_users: 0,
      total_llm_resources: 0,
      total_requests: 0,
      success_rate: 0,
      avg_response_time: 0
    })
    recentRequests.value = []
    systemInfo.value = null
  } finally {
    refreshing.value = false
  }
}

const handleRefresh = () => {
  getSystemStats()
}

const navigateToRequests = () => {
  router.push('/requests')
}

onMounted(() => {
  getSystemStats()
})
</script>

<style scoped>
.dashboard {
  animation: fadeIn 0.4s ease-out;
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

.dashboard-header {
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

/* Claude Style Stats Cards */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: var(--claude-ivory);
  border: 1px solid var(--claude-border-cream);
  border-radius: var(--radius-comfortable);
  padding: 20px;
  transition: all 0.2s ease;
  box-shadow: var(--shadow-whisper);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: rgba(0, 0, 0, 0.08) 0px 8px 32px;
  border-color: var(--claude-border-warm);
}

.stat-card-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-comfortable);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
}

.stat-icon-terracotta {
  background: #f5e6df;
  color: var(--claude-terracotta);
}

.stat-icon-green {
  background: #e8f0e4;
  color: #7a9a6d;
}

.stat-icon-coral {
  background: #f5e8e6;
  color: var(--claude-coral);
}

.stat-icon-success {
  background: #e8f4e8;
  color: #6d9a7a;
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-family: var(--font-serif);
  font-size: 36px;
  font-weight: 500;
  line-height: 1.2;
  color: var(--claude-text-primary);
}

.stat-label {
  font-family: var(--font-sans);
  font-size: 14px;
  color: var(--claude-text-secondary);
}

/* Claude Style Performance Card */
.performance-card {
  margin-bottom: 24px;
}

.performance-metrics {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
}

.metric-card {
  background: #f8f7f3;
  border: 1px solid var(--claude-border-cream);
  border-radius: var(--radius-comfortable);
  padding: 16px;
  transition: all 0.2s ease;
}

.metric-card:hover {
  background: var(--claude-ivory);
  border-color: var(--claude-terracotta);
}

.metric-label {
  font-family: var(--font-sans);
  font-size: 14px;
  color: var(--claude-text-secondary);
  margin-bottom: 8px;
}

.metric-value {
  font-family: var(--font-sans);
  font-size: 20px;
  font-weight: 500;
  color: var(--claude-text-primary);
}

.metric-value-serif {
  font-family: var(--font-serif);
  font-size: 24px;
}

.metric-status {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #7a9a6d;
}

.status-dot {
  width: 8px;
  height: 8px;
  background: #7a9a6d;
  border-radius: 50%;
  box-shadow: 0 0 0 3px rgba(122, 154, 109, 0.2);
}

/* Claude Style Requests Card */
.requests-card {
  margin-bottom: 24px;
}

.card-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-family: var(--font-serif);
  font-size: 20px;
  font-weight: 500;
  line-height: 1.2;
  color: var(--claude-text-primary);
}

/* Responsive */
@media (max-width: 768px) {
  .stats-cards {
    grid-template-columns: 1fr;
  }

  .performance-metrics {
    grid-template-columns: repeat(2, 1fr);
  }

  .dashboard-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
}
</style>