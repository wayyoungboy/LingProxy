<template>
  <div class="dashboard">
    <el-card class="dashboard-card">
      <template #header>
        <div class="dashboard-header">
          <h2>系统仪表盘</h2>
          <el-button
            type="primary"
            size="small"
            @click="handleRefresh"
            :loading="refreshing"
          >
            <el-icon><Refresh /></el-icon>
            刷新数据
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
              <div class="stat-card-value">{{ stats.total_users }}</div>
              <div class="stat-card-label">总用户数</div>
            </div>
          </div>
        </el-card>
        
        <el-card class="stat-card stat-card-success" shadow="hover">
          <div class="stat-card-content">
            <div class="stat-card-icon">
              <el-icon><Cpu /></el-icon>
            </div>
            <div class="stat-card-info">
              <div class="stat-card-value">{{ stats.total_llm_resources }}</div>
              <div class="stat-card-label">LLM资源数</div>
            </div>
          </div>
        </el-card>
        
        <el-card class="stat-card stat-card-warning" shadow="hover">
          <div class="stat-card-content">
            <div class="stat-card-icon">
              <el-icon><Message /></el-icon>
            </div>
            <div class="stat-card-info">
              <div class="stat-card-value">{{ stats.total_requests }}</div>
              <div class="stat-card-label">总请求数</div>
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
              <div class="stat-card-label">成功率</div>
            </div>
          </div>
        </el-card>
      </div>
      
      <!-- 图表区域 -->
      <div class="charts-section">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="chart-header">
              <h3>系统性能</h3>
            </div>
          </template>
          <div class="chart-content">
            <div class="performance-metrics">
              <div class="metric-item">
                <div class="metric-label">平均响应时间</div>
                <div class="metric-value">{{ stats.avg_response_time }}ms</div>
              </div>
              <div class="metric-item">
                <div class="metric-label">今日请求数</div>
                <div class="metric-value">{{ todayRequests }}</div>
              </div>
              <div class="metric-item">
                <div class="metric-label">系统状态</div>
                <div class="metric-value status-online">正常运行</div>
              </div>
              <div class="metric-item">
                <div class="metric-label">运行时间</div>
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
            <div class="requests-header">
              <h3>最近请求</h3>
              <el-button
                type="primary"
                size="small"
                @click="navigateToRequests"
              >
                查看全部
              </el-button>
            </div>
          </template>
          
          <el-table
            :data="recentRequests"
            style="width: 100%"
            border
            size="small"
          >
            <el-table-column prop="id" label="请求ID" width="180" />
            <el-table-column prop="user_id" label="用户ID" width="180" />
            <el-table-column prop="endpoint" label="端点" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag
                  :type="scope.row.status === 'success' ? 'success' : 'danger'"
                >
                  {{ scope.row.status === 'success' ? '成功' : '失败' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="duration" label="响应时间" width="120">
              <template #default="scope">
                {{ scope.row.duration }}ms
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
          
          <div v-if="recentRequests.length === 0" class="no-data">
            <el-empty description="暂无请求记录" />
          </div>
        </el-card>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  User,
  Cpu,
  Message,
  Check,
  DataLine,
  Timer,
  WarningFilled,
  Monitor
} from '@element-plus/icons-vue'
import api from '../api'

const router = useRouter()
const refreshing = ref(false)
const stats = reactive({
  total_users: 0,
  total_llm_resources: 0,
  total_requests: 0,
  success_rate: 0,
  avg_response_time: 0
})
const recentRequests = ref([])

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
  // 从系统启动时间计算运行时间（如果有的话）
  // 暂时返回空，等待后端提供启动时间信息
  return ''
})

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 获取系统统计信息
const getSystemStats = async () => {
  try {
    refreshing.value = true
    
    // 从API获取系统统计信息
    const response = await api.getSystemStats()
    if (response && response.data) {
      Object.assign(stats, response.data)
    }
    
    // 获取最近的请求记录
    const requestsResponse = await api.getRequests({ limit: 10 })
    if (requestsResponse && requestsResponse.data) {
      // 只显示最近10条
      recentRequests.value = Array.isArray(requestsResponse.data) ? requestsResponse.data : []
    }
    
  } catch (error) {
    console.error('获取系统统计信息失败:', error)
    ElMessage.error('获取系统统计信息失败')
    
    // 显示真实的空数据状态
    Object.assign(stats, {
      total_users: 0,
      total_llm_resources: 0,
      total_requests: 0,
      success_rate: 0,
      avg_response_time: 0
    })
    recentRequests.value = []
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

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dashboard-header h2 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
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

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-header h3 {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
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

.requests-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.requests-header h3 {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
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
  
  .dashboard-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .chart-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .requests-header {
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