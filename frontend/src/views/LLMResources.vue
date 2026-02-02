<template>
  <div class="llm-resources">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>LLM资源管理</span>
          <el-button type="primary" @click="handleAddResource">
            <el-icon><Plus /></el-icon>
            添加LLM资源
          </el-button>
        </div>
      </template>
      
      <!-- 搜索和筛选 -->
      <div class="search-filter">
        <el-input
          v-model="searchQuery"
          placeholder="搜索资源名称"
          :prefix-icon="Search"
          style="width: 240px; margin-right: 10px"
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
        <el-table-column prop="provider" label="服务提供商" width="120">
          <template #default="scope">
            <el-tag type="info">{{ getProviderLabel(scope.row.provider) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="model" label="默认模型" width="150">
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
              type="success"
              size="small"
              @click="viewModels(scope.row.id)"
              style="margin-right: 5px"
            >
              查看模型
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
        <el-form-item label="服务提供商" prop="provider">
          <el-select 
            v-model="resourceForm.provider" 
            placeholder="请选择服务提供商"
            @change="handleProviderChange"
          >
            <el-option label="OpenAI" value="openai"></el-option>
            <el-option label="Zai" value="zai"></el-option>
            <el-option label="Anthropic" value="anthropic"></el-option>
            <el-option label="Google" value="google"></el-option>
            <el-option label="Azure" value="azure"></el-option>
            <el-option label="自定义" value="custom"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="默认模型" prop="model">
          <el-input 
            v-model="resourceForm.model" 
            placeholder="请输入默认模型名称，如：gpt-4, gpt-3.5-turbo"
          ></el-input>
        </el-form-item>
        <el-form-item label="基础URL" prop="base_url">
          <el-input 
            v-model="resourceForm.base_url" 
            placeholder="选择服务提供商后将自动填充"
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
import { Plus, Search, View, Hide } from '@element-plus/icons-vue'
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
  provider: 'openai',
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
  provider: [
    { required: true, message: '请选择服务提供商', trigger: 'change' }
  ],
  model: [
    { required: true, message: '请输入默认模型名称', trigger: 'blur' }
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
      resources.value = response.data.items || []
      total.value = response.data.total || 0
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
  
  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(resource => 
      resource.name.toLowerCase().includes(query)
    )
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

// 服务提供商标签映射
const providerLabels = {
  openai: 'OpenAI',
  zai: 'Zai',
  anthropic: 'Anthropic',
  google: 'Google',
  azure: 'Azure',
  custom: '自定义'
}

// 获取服务提供商标签
const getProviderLabel = (provider) => {
  return providerLabels[provider] || provider
}

// 服务提供商到BaseURL的映射
const providerToBaseURL = {
  openai: 'https://api.openai.com/v1',
  zai: 'https://api.zai.com/v1',
  anthropic: 'https://api.anthropic.com/v1',
  google: 'https://generativelanguage.googleapis.com/v1',
  azure: 'https://your-resource-name.openai.azure.com',
  custom: ''
}

// 处理服务提供商变化
const handleProviderChange = (provider) => {
  if (!provider) return
  
  if (provider === 'custom') {
    // 自定义提供商时，不清空base_url，让用户自己填写或保留原有值
    // 如果是添加模式，清空；如果是编辑模式，保留原值
    if (isAddMode.value) {
      resourceForm.base_url = ''
    }
  } else {
    // 选择预设提供商时，自动填充对应的base_url
    const defaultURL = providerToBaseURL[provider]
    if (defaultURL) {
      // 如果当前base_url为空或者是其他预设提供商的URL，则替换
      // 如果用户已经手动修改过，则询问是否替换（这里简化处理，直接替换）
      if (isAddMode.value || !resourceForm.base_url || Object.values(providerToBaseURL).includes(resourceForm.base_url)) {
        resourceForm.base_url = defaultURL
      }
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
    provider: 'openai',
    model: '',
    base_url: providerToBaseURL['openai'], // 默认填充OpenAI的URL
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
const viewModels = (resourceId) => {
  router.push(`/models?resource_id=${resourceId}`)
}

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