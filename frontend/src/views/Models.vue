<template>
  <div class="models">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>模型管理</span>
          <el-button type="primary" @click="handleAddModel">
            <el-icon><Plus /></el-icon>
            添加模型
          </el-button>
        </div>
      </template>
      
      <!-- 搜索和筛选 -->
      <div class="search-filter">
        <el-input
          v-model="searchQuery"
          placeholder="搜索模型名称"
          prefix-icon="Search"
          style="width: 240px; margin-right: 10px"
        ></el-input>
        <el-select
          v-model="typeFilter"
          placeholder="筛选类型"
          style="width: 120px; margin-right: 10px"
        >
          <el-option label="全部" value=""></el-option>
          <el-option label="聊天" value="chat"></el-option>
          <el-option label="补全" value="completion"></el-option>
          <el-option label="嵌入" value="embedding"></el-option>
          <el-option label="图像" value="image"></el-option>
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
      
      <!-- 模型列表 -->
      <el-table
        v-loading="loading"
        :data="filteredModels"
        style="width: 100%; margin-top: 20px"
        border
      >
        <el-table-column prop="id" label="ID" width="180" />
        <el-table-column prop="name" label="模型名称" />
        <el-table-column prop="type" label="模型类型" width="100">
          <template #default="scope">
            <el-tag>{{ scope.row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="llm_resource_id" label="LLM资源" width="180">
          <template #default="scope">
            {{ getLLMResourceName(scope.row.llm_resource_id) }}
          </template>
        </el-table-column>
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
              @click="handleEditModel(scope.row)"
              style="margin-right: 5px"
            >
              编辑
            </el-button>
            <el-button
              type="info"
              size="small"
              @click="handleViewPricing(scope.row.id)"
              style="margin-right: 5px"
            >
              价格
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDeleteModel(scope.row.id)"
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
    
    <!-- 添加/编辑模型对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
    >
      <el-form
        ref="modelFormRef"
        :model="modelForm"
        :rules="modelRules"
        label-width="100px"
      >
        <el-form-item label="模型名称" prop="name">
          <el-input v-model="modelForm.name" placeholder="请输入模型名称"></el-input>
        </el-form-item>
        <el-form-item label="模型类型" prop="type">
          <el-select v-model="modelForm.type" placeholder="请选择模型类型">
            <el-option label="聊天" value="chat"></el-option>
            <el-option label="补全" value="completion"></el-option>
            <el-option label="嵌入" value="embedding"></el-option>
            <el-option label="图像" value="image"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="LLM资源" prop="llm_resource_id">
          <el-select v-model="modelForm.llm_resource_id" placeholder="请选择LLM资源">
            <el-option
              v-for="resource in llmResources"
              :key="resource.id"
              :label="resource.name"
              :value="resource.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="modelForm.status" placeholder="请选择状态">
            <el-option label="活跃" value="active"></el-option>
            <el-option label="禁用" value="inactive"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="modelForm.description"
            placeholder="请输入模型描述"
            type="textarea"
            rows="3"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveModel">确定</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 模型价格对话框 -->
    <el-dialog
      v-model="pricingVisible"
      title="模型价格信息"
      width="500px"
    >
      <el-card v-if="modelPricing" class="pricing-card">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="模型ID">
            {{ modelPricing.model_id }}
          </el-descriptions-item>
          <el-descriptions-item label="输入价格">
            {{ modelPricing.input_price }} 元/千 tokens
          </el-descriptions-item>
          <el-descriptions-item label="输出价格">
            {{ modelPricing.output_price }} 元/千 tokens
          </el-descriptions-item>
          <el-descriptions-item label="上下文窗口">
            {{ modelPricing.context_window }} tokens
          </el-descriptions-item>
          <el-descriptions-item label="最大输出">
            {{ modelPricing.max_output }} tokens
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(modelPricing.created_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
      <el-empty v-else description="暂无价格信息"></el-empty>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="pricingVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, PriceTag } from '@element-plus/icons-vue'
import api from '../api'

const loading = ref(false)
const models = ref([])
const llmResources = ref([])
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
const modelFormRef = ref(null)
const modelForm = reactive({
  id: '',
  name: '',
  type: 'chat',
  llm_resource_id: '',
  status: 'active',
  description: ''
})

// 价格对话框相关
const pricingVisible = ref(false)
const modelPricing = ref(null)

// 表单验证规则
const modelRules = {
  name: [
    { required: true, message: '请输入模型名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择模型类型', trigger: 'change' }
  ],
  llm_resource_id: [
    { required: true, message: '请选择LLM资源', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 获取模型列表
const getModelList = async () => {
  try {
    loading.value = true
    const response = await api.getModels({ page: currentPage.value, page_size: pageSize.value })
    if (response && response.data) {
      models.value = response.data.items || []
      total.value = response.data.total || 0
    }
  } catch (error) {
    console.error('获取模型列表失败:', error)
    ElMessage.error('获取模型列表失败')
  } finally {
    loading.value = false
  }
}

// 获取LLM资源列表
const getLLMResourceList = async () => {
  try {
    const response = await api.getLLMResources({ page: 1, page_size: 100 })
    if (response && response.data) {
      llmResources.value = response.data.items || []
    }
  } catch (error) {
    console.error('获取LLM资源列表失败:', error)
  }
}

// 过滤模型
const filteredModels = computed(() => {
  let result = models.value
  
  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(model => 
      model.name.toLowerCase().includes(query)
    )
  }
  
  // 类型过滤
  if (typeFilter.value) {
    result = result.filter(model => model.type === typeFilter.value)
  }
  
  // 状态过滤
  if (statusFilter.value) {
    result = result.filter(model => model.status === statusFilter.value)
  }
  
  // 分页
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return result.slice(startIndex, endIndex)
})

// 获取LLM资源名称
const getLLMResourceName = (resourceId) => {
  const resource = llmResources.value.find(r => r.id === resourceId)
  return resource ? resource.name : '未知'
}

// 处理添加模型
const handleAddModel = () => {
  isAddMode.value = true
  dialogTitle.value = '添加模型'
  // 重置表单
  Object.assign(modelForm, {
    id: '',
    name: '',
    type: 'chat',
    llm_resource_id: '',
    status: 'active',
    description: ''
  })
  dialogVisible.value = true
}

// 处理编辑模型
const handleEditModel = (model) => {
  isAddMode.value = false
  dialogTitle.value = '编辑模型'
  // 填充表单
  Object.assign(modelForm, model)
  dialogVisible.value = true
}

// 处理保存模型
const handleSaveModel = async () => {
  try {
    // 验证表单
    await modelFormRef.value.validate()
    
    if (isAddMode.value) {
      // 创建模型
      await api.createModel(modelForm)
      ElMessage.success('模型创建成功')
    } else {
      // 更新模型
      await api.updateModel(modelForm.id, modelForm)
      ElMessage.success('模型更新成功')
    }
    
    // 关闭对话框
    dialogVisible.value = false
    // 重新获取模型列表
    getModelList()
  } catch (error) {
    console.error('保存模型失败:', error)
    ElMessage.error('保存模型失败')
  }
}

// 处理删除模型
const handleDeleteModel = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个模型吗？',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'danger'
      }
    )
    
    await api.deleteModel(id)
    ElMessage.success('模型删除成功')
    // 重新获取模型列表
    getModelList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除模型失败:', error)
      ElMessage.error('删除模型失败')
    }
  }
}

// 处理查看价格
const handleViewPricing = async (id) => {
  try {
    const response = await api.getModelPricing(id)
    if (response && response.data) {
      modelPricing.value = response.data
    } else {
      modelPricing.value = null
    }
    pricingVisible.value = true
  } catch (error) {
    console.error('获取模型价格失败:', error)
    ElMessage.error('获取模型价格失败')
  }
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
  getModelList()
  getLLMResourceList()
})
</script>

<style scoped>
.models {
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

.pricing-card {
  margin-bottom: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .search-filter {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .el-input,
  .el-select {
    width: 100% !important;
  }
}
</style>