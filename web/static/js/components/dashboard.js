// 仪表盘组件
const DashboardView = {
    props: ['stats', 'charts'],
    template: `
        <div>
            <!-- 统计卡片 -->
            <a-row :gutter="16">
                <a-col :span="6">
                    <a-card class="stat-card">
                        <div class="stat-icon" style="color: #1890ff;">
                            <user-outlined />
                        </div>
                        <div class="stat-value">{{ stats.totalUsers }}</div>
                        <div class="stat-label">总用户数</div>
                    </a-card>
                </a-col>
                <a-col :span="6">
                    <a-card class="stat-card">
                        <div class="stat-icon" style="color: #52c41a;">
                            <api-outlined />
                        </div>
                        <div class="stat-value">{{ stats.totalProviders }}</div>
                        <div class="stat-label">提供商数量</div>
                    </a-card>
                </a-col>
                <a-col :span="6">
                    <a-card class="stat-card">
                        <div class="stat-icon" style="color: #faad14;">
                            <history-outlined />
                        </div>
                        <div class="stat-value">{{ stats.totalRequests }}</div>
                        <div class="stat-label">总请求数</div>
                    </a-card>
                </a-col>
                <a-col :span="6">
                    <a-card class="stat-card">
                        <div class="stat-icon" style="color: #722ed1;">
                            <check-circle-outlined />
                        </div>
                        <div class="stat-value">{{ stats.successRate }}%</div>
                        <div class="stat-label">成功率</div>
                    </a-card>
                </a-col>
            </a-row>

            <!-- 图表区域 -->
            <a-row :gutter="16" style="margin-top: 24px;">
                <a-col :span="12">
                    <a-card title="每日请求量" class="chart-container">
                        <div id="daily-chart" style="height: 300px;"></div>
                    </a-card>
                </a-col>
                <a-col :span="12">
                    <a-card title="提供商使用分布" class="chart-container">
                        <div id="provider-chart" style="height: 300px;"></div>
                    </a-card>
                </a-col>
            </a-row>

            <!-- 实时状态 -->
            <a-row :gutter="16" style="margin-top: 24px;">
                <a-col :span="24">
                    <a-card title="实时状态">
                        <a-table
                            :columns="statusColumns"
                            :data-source="statusData"
                            :pagination="false"
                            size="small"
                        >
                            <template #bodyCell="{ column, record }">
                                <template v-if="column.key === 'status'">
                                    <a-tag :color="record.status === '正常' ? 'green' : 'red'">
                                        {{ record.status }}
                                    </a-tag>
                                </template>
                            </template>
                        </a-table>
                    </a-card>
                </a-col>
            </a-row>
        </div>
    `,
    data() {
        return {
            statusColumns: [
                { title: '服务', dataIndex: 'service', key: 'service' },
                { title: '状态', dataIndex: 'status', key: 'status' },
                { title: '响应时间', dataIndex: 'responseTime', key: 'responseTime' },
                { title: '最后检查', dataIndex: 'lastCheck', key: 'lastCheck' }
            ],
            statusData: [
                { service: 'OpenAI API', status: '正常', responseTime: '120ms', lastCheck: '2分钟前' },
                { service: 'Claude API', status: '正常', responseTime: '200ms', lastCheck: '1分钟前' },
                { service: 'Gemini API', status: '正常', responseTime: '150ms', lastCheck: '3分钟前' }
            ]
        }
    },
    mounted() {
        this.initCharts()
    },
    methods: {
        initCharts() {
            // 初始化AntV图表
            console.log('初始化仪表盘图表...')
        }
    }
}