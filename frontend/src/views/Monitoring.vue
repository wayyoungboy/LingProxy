<template>
  <div class="monitoring">
    <div class="monitoring-header">
      <h1 class="page-title">{{ $t('monitoring.title') }}</h1>
      <div class="header-actions">
        <el-select
          v-model="selectedPeriod"
          @change="fetchMonitorData"
          style="width: 160px"
        >
          <el-option :label="$t('monitoring.oneHour')" value="1h" />
          <el-option :label="$t('monitoring.sixHours')" value="6h" />
          <el-option :label="$t('monitoring.twentyFourHours')" value="24h" />
          <el-option :label="$t('monitoring.sevenDays')" value="7d" />
        </el-select>
        <el-button type="primary" @click="fetchMonitorData" :loading="loading">
          <el-icon><Refresh /></el-icon>
          {{ $t('dashboard.refreshData') }}
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-card-content">
          <div class="stat-icon stat-icon-terracotta">
            <el-icon><Message /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ formatNumber(monitorData.totalRequests) }}</div>
            <div class="stat-label">{{ $t('monitoring.totalRequests') }}</div>
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-card-content">
          <div class="stat-icon stat-icon-success">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ monitorData.successRate?.toFixed(1) || '--' }}%</div>
            <div class="stat-label">{{ $t('monitoring.successRate') }}</div>
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-card-content">
          <div class="stat-icon stat-icon-green">
            <el-icon><Timer /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ monitorData.avgResponseTime || '--' }}ms</div>
            <div class="stat-label">{{ $t('monitoring.avgResponseTime') }}</div>
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-card-content">
          <div class="stat-icon stat-icon-coral">
            <el-icon><Files /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ formatNumber(monitorData.totalTokens) }}</div>
            <div class="stat-label">{{ $t('monitoring.totalTokens') }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 限流器状态 -->
    <el-card v-if="monitorData.rateLimiter" class="rate-limiter-card">
      <template #header>
        <span class="section-title">{{ $t('monitoring.rateLimiterStatus') }}</span>
      </template>
      <div class="rate-limiter-stats">
        <div class="rate-item">
          <span class="rate-label">{{ $t('monitoring.enabled') }}</span>
          <el-tag :type="monitorData.rateLimiter.enabled ? 'success' : 'info'" size="small">
            {{ monitorData.rateLimiter.enabled ? $t('monitoring.enabled') : $t('monitoring.disabled') }}
          </el-tag>
        </div>
        <div class="rate-item">
          <span class="rate-label">{{ $t('monitoring.rpmLimit') }}</span>
          <span class="rate-value">{{ monitorData.rateLimiter.rpm_limit || '--' }}</span>
        </div>
        <div class="rate-item">
          <span class="rate-label">{{ $t('monitoring.refillRate') }}</span>
          <span class="rate-value">{{ (monitorData.rateLimiter.refill_rate || 0).toFixed(1) }}/s</span>
        </div>
        <div class="rate-item">
          <span class="rate-label">{{ $t('monitoring.concurrency') }}</span>
          <span class="rate-value">{{ monitorData.rateLimiter.concurrency || '--' }}</span>
        </div>
        <div class="rate-item">
          <span class="rate-label">{{ $t('monitoring.activeClients') }}</span>
          <span class="rate-value">{{ monitorData.rateLimiter.active_clients || 0 }}</span>
        </div>
      </div>
    </el-card>

    <!-- 图表区域 -->
    <div v-if="hasData" class="charts-grid">
      <el-card class="chart-card">
        <template #header>
          <span class="section-title">{{ $t('monitoring.requestsOverTime') }}</span>
        </template>
        <div ref="requestsChartRef" class="chart-container"></div>
      </el-card>

      <el-card class="chart-card">
        <template #header>
          <span class="section-title">{{ $t('monitoring.successRateOverTime') }}</span>
        </template>
        <div ref="successRateChartRef" class="chart-container"></div>
      </el-card>

      <el-card class="chart-card">
        <template #header>
          <span class="section-title">{{ $t('monitoring.avgDurationOverTime') }}</span>
        </template>
        <div ref="durationChartRef" class="chart-container"></div>
      </el-card>

      <el-card class="chart-card">
        <template #header>
          <span class="section-title">{{ $t('monitoring.tokensOverTime') }}</span>
        </template>
        <div ref="tokensChartRef" class="chart-container"></div>
      </el-card>
    </div>

    <el-empty v-else :description="$t('monitoring.noData')" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, computed, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'
import { Refresh, Message, Check, Timer, Files } from '@element-plus/icons-vue'
import api from '../api'
import { formatNumber } from '../utils/index'

const { t } = useI18n()

const loading = ref(false)
const selectedPeriod = ref('1h')
const timeline = ref([])
const rateLimiter = ref(null)
const chartInstances = ref([])

const requestsChartRef = ref(null)
const successRateChartRef = ref(null)
const durationChartRef = ref(null)
const tokensChartRef = ref(null)

const hasData = computed(() => timeline.value.length > 0)

const monitorData = computed(() => {
  if (!timeline.value.length) {
    return { totalRequests: 0, successRate: 0, avgResponseTime: 0, totalTokens: 0, rateLimiter: null }
  }

  let totalReqs = 0
  let totalSuccess = 0
  let totalFail = 0
  let totalTokensVal = 0
  let totalDur = 0

  for (const b of timeline.value) {
    totalReqs += b.total_requests
    totalSuccess += b.success_requests
    totalFail += b.failed_requests
    totalTokensVal += b.total_tokens
    totalDur += b.total_duration_ms
  }

  const avgDur = totalReqs > 0 ? Math.round(totalDur / totalReqs) : 0
  const successRate = totalReqs > 0 ? (totalSuccess / totalReqs * 100) : 0

  return {
    totalRequests: totalReqs,
    successRate,
    avgResponseTime: avgDur,
    totalTokens: totalTokensVal,
    rateLimiter: rateLimiter.value
  }
})

const fetchMonitorData = async () => {
  try {
    loading.value = true
    const response = await api.getMonitorStats(selectedPeriod.value)
    const data = response?.data || response
    timeline.value = data.timeline || []
    rateLimiter.value = data.rate_limiter || null
    await nextTick()
    renderCharts()
  } catch (error) {
    console.error('获取监控数据失败:', error)
    timeline.value = []
    rateLimiter.value = null
  } finally {
    loading.value = false
  }
}

const renderCharts = () => {
  disposeCharts()
  if (!hasData.value) return

  const data = timeline.value
  const times = data.map(d => {
    const date = new Date(d.timestamp)
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  })

  createChart(requestsChartRef.value, buildRequestChart(times, data))
  createChart(successRateChartRef.value, buildSuccessRateChart(times, data))
  createChart(durationChartRef.value, buildDurationChart(times, data))
  createChart(tokensChartRef.value, buildTokensChart(times, data))
}

const createChart = (domElement, option) => {
  if (!domElement) return
  const chart = echarts.init(domElement)
  chart.setOption(option)
  chartInstances.value.push(chart)
  window.addEventListener('resize', () => chart.resize())
}

const disposeCharts = () => {
  chartInstances.value.forEach(c => c.dispose())
  chartInstances.value = []
}

const buildRequestChart = (times, data) => ({
  tooltip: { trigger: 'axis' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: times, boundaryGap: false },
  yAxis: { type: 'value' },
  series: [
    {
      name: t('monitoring.successRequests'),
      type: 'line',
      stack: 'Total',
      areaStyle: { color: '#e8f4e8' },
      data: data.map(d => d.success_requests),
      lineStyle: { color: '#6d9a7a' },
      itemStyle: { color: '#6d9a7a' },
      smooth: true
    },
    {
      name: t('monitoring.failedRequests'),
      type: 'line',
      stack: 'Total',
      areaStyle: { color: '#f5e6df' },
      data: data.map(d => d.failed_requests),
      lineStyle: { color: '#d4756d' },
      itemStyle: { color: '#d4756d' },
      smooth: true
    }
  ]
})

const buildSuccessRateChart = (times, data) => ({
  tooltip: { trigger: 'axis', formatter: p => `${p[0].name}<br/>${p[0].marker}${p[0].seriesName}: ${p[0].value.toFixed(1)}%` },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: times, boundaryGap: false },
  yAxis: { type: 'value', max: 100 },
  series: [{
    name: t('monitoring.successRate'),
    type: 'line',
    areaStyle: { color: '#f0f4ff' },
    data: data.map(d => d.success_rate),
    lineStyle: { color: '#7a9a6d' },
    itemStyle: { color: '#7a9a6d' },
    smooth: true
  }]
})

const buildDurationChart = (times, data) => ({
  tooltip: { trigger: 'axis', formatter: p => `${p[0].name}<br/>${p[0].marker}${p[0].seriesName}: ${p[0].value}ms` },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: times, boundaryGap: false },
  yAxis: { type: 'value', name: 'ms' },
  series: [{
    name: t('monitoring.avgResponseTime'),
    type: 'line',
    areaStyle: { color: '#fff5f0' },
    data: data.map(d => d.avg_duration_ms),
    lineStyle: { color: '#d4756d' },
    itemStyle: { color: '#d4756d' },
    smooth: true
  }]
})

const buildTokensChart = (times, data) => ({
  tooltip: { trigger: 'axis' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: times, boundaryGap: false },
  yAxis: { type: 'value' },
  series: [{
    name: t('monitoring.totalTokens'),
    type: 'bar',
    areaStyle: { color: '#e8f0e4' },
    data: data.map(d => d.total_tokens),
    itemStyle: { color: '#7a9a6d' },
    smooth: true
  }]
})

onMounted(() => {
  fetchMonitorData()
})

onUnmounted(() => {
  disposeCharts()
})
</script>

<style scoped>
.monitoring {
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.monitoring-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.page-title {
  font-family: var(--font-serif);
  font-size: 32px;
  font-weight: 500;
  line-height: 1.1;
  color: var(--claude-text-primary);
}

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

.stat-icon-terracotta { background: #f5e6df; color: var(--claude-terracotta); }
.stat-icon-green { background: #e8f0e4; color: #7a9a6d; }
.stat-icon-coral { background: #f5e8e6; color: var(--claude-coral); }
.stat-icon-success { background: #e8f4e8; color: #6d9a7a; }

.stat-info { flex: 1; }

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

.rate-limiter-card { margin-bottom: 24px; }

.rate-limiter-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
}

.rate-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.rate-label {
  font-size: 14px;
  color: var(--claude-text-secondary);
}

.rate-value {
  font-family: var(--font-serif);
  font-size: 20px;
  font-weight: 500;
  color: var(--claude-text-primary);
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(480px, 1fr));
  gap: 16px;
}

.chart-card { margin-bottom: 8px; }

.chart-container {
  width: 100%;
  height: 300px;
}

.section-title {
  font-family: var(--font-serif);
  font-size: 20px;
  font-weight: 500;
  line-height: 1.2;
  color: var(--claude-text-primary);
}

@media (max-width: 768px) {
  .stats-cards { grid-template-columns: 1fr; }
  .charts-grid { grid-template-columns: 1fr; }
  .monitoring-header { flex-direction: column; gap: 16px; align-items: flex-start; }
  .rate-limiter-stats { grid-template-columns: repeat(2, 1fr); }
}
</style>
