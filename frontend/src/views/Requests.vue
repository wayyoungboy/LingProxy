<template>
  <div class="requests-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>请求管理</span>
          <el-button type="primary" @click="exportRequests">导出请求</el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="searchForm" class="mb-4">
        <el-form-item label="请求路径">
          <el-input v-model="searchForm.path" placeholder="请输入请求路径" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态">
            <el-option label="成功" value="success" />
            <el-option label="失败" value="error" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 250px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getRequests">查询</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
      
      <el-table :data="requestsList" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="path" label="请求路径" />
        <el-table-column prop="method" label="方法" width="100" />
        <el-table-column prop="model" label="模型" width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'success' ? 'success' : 'danger'">
              {{ scope.row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="response_time" label="响应时间(ms)" width="120" />
        <el-table-column prop="created_at" label="请求时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewRequestDetail(scope.row)">
              查看详情
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
      title="请求详情"
      width="800px"
    >
      <div v-if="currentRequest" class="request-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="请求ID">{{ currentRequest.id }}</el-descriptions-item>
          <el-descriptions-item label="请求路径">{{ currentRequest.path }}</el-descriptions-item>
          <el-descriptions-item label="请求方法">{{ currentRequest.method }}</el-descriptions-item>
          <el-descriptions-item label="使用模型">{{ currentRequest.model }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ currentRequest.status }}</el-descriptions-item>
          <el-descriptions-item label="响应时间">{{ currentRequest.response_time }}ms</el-descriptions-item>
          <el-descriptions-item label="请求时间">{{ currentRequest.created_at }}</el-descriptions-item>
          <el-descriptions-item label="请求参数">
            <el-scrollbar height="200px">
              <pre>{{ formatJson(currentRequest.request_params) }}</pre>
            </el-scrollbar>
          </el-descriptions-item>
          <el-descriptions-item label="响应数据">
            <el-scrollbar height="200px">
              <pre>{{ formatJson(currentRequest.response_data) }}</pre>
            </el-scrollbar>
          </el-descriptions-item>
          <el-descriptions-item label="错误信息" v-if="currentRequest.error_message">
            <el-scrollbar height="200px">
              <pre>{{ currentRequest.error_message }}</pre>
            </el-scrollbar>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="detailDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

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
      page: pagination.current,
      page_size: pagination.size,
      path: searchForm.path,
      status: searchForm.status
    }
    
    if (searchForm.dateRange && searchForm.dateRange.length === 2) {
      params.start_time = searchForm.dateRange[0]
      params.end_time = searchForm.dateRange[1]
    }
    
    const response = await api.getRequests(params)
    requestsList.value = response.data.items
    pagination.total = response.data.total
  } catch (error) {
    ElMessage.error('获取请求列表失败')
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
    ElMessage.error('获取请求详情失败')
  }
}

const exportRequests = async () => {
  try {
    await api.exportRequests(searchForm)
    ElMessage.success('导出请求成功')
  } catch (error) {
    ElMessage.error('导出请求失败')
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