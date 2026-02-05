<template>
  <div class="dashboard">
    <el-card class="dashboard-card">
      <template #header>
        <div class="card-header">
          <span class="page-title">{{ $t('dashboard.systemDashboard') }}</span>
          <el-button
            type="primary"
            size="small"
            @click="handleRefresh"
            :loading="refreshing"
          >
            <el-icon><Refresh /></el-icon>
            {{ $t('dashboard.refreshData') }}
          </el-button>
        </div>
      </template>
      
      <!-- 统计卡片 -->
      <div class="stats-cards">
        <el-card class="stat-card stat-card-primary" shadow="hover">
          <div class="stat-card-content">
            <div class="stat-card-icon">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-card-info">
              <div class="stat-card-value">{{ formatNumber(stats.total_users) }}</div>
              <div class="stat-card-label">{{ $t('dashboard.totalUsers') }}</div>
            </div>
          </div>
        </el-card>
        
        <el-card class="stat-card stat-card-success" shadow="hover">
          <div class="stat-card-content">
            <div class="stat-card-icon">
              <el-icon><Cpu /></el-icon>
            </div>
            <div class="stat-card-info">
              <div class="stat-card-value">{{ formatNumber(stats.total_llm_resources) }}</div>
              <div class="stat-card-label">{{ $t('dashboard.totalLLMResources') }}</div>
            </div>
          </div>
        </el-card>
        
        <el-card class="stat-card stat-card-warning" shadow="hover">
          <div class="stat-card-content">
            <div class="stat-card-icon">
              <el-icon><Message /></el-icon>
            </div>
            <div class="stat-card-info">
              <div class="stat-card-value">{{ formatNumber(stats.total_requests) }}</div>
              <div class="stat-card-label">{{ $t('dashboard.totalRequests') }}</div>
            </div>
          </div>
        </el-card>
        
        <el-card class="stat-card stat-card-danger" shadow="hover">
          <div class="stat-card-content">
            <div class="stat-card-icon">
              <el-icon><Check /></el-icon>
            </div>
            <div class="stat-card-info">
              <div class="stat-card-value">{{ stats.success_rate }}%</div>
              <div class="stat-card-label">{{ $t('dashboard.successRate') }}</div>
            </div>
          </div>
        </el-card>
      </div>
      
      <!-- 图表区域 -->
      <div class="charts-section">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="page-title">{{ $t('dashboard.systemPerformance') }}</span>
            </div>
          </template>
          <div class="chart-content">
            <div class="performance-metrics">
              <div class="metric-item">
                <div class="metric-label">{{ $t('dashboard.avgResponseTime') }}</div>
                <div class="metric-value">{{ stats.avg_response_time }}ms</div>
              </div>
              <div class="metric-item">
                <div class="metric-label">{{ $t('dashboard.todayRequests') }}</div>
                <div class="metric-value">{{ formatNumber(todayRequests) }}</div>
              </div>
              <div class="metric-item">
                <div class="metric-label">{{ $t('dashboard.systemStatus') }}</div>
                <div class="metric-value status-online">{{ $t('dashboard.running') }}</div>
              </div>
              <div class="metric-item">
                <div class="metric-label">{{ $t('dashboard.runningTime') }}</div>
                <div class="metric-value">{{ uptime || '--' }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </div>
      
      <!-- 最近请求 -->
      <div class="recent-requests-section">
        <el-card class="requests-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="page-title">{{ $t('dashboard.recentRequests') }}</span>
              <el-button
                type="primary"
                size="small"
                @click="navigateToRequests"
              >
                {{ $t('dashboard.viewAll') }}
              </el-button>
            </div>
          </template>
          
          <el-table
            :data="recentRequests"
            style="width: 100%"
            border
            stripe
            size="small"
          >
            <el-table-column prop="id" :label="$t('requests.requestId')" width="180" />
            <el-table-column prop="user_id" :label="$t('dashboard.userId')" width="180" />
            <el-table-column prop="endpoint" :label="$t('llmResources.endpoint')" />
            <el-table-column prop="status" :label="$t('common.status')" width="100">
              <template #default="scope">
                <el-tag
                  :type="scope.row.status === 'success' ? 'success' : 'danger'"
                >
                  {{ scope.row.status === 'success' ? $t('requests.success') : $t('requests.failed') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="duration" :label="$t('requests.duration')" width="120">
              <template #default="scope">
                {{ scope.row.duration }}ms
              </template>
            </el-table-column>
            <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
          
          <div v-if="recentRequests.length === 0" class="no-data">
            <el-empty :description="$t('dashboard.noRequests')" />
          </div>
        </el-card>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  User,
  Cpu,
  Message,
  Check
} from '@element-plus/icons-vue'
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
  // 从最近请求中计算今日请求数
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return recentRequests.value.filter(req => {
    const reqDate = new Date(req.created_at)
    return reqDate >= today
  }).length
})

const uptime = computed(() => {
  // 从系统信息中获取运行时间
  return systemInfo.value?.uptime || '--'
})

// 获取系统统计信息
const getSystemStats = async () => {
  try {
    refreshing.value = true
    
    // 并行请求多个接口以提高性能
    const [statsResponse, systemInfoResponse, requestsResponse] = await Promise.allSettled([
      api.getSystemStats(),
      api.getSystemInfo(),
      api.getRequests({ limit: 10 })
    ])
    
    // 处理统计信息
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
    
    // 处理系统信息
    if (systemInfoResponse.status === 'fulfilled' && systemInfoResponse.value) {
      systemInfo.value = systemInfoResponse.value.data || systemInfoResponse.value
    }
    
    // 处理请求记录
    if (requestsResponse.status === 'fulfilled' && requestsResponse.value) {
      const requestsData = requestsResponse.value.data || requestsResponse.value
      recentRequests.value = Array.isArray(requestsData) ? requestsData : 
                           (Array.isArray(requestsData?.items) ? requestsData.items : [])
    }
    
  } catch (error) {
    console.error('获取系统统计信息失败:', error)
    // 错误信息已在API拦截器中处理，这里只重置数据
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



// 处理刷新数据
const handleRefresh = () => {
  getSystemStats()
}

// 跳转到请求管理页面
const navigateToRequests = () => {
  router.push('/requests')
}

// 组件挂载时获取数据
onMounted(() => {
  getSystemStats()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.dashboard-card {
  margin-bottom: 20px;
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

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-card-primary {
  border-left: 4px solid #409EFF;
}

.stat-card-success {
  border-left: 4px solid #67C23A;
}

.stat-card-warning {
  border-left: 4px solid #E6A23C;
}

.stat-card-danger {
  border-left: 4px solid #F56C6C;
}

.stat-card-content {
  display: flex;
  align-items: center;
  padding: 20px;
}

.stat-card-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  font-size: 24px;
}

.stat-card-primary .stat-card-icon {
  background-color: rgba(64, 158, 255, 0.1);
  color: #409EFF;
}

.stat-card-success .stat-card-icon {
  background-color: rgba(103, 194, 58, 0.1);
  color: #67C23A;
}

.stat-card-warning .stat-card-icon {
  background-color: rgba(230, 162, 60, 0.1);
  color: #E6A23C;
}

.stat-card-danger .stat-card-icon {
  background-color: rgba(245, 108, 108, 0.1);
  color: #F56C6C;
}

.stat-card-info {
  flex: 1;
}

.stat-card-value {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 4px;
  color: #303133;
}

.stat-card-label {
  font-size: 14px;
  color: #606266;
}

/* 图表区域 */
.charts-section {
  margin-bottom: 24px;
}

.chart-card {
  border-radius: 8px;
  overflow: hidden;
}


.chart-content {
  padding: 20px;
}

.performance-metrics {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 20px;
}

.metric-item {
  text-align: center;
  padding: 16px;
  background-color: #f9f9f9;
  border-radius: 8px;
}

.metric-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
}

.metric-value {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.status-online {
  color: #67C23A;
}

/* 最近请求 */
.recent-requests-section {
  margin-bottom: 20px;
}

.requests-card {
  border-radius: 8px;
  overflow: hidden;
}


.no-data {
  padding: 40px 0;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-cards {
    grid-template-columns: 1fr;
  }
  
  .performance-metrics {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}

@media (max-width: 480px) {
  .performance-metrics {
    grid-template-columns: 1fr;
  }
  
  .stat-card-content {
    padding: 16px;
  }
  
  .stat-card-value {
    font-size: 20px;
  }
}
</style>