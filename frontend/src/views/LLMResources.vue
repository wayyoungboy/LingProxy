<template>
  <div class="llm-resources">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>LLM资源管理</span>
          <div class="header-actions">
            <el-button type="success" @click="handleDownloadTemplate">
              <el-icon><Download /></el-icon>
              下载导入模板
            </el-button>
            <el-upload
              ref="uploadRef"
              :action="uploadAction"
              :headers="uploadHeaders"
              :on-success="handleImportSuccess"
              :on-error="handleImportError"
              :before-upload="beforeUpload"
              :show-file-list="false"
              accept=".xlsx,.xls"
            >
              <el-button type="warning">
                <el-icon><Upload /></el-icon>
                批量导入
              </el-button>
            </el-upload>
            <el-button type="primary" @click="handleAddResource">
              <el-icon><Plus /></el-icon>
              添加LLM资源
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 搜索和筛选 -->
      <div class="search-filter">
        <el-input
          v-model="searchQuery"
          placeholder="搜索资源名称或基础URL"
          :prefix-icon="Search"
          style="width: 300px; margin-right: 10px"
        ></el-input>
        <el-select
          v-model="typeFilter"
          placeholder="筛选模型类别"
          style="width: 140px; margin-right: 10px"
        >
          <el-option label="全部" value=""></el-option>
          <el-option label="对话" value="chat"></el-option>
          <el-option label="生图" value="image"></el-option>
          <el-option label="嵌入" value="embedding"></el-option>
          <el-option label="重排序" value="rerank"></el-option>
          <el-option label="语音" value="audio"></el-option>
          <el-option label="视频" value="video"></el-option>
        </el-select>
        <el-select
          v-model="statusFilter"
          placeholder="筛选状态"
          style="width: 120px"
        >
          <el-option label="全部" value=""></el-option>
          <el-option label="活跃" value="active"></el-option>
          <el-option label="禁用" value="inactive"></el-option>
        </el-select>
      </div>
      
      <!-- LLM资源列表 -->
      <el-table
        v-loading="loading"
        :data="filteredResources"
        style="width: 100%; margin-top: 20px"
        border
      >
        <el-table-column prop="id" label="ID" width="180" />
        <el-table-column prop="name" label="资源名称" />
        <el-table-column prop="type" label="模型类别" width="120">
          <template #default="scope">
            <el-tag>{{ getTypeLabel(scope.row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="driver" label="驱动" width="100">
          <template #default="scope">
            <el-tag type="info">{{ getDriverLabel(scope.row.driver) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="model" label="模型标识" width="150">
          <template #default="scope">
            <el-tag type="warning" v-if="scope.row.model">{{ scope.row.model }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="base_url" label="基础URL" />
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag
              :type="scope.row.status === 'active' ? 'success' : 'danger'"
            >
              {{ scope.row.status === 'active' ? '活跃' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="handleEditResource(scope.row)"
              style="margin-right: 5px"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDeleteResource(scope.row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 添加/编辑LLM资源对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
    >
      <el-form
        ref="resourceFormRef"
        :model="resourceForm"
        :rules="resourceRules"
        label-width="100px"
      >
        <el-form-item label="资源名称" prop="name">
          <el-input v-model="resourceForm.name" placeholder="请输入资源名称"></el-input>
        </el-form-item>
        <el-form-item label="模型类别" prop="type">
          <el-select 
            v-model="resourceForm.type" 
            placeholder="请选择模型类别"
          >
            <el-option label="对话" value="chat"></el-option>
            <el-option label="生图" value="image"></el-option>
            <el-option label="嵌入" value="embedding"></el-option>
            <el-option label="重排序" value="rerank"></el-option>
            <el-option label="语音" value="audio"></el-option>
            <el-option label="视频" value="video"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="驱动" prop="driver">
          <el-select 
            v-model="resourceForm.driver" 
            placeholder="请选择驱动"
            @change="handleDriverChange"
          >
            <el-option label="OpenAI" value="openai"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="模型标识" prop="model">
          <el-input 
            v-model="resourceForm.model" 
            placeholder="请输入模型标识，如：gpt-4, gpt-3.5-turbo（此资源对应的模型标识）"
          ></el-input>
        </el-form-item>
        <el-form-item label="基础URL" prop="base_url">
          <el-input 
            v-model="resourceForm.base_url" 
            placeholder="选择驱动后将自动填充"
          ></el-input>
        </el-form-item>
        <el-form-item label="API密钥" prop="api_key">
          <el-input
            v-model="resourceForm.api_key"
            placeholder="请输入API密钥"
            :type="apiKeyVisible ? 'text' : 'password'"
          >
            <template #append>
              <el-button @click="apiKeyVisible = !apiKeyVisible" type="text">
                {{ apiKeyVisible ? '隐藏' : '显示' }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="resourceForm.status" placeholder="请选择状态">
            <el-option label="活跃" value="active"></el-option>
            <el-option label="禁用" value="inactive"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveResource">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, View, Hide, Upload, Download } from '@element-plus/icons-vue'
import api from '../api'

const router = useRouter()
const loading = ref(false)
const resources = ref([])
const searchQuery = ref('')
const typeFilter = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const uploadRef = ref(null)

// 上传配置
const uploadAction = '/api/v1/llm-resources/import'
const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
})

// 对话框相关
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isAddMode = ref(false)
const resourceFormRef = ref(null)
const apiKeyVisible = ref(false)
const resourceForm = reactive({
  id: '',
  name: '',
  type: 'chat',
  driver: 'openai',
  model: '',
  base_url: '',
  api_key: '',
  status: 'active'
})

// 表单验证规则
const resourceRules = {
  name: [
    { required: true, message: '请输入资源名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择模型类别', trigger: 'change' }
  ],
  driver: [
    { required: true, message: '请选择驱动', trigger: 'change' }
  ],
  model: [
    { required: true, message: '请输入模型标识', trigger: 'blur' }
  ],
  base_url: [
    { required: true, message: '请输入基础URL', trigger: 'blur' }
  ],
  api_key: [
    { required: true, message: '请输入API密钥', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 获取LLM资源列表
const getResourceList = async () => {
  try {
    loading.value = true
    const response = await api.getLLMResources({ page: currentPage.value, page_size: pageSize.value })
    if (response && response.data) {
      // 后端返回的data可能是数组，也可能是分页对象
      if (Array.isArray(response.data)) {
        // 直接是数组格式
        resources.value = response.data
        total.value = response.data.length
      } else if (response.data.items) {
        // 分页格式
        resources.value = response.data.items || []
        total.value = response.data.total || 0
      } else {
        // 其他格式，尝试直接使用data
        resources.value = []
        total.value = 0
      }
    }
  } catch (error) {
    console.error('获取LLM资源列表失败:', error)
    ElMessage.error('获取LLM资源列表失败')
  } finally {
    loading.value = false
  }
}

// 过滤资源
const filteredResources = computed(() => {
  let result = resources.value
  
  // 搜索过滤（支持资源名称、基础URL和模型标识的模糊搜索）
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase().trim()
    result = result.filter(resource => {
      // 搜索资源名称
      const nameMatch = resource.name && resource.name.toLowerCase().includes(query)
      // 搜索基础URL（模糊匹配）
      const urlMatch = resource.base_url && resource.base_url.toLowerCase().includes(query)
      // 搜索模型标识（可选，增强搜索体验）
      const modelMatch = resource.model && resource.model.toLowerCase().includes(query)
      // 只要有一个字段匹配就返回true
      return nameMatch || urlMatch || modelMatch
    })
  }
  
  // 类型过滤
  if (typeFilter.value) {
    result = result.filter(resource => resource.type === typeFilter.value)
  }
  
  // 状态过滤
  if (statusFilter.value) {
    result = result.filter(resource => resource.status === statusFilter.value)
  }
  
  // 分页
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return result.slice(startIndex, endIndex)
})

// 模型类别标签映射
const typeLabels = {
  chat: '对话',
  image: '生图',
  embedding: '嵌入',
  rerank: '重排序',
  audio: '语音',
  video: '视频'
}

// 获取模型类别标签
const getTypeLabel = (type) => {
  return typeLabels[type] || type
}

// 驱动标签映射
const driverLabels = {
  openai: 'OpenAI'
}

// 获取驱动标签
const getDriverLabel = (driver) => {
  return driverLabels[driver] || driver
}

// 驱动到BaseURL的映射
const driverToBaseURL = {
  openai: 'https://api.openai.com/v1'
}

// 处理驱动变化
const handleDriverChange = (driver) => {
  if (!driver) return
  
  // 根据驱动自动填充BaseURL
  const defaultURL = driverToBaseURL[driver]
  if (defaultURL) {
    // 只有在添加模式，或者当前base_url为空，或者是默认URL时才自动填充
    if (isAddMode.value || !resourceForm.base_url || Object.values(driverToBaseURL).includes(resourceForm.base_url)) {
      resourceForm.base_url = defaultURL
    }
  }
}

// 处理添加资源
const handleAddResource = () => {
  isAddMode.value = true
  dialogTitle.value = '添加LLM资源'
  // 重置表单
  Object.assign(resourceForm, {
    id: '',
    name: '',
    type: 'chat',
    driver: 'openai',
    model: '',
    base_url: driverToBaseURL['openai'], // 默认填充OpenAI的URL
    api_key: '',
    status: 'active'
  })
  apiKeyVisible.value = false
  dialogVisible.value = true
}

// 处理编辑资源
const handleEditResource = (resource) => {
  isAddMode.value = false
  dialogTitle.value = '编辑LLM资源'
  // 填充表单
  Object.assign(resourceForm, resource)
  apiKeyVisible.value = false
  dialogVisible.value = true
}

// 处理保存资源
const handleSaveResource = async () => {
  try {
    // 验证表单
    await resourceFormRef.value.validate()
    
    if (isAddMode.value) {
      // 创建资源
      await api.createLLMResource(resourceForm)
      ElMessage.success('LLM资源创建成功')
    } else {
      // 更新资源
      await api.updateLLMResource(resourceForm.id, resourceForm)
      ElMessage.success('LLM资源更新成功')
    }
    
    // 关闭对话框
    dialogVisible.value = false
    // 重新获取资源列表
    getResourceList()
  } catch (error) {
    console.error('保存LLM资源失败:', error)
    ElMessage.error('保存LLM资源失败')
  }
}

// 处理删除资源
const handleDeleteResource = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个LLM资源吗？',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'danger'
      }
    )
    
    await api.deleteLLMResource(id)
    ElMessage.success('LLM资源删除成功')
    // 重新获取资源列表
    getResourceList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除LLM资源失败:', error)
      ElMessage.error('删除LLM资源失败')
    }
  }
}

// 查看资源下的模型

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (current) => {
  currentPage.value = current
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 下载导入模板
const handleDownloadTemplate = async () => {
  try {
    const response = await api.downloadLLMResourcesTemplate()
    
    // 响应拦截器对于blob响应返回的是整个response对象
    // response.data应该是Blob类型
    if (!response || !response.data) {
      ElMessage.error('下载失败：响应格式不正确')
      return
    }
    
    const blob = response.data instanceof Blob ? response.data : new Blob([response.data])
    
    // 验证blob是否有效
    if (!blob || blob.size === 0) {
      ElMessage.error('下载的文件为空，请重试')
      return
    }
    
    // 验证文件类型（Excel文件应该以ZIP格式开头）
    const firstBytes = await blob.slice(0, 4).arrayBuffer()
    const uint8Array = new Uint8Array(firstBytes)
    
    // Excel文件（.xlsx）是ZIP格式，ZIP文件头是 50 4B 03 04 (PK..)
    const isValidExcel = uint8Array[0] === 0x50 && uint8Array[1] === 0x4B && 
                          uint8Array[2] === 0x03 && uint8Array[3] === 0x04
    
    if (!isValidExcel) {
      // 可能是错误响应，尝试读取错误信息
      const text = await blob.text()
      try {
        const errorData = JSON.parse(text)
        ElMessage.error(errorData.error || '下载模板失败')
      } catch (e) {
        ElMessage.error('下载的文件格式不正确，请重试')
      }
      return
    }
    
    // 创建下载链接
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'llm_resources_import_template.xlsx'
    link.style.display = 'none'
    document.body.appendChild(link)
    link.click()
    
    // 清理
    setTimeout(() => {
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    }, 100)
    
    ElMessage.success('模板下载成功')
  } catch (error) {
    console.error('下载模板失败:', error)
    ElMessage.error(error.response?.data?.error || error.message || '下载模板失败')
  }
}

// 上传前验证
const beforeUpload = (file) => {
  const isExcel = file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
                  file.type === 'application/vnd.ms-excel' ||
                  file.name.endsWith('.xlsx') ||
                  file.name.endsWith('.xls')
  if (!isExcel) {
    ElMessage.error('只能上传Excel文件（.xlsx或.xls格式）')
    return false
  }
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过10MB')
    return false
  }
  return true
}

// 导入成功
const handleImportSuccess = (response) => {
  if (response && response.success !== undefined) {
    const message = `导入完成！成功: ${response.success}条，失败: ${response.fail}条`
    if (response.fail > 0 && response.errors && response.errors.length > 0) {
      ElMessageBox.alert(
        `${message}\n\n错误详情：\n${response.errors.slice(0, 10).join('\n')}${response.errors.length > 10 ? '\n...' : ''}`,
        '导入结果',
        {
          confirmButtonText: '确定',
          type: response.success > 0 ? 'warning' : 'error'
        }
      )
    } else {
      ElMessage.success(message)
    }
    // 刷新列表
    getResourceList()
  } else {
    ElMessage.success('导入成功')
    getResourceList()
  }
}

// 导入失败
const handleImportError = (error) => {
  console.error('导入失败:', error)
  const errorMsg = error.response?.data?.error || '导入失败，请检查文件格式'
  ElMessage.error(errorMsg)
}

// 组件挂载时获取数据
onMounted(() => {
  getResourceList()
})
</script>

<style scoped>
.llm-resources {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.search-filter {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}
</style>